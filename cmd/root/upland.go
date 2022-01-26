package root

import (
	"fmt"

	"github.com/spf13/cobra"
)
	
var (
	dt bool = false
	qt bool = false
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
		
		The UPLD-PIPELINE will also scrape the Upland website and collect data via a headless browser.
		using chromedp and chromedp-go for headless browsing.`,
	Run: UplandPipeline,
	}
}

func init() {
	upldCmd := NewUPLDCmd()
	upldCmd.Flags().BoolVarP(&dt, "collect", "d", false, "will get all of the recent properties listed for sale")
	upldCmd.Flags().BoolVarP(&qt, "live", "q", false, "live scrape the upland website")
	RootCmd.AddCommand(upldCmd)
}

func UplandPipeline(cmd *cobra.Command, args []string) {
	if dt {
		TermUIpanel("Collecting data from blockchain")
	}
	if qt {
		fmt.Println("offline until implemented")
	}
}