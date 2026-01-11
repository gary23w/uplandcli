package home

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gary23w/uplandcli/internal/req"

	ui "github.com/gizak/termui/v3"

	"github.com/gizak/termui/v3/widgets"
)

var floatRe = regexp.MustCompile(`\d+\.\d+`)

func TermUIGrid(bypassed bool) { 
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	stringMachine := func() ([]string, []float64, []string) {
		var data []string
		var newFloats []float64
		var isFiat []string
		var EOSHttpReq req.EOSHTTPREQ
		EOSHttpReq.Bypass_sql = bypassed
		newData := EOSHttpReq.CollectJsonFromAPI()
		data = make([]string, 0, len(newData)+1)
		newFloats = make([]float64, 0, len(newData))
		isFiat = make([]string, 0, len(newData))

		// Preserve existing visual ordering without O(n^2) prepends.
		for i := len(newData) - 1; i >= 0; i-- {
			v := newData[i]
			var l string
			if v.Type != "NULL-DATA" {
				l = fmt.Sprintf("| https://play.upland.me/?prop_id=%s |%s|%s|%s|\n", v.ID, v.Address, v.UPX, v.FIAT)
			} else {
				l = " "
			}
			data = append(data, l)
			getFloat := floatRe.FindString(v.UPX)
			getFloatFIAT := floatRe.FindString(v.FIAT)
			newFiat, err := strconv.ParseFloat(getFloatFIAT, 64)
			if err != nil {
				log.Println("Error converting FIAT to float64")
			} else {
				if newFiat > 0 {
					isFiat = append(isFiat, l)
				}
			}
			f, err := strconv.ParseFloat(getFloat, 64)
			if err != nil {
				log.Println("Error converting UPX to float64")
			} else {
				newFloats = append(newFloats, f)
			}
		}
		data = append(data, "------------------------------------------------------")
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
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	results := make(chan struct {
		data1  []string
		data2  []float64
		isFiat []string
	}, 1)
	var inFlight int32
	
	startFetch := func() {
		if !atomic.CompareAndSwapInt32(&inFlight, 0, 1) {
			return
		}
		go func() {
			defer atomic.StoreInt32(&inFlight, 0)
			data1, data2, isFiat := stringMachine()
			select {
			case results <- struct {
				data1  []string
				data2  []float64
				isFiat []string
			}{data1: data1, data2: data2, isFiat: isFiat}:
			default:
				// If UI is behind, drop this tick's update.
			}
		}()
	}
	startFetch()

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

		case <-ticker.C:
			startFetch()
		case r := <-results:
			if len(r.data2) > 0 {
				if len(lc.Data) >= 3 && len(lc.Data) < 5 {
					lc.Data = append(lc.Data[:3], r.data2...)
					pc.Data = append(pc.Data[:3], r.data2...)
				} else {
					lc.Data = r.data2
					pc.Data = r.data2
				}
			}
			if len(r.isFiat) > 0 {
				l2.Rows = append(l2.Rows, r.isFiat...)
			}
			l.Rows = append(r.data1, l.Rows...)
			ui.Render(grid)
		}
	}
}