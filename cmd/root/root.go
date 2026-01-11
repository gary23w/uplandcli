package root

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// Verbose is the log level
var Verbose string = "warning"

// Version of the application
var Version string = "0.1.0"

var RootCmd = &cobra.Command{
Use:     "uplandcli",
Version: Version,
Short:   "Upland CLI interacts with the blockchain",
Long: `++UPLND-CLI++
=================
The UPLND-CLI will interact with the blockchain to collect data related to Upland properties.
This data will be used to populate the CLI based user interface.
Dy`,
PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
fmt.Printf("Setting up UPLD-CLI-%v: \n", Version)

config := zap.Config{}
//config.Level.SetLevel(zap.NewAtomicLevelAt(zap.DebugLevel))
//configure the logger
jsonfile, err := os.Open("conf/logging.json")
	if err != nil {
		log.Fatalln("Couldn't open the json file", err)
		return err
	}
	defer jsonfile.Close()
	if err := json.NewDecoder(jsonfile).Decode(&config); err != nil {
		return err
	}
	logger, err := config.Build()
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(logger)
	logger.Debug("Running prerun with log level:", zap.String("log_level", Verbose))
	return nil
},
}

func ExecuteCLI() {
	zap.S().Debug("Running the main execute function with log level:", zap.String("log_level", Verbose))
	if err := RootCmd.Execute(); err != nil {
		zap.S().Error("upload failed", zap.Error(err))
		os.Exit(1)
	}
	zap.S().Debugf("Completed the main execute function with log level: %s", Verbose)
}

func SetLogs(notice string) {
	file, err := os.OpenFile("tmp/debug.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    log.SetOutput(file)
    log.Print(notice)
}