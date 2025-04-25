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
		Use:   "device-id sg-device",
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

	rootCmd.AddCommand(
		deviceIdentificationCmd,
		elementStatusCmd,
	)
	if err := rootCmd.Execute(); err != nil {
		util.PanicIfError(err)
	}
}
