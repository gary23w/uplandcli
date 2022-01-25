package root

import (
	"eos_bot/internal/models"
	"eos_bot/internal/req"
	"fmt"
	"log"
	"time"

	ui "github.com/gizak/termui/v3"

	"github.com/gizak/termui/v3/widgets"
)


func TermUIpanel(ty string) {
	var check string
	for loop := true; loop; {
		if err := ui.Init(); err != nil {
			log.Fatalf("failed to initialize termui: %v", err)
		}
		termWidth, termHeight := ui.TerminalDimensions()
		defer ui.Close()

		var newRows []string
		var data []models.DataPackageBLOCK
		data = req.CollectJsonFromAPI()
		// loop through data list
		for _, v := range data {
			l := fmt.Sprintf("[*] %s: %s\n", v.Type, v.ID)
			newRows = append(newRows, l)
		}
		l := widgets.NewList()
		
		l.Title = ty + check
		// check newRows size	
		if len(newRows) > 0 {
			l.Rows = newRows
		} else {
			l.Rows = []string{"[*] No data found"}
		}
		l.Rows = newRows
		l.TextStyle = ui.NewStyle(ui.ColorYellow)
		l.WrapText = false
		l.SetRect(0, 0, termWidth, termHeight)
		
		ui.Render(l)
		previousKey := ""

		uiEvents := ui.PollEvents()
		var timer int = 0
		for {
			e := <-uiEvents
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
				case "g":
					if previousKey == "g" {
						l.ScrollTop()
					}
				case "<Home>":
					l.ScrollTop()
				case "G", "<End>":
					l.ScrollBottom()
			}
			if previousKey == "g" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}
			ui.Render(l)
			if timer == 10 {
				break
			}
			timer++
			time.Sleep(1 * time.Second)
		}
		check = fmt.Sprintf(" :: %v", time.Now().Format(time.RFC3339))
	}
}