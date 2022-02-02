package interface_main

import (
	"eos_bot/internal/req"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"time"

	ui "github.com/gizak/termui/v3"

	"github.com/gizak/termui/v3/widgets"
)

func TermUIGrid(bypass bool) { 
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	stringMachine := func() ([]string, []float64, []string) {
		var data []string
		var newFloats []float64
		var isFiat []string
		data = append(data, "------------------------------------------------------")
		newData := req.CollectJsonFromAPI(bypass)
		for _, v := range newData {
			var l string
			if v.Type != "NULL-DATA" {
				l = fmt.Sprintf("| https://play.upland.me/?prop_id=%s |%s|%s|%s|\n", v.ID, v.Address, v.UPX, v.FIAT)
			} else {
				l = " "
			}
			data = append([]string{l}, data...)
			myFloat := regexp.MustCompile(`\d+\.\d+`)
			getFloat := myFloat.FindString(v.UPX)
			getFloatFIAT := myFloat.FindString(v.FIAT)
			newFiat, err := strconv.ParseFloat(getFloatFIAT, 64)
			if err != nil {
				log.Println("Error converting FIAT to float64")
			} else {
				if newFiat > 0 {
					isFiat = append([]string{l}, isFiat... )
				}
			}
			f, err := strconv.ParseFloat(getFloat, 64)
			if err != nil {
				log.Println("Error converting UPX to float64")
			} else {
				newFloats = append(newFloats, f)
			}
		}
		return data, newFloats, isFiat
	}
	
	data1, _, _ := stringMachine()

	lc := widgets.NewBarChart()
	lc.Title = "Trends"
	lc.BarWidth = 5
	lc.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
	lc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	lc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}

	l := widgets.NewList()
	l.Title = "Latest Properties"
	l.Rows = append(l.Rows, data1...)
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false

	pc := widgets.NewPieChart()
	pc.Title = "Pie Chart"
	pc.AngleOffset = -.5 * math.Pi

	l2 := widgets.NewList()
	l2.Title = "FIAT props"
	l2.TextStyle = ui.NewStyle(ui.ColorYellow)
	l2.WrapText = false



	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)



	grid.Set(

		ui.NewRow(2.0/3,
			ui.NewCol(1.0/4, lc),
			ui.NewCol(1.0/2, l),
			ui.NewCol(1.0/4, pc),

		),
		ui.NewRow(1.0/1,
			ui.NewCol(1.0/1, l2),
		),
	)
	ui.Render(grid)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "j", "<Down>":
				l.ScrollDown()
			case "k", "<Up>":
				l.ScrollUp()
			case "<C-d>":
				l.ScrollHalfPageDown()
			case "<C-u>":
				l.ScrollHalfPageUp()
			case "<C-f>":
				l.ScrollPageDown()
			case "<C-b>":
				l.ScrollPageUp()
			case "<Home>":
				l.ScrollTop()
			case "G", "<End>":
				l.ScrollBottom()
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(grid)
			}

		case <-ticker:
			go func() {
				data1, data2, isFiat := stringMachine()
				if len(data2) > 0 {
					if len(lc.Data) >= 3 && len(lc.Data) < 5 {
						lc.Data = append(lc.Data[:3], data2...)
						pc.Data = append(pc.Data[:3], data2...)
					} else {
						lc.Data = data2
						pc.Data = data2
					}
				}
				if len(isFiat) > 0 {
					l2.Rows = append(l2.Rows, isFiat...)
				}
				l.Rows = append(data1, l.Rows...)	
			}()
			ui.Render(grid)
		}
	}
}