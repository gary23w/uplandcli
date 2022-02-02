package live

import (
	"fmt"
	"time"

	"eos_bot/internal/models"
	"eos_bot/internal/req"
)

func TailDatabaseTables(bypass bool) { 
	var data []models.DataPackageBLOCK
	for {
		data = req.CollectJsonFromAPI(bypass)
			for _, v := range data {
				var l string
				if v.Type != "NULL-DATA" {
					l = fmt.Sprintf("[*] %s | %s | https://play.upland.me/?prop_id=%s | %s\n", v.UPX, v.FIAT, v.ID, v.Address)
					// loop must sleep or it will break.
					//time.Sleep(3 * time.Second)
				} else {
					ti := time.Now()
					l = "[*] No data available at " + ti.String() + "\n"
					time.Sleep(time.Minute)
				}
				fmt.Println(l)
			}
	}
}



