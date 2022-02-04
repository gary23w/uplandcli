package live

import (
	"fmt"
	"time"

	"eos_bot/internal/models"
	"eos_bot/internal/req"
)

func TailDatabaseTables(bypassed bool) { 
	var data []models.DataPackageBLOCK
	for {
		var EOSHttpReq req.EOSHTTPREQ
		EOSHttpReq.Bypass_sql = bypassed
		data = EOSHttpReq.CollectJsonFromAPI()
			for _, v := range data {
				var l string
				if v.Type != "NULL-DATA" {
					l = fmt.Sprintf("[*] %s | %s | https://play.upland.me/?prop_id=%s | %s\n", v.UPX, v.FIAT, v.ID, v.Address)
				} else {
					ti := time.Now()
					l = "[*] No data available at " + ti.String() + "\n"
					time.Sleep(time.Second * 3)
				}
				fmt.Println(l)
			}
	}
}



