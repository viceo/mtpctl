package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/viceo/tplibcmd/cli"
	"github.com/viceo/tplibcmd/util"
)

func main() {
	defer util.ErrorHandler()

	var rootCmd = &cobra.Command{
		Use:   "tplibcmd",
		Short: "CLI tool to interact with tape libraries",
	}

	var deviceIdentificationCmd = &cobra.Command{
		Use:   "device-id [sg-device]",
		Short: "Get device identification",
		Args:  cobra.ExactArgs(1),
		Run:   cli.DeviceIdentificationCmd,
	}

	var elementStatusCmd = &cobra.Command{
		Use:   "element-status",
		Short: "Get element status page",
		Run:   cli.ExecElementStatusCmd,
	}

	rootCmd.AddCommand(
		deviceIdentificationCmd,
		elementStatusCmd,
	)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
