package live

import (
	"fmt"
	"strings"
	"time"

	"github.com/gary23w/uplandcli/internal/models"
	"github.com/gary23w/uplandcli/internal/req"
)

func TailDatabaseTables(bypassed bool) { 
	var data []models.DataPackageBLOCK
	for {
		var EOSHttpReq req.EOSHTTPREQ
		EOSHttpReq.Bypass_sql = bypassed
		data = EOSHttpReq.CollectJsonFromAPI()
			for _, v := range data {
				if v.Type != "NULL-DATA" {
					var b strings.Builder
					b.Grow(128)
					b.WriteString("[*] ")
					b.WriteString(v.UPX)
					b.WriteString(" | ")
					b.WriteString(v.FIAT)
					b.WriteString(" | https://play.upland.me/?prop_id=")
					b.WriteString(v.ID)
					b.WriteString(" | ")
					b.WriteString(v.Address)
					fmt.Println(b.String())
				} else {
					ti := time.Now()
					fmt.Println("[*] No data available at " + ti.String())
					time.Sleep(3 * time.Second)
				}
			}
	}
}



