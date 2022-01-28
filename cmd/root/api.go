package root

import (
	api "eos_bot/api/props_crud"

	"github.com/spf13/cobra"
)
	
var (
	htt bool = false
)

func NewAPICMD() *cobra.Command {
return &cobra.Command{
	Use:   "api",
	Short: "launch a crud api to interact with the database",
	Long: `++UPLD-DB UPLAND++
		========================
		The api command is used to deploy a crud api to interact with the database.`,
	Run: APIroutes,
	}
}

func init() {
	apiCMD := NewAPICMD()
	apiCMD.Flags().BoolVarP(&htt, "deploy", "d", false, "will initialize a crud api to interact with the database")
	RootCmd.AddCommand(apiCMD)
}

func APIroutes(cmd *cobra.Command, args []string) {
	if htt {
		api.StartAPI()
	}
}