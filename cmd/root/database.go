package root

import (
	"fmt"

	db "eos_bot/internal/database"

	"github.com/spf13/cobra"
)
	
var (
	ht bool = false
	dtt bool = false
	ct bool = false
)

func NewDBCMD() *cobra.Command {
return &cobra.Command{
	Use:   "database",
	Short: "setup a postgresql database on heroku",
	Long: `++UPLD-DB UPLAND++
		========================
		The DB command is used to setup and initialize a postgresql database on heroku.`,
	Run: DBroutes,
	}
}

func init() {
	dbCMD := NewDBCMD()
	dbCMD.Flags().BoolVarP(&ht, "deploy", "d", false, "will setup a postgresql database on heroku")
	dbCMD.Flags().BoolVarP(&dtt, "destroy", "u", false, "will attempt to destroy the database")
	dbCMD.Flags().BoolVarP(&ct, "check", "q", false, "checks to see if a database is already active")
	RootCmd.AddCommand(dbCMD)
}

func DBroutes(cmd *cobra.Command, args []string) {
	if ht {
		db.DeployHeroku("upland-cli")
	}
	if dtt {
		db.DestroyPostgres()
	}
	if ct {
		fmt.Println("offline until implemented")
	}
}