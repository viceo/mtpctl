package main

import (
	"github.com/spf13/cobra"
	"github.com/viceo/tplibcmd/cmd"
	"github.com/viceo/tplibcmd/util"
)

func main() {
	defer util.ErrorTrap()

	var rootCmd = &cobra.Command{
		Use:   "tplibcmd",
		Short: "CLI tool to interact with tape libraries",
	}

	var deviceIdentificationCmd = &cobra.Command{
		Use:   "device-id [sg-device]",
		Short: "Get device identification",
		Args:  cobra.ExactArgs(1),
		Run:   cmd.DeviceIdentificationCmd,
		// SilenceErrors: true,
	}

	var elementStatusCmd = &cobra.Command{
		Use:   "element-status [sg-device]",
		Short: "Get element status page",
		Args:  cobra.ExactArgs(1),
		Run:   cmd.ExecElementStatusCmd,
	}

	var logSenseCdm = &cobra.Command{
		Use:   "log-sense [sg-device]",
		Short: "Get log sense flags (from a tape drive)",
		Args:  cobra.ExactArgs(1),
		Run:   cmd.LogSenseCmd,
	}

	rootCmd.AddCommand(
		deviceIdentificationCmd,
		elementStatusCmd,
		logSenseCdm,
	)
	if err := rootCmd.Execute(); err != nil {
		util.PanicIfError(err)
	}
}
