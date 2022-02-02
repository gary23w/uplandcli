package root

import (
	api "eos_bot/api/props_crud"
	uiInt "eos_bot/cmd/interface_main"
	"eos_bot/internal/live"

	"github.com/spf13/cobra"
)
	
var (
	dt bool = false
	qt bool = false
	apiq bool = false
	bypass bool = false
)

func NewUPLDCmd() *cobra.Command {
return &cobra.Command{
	Use:   "upland",
	Short: "run the upland pipeline",
	Long: `++UPLD-CLI UPLAND++
		========================
		The UPLD-PIPELINE will query from the blockchain and collect data related
		to Upland properties. This data will be used to populate the CLI based user interface.

		Example:
			upldcli upland --collect 
			upldcli upland --live
			upldcli upland --live -a  // run API in async mode
			upldcli upland --live -a -b  // run API in async mode and bypass database connections
		
		The UPLD-PIPELINE will also scrape the Upland website and collect data via a headless browser.
		using chromedp and chromedp-go for headless browsing. This is a future implementation, and should be available soon.`,
	Run: UplandPipeline,
	}
}

func init() {
	upldCmd := NewUPLDCmd()
	upldCmd.Flags().BoolVarP(&dt, "collect", "d", false, "will get all of the recent properties listed for sale.")
	upldCmd.Flags().BoolVarP(&qt, "live", "q", false, "live mode which tails collected data in your shell.")
	upldCmd.Flags().BoolVarP(&apiq, "api", "a", false, "run API in async mode")
	upldCmd.Flags().BoolVarP(&bypass, "bypass", "b", false, "bypass database connections and inserts")

	RootCmd.AddCommand(upldCmd)
}

func UplandPipeline(cmd *cobra.Command, args []string) {
		if dt {
			//uiInt.TermUIpanel("Collecting data from blockchain", bypass)
			SetLogs("[UPLD-CLI] USING TERM UI. LIVE MODE DISABLED")
			uiInt.TermUIGrid(bypass)
		}
		if qt {
			if apiq {
				go api.StartAPI()
			}
			live.TailDatabaseTables(bypass)
		}
}