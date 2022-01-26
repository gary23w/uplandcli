package root

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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
Dy`,
PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
fmt.Printf("Setting up UPLD-CLI-%v: \n", Version)

config := zap.Config{}
//config.Level.SetLevel(zap.NewAtomicLevelAt(zap.DebugLevel))
//configure the logger
jsonfile, err := os.Open("utils/logging.json")
	if err != nil {
		log.Fatalln("Couldn't open the json file", err)
		return err
	}
	defer jsonfile.Close()

	var buf strings.Builder
	_, err = io.Copy(&buf, jsonfile)
	if err != nil {
		// handle error
		log.Fatalln("Could not convert file to string", err)
		return err
	}
	s := buf.String()
	if err := json.Unmarshal([]byte(s), &config); err != nil {
		panic(err)
	}
	logger, err := config.Build()
	// zap.ReplaceGlobals(logger)
	if err != nil {
		panic(err)
	}
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

