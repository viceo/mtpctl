package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/viceo/tplibcmd/dto"
	"github.com/viceo/tplibcmd/scsi"
	"github.com/viceo/tplibcmd/util"
)

func LogSenseCmd(_ *cobra.Command, args []string) {
	sgDevice := args[0]
	// Open SCSI device
	device, err := os.OpenFile(sgDevice, os.O_RDWR, 0)
	util.PanicIfError(err)
	defer device.Close()

	devIdPage := scsi.RunDeviceIdentification(device)
	if devIdPage.Page.PheripherialDeviceType != 1 {
		panic(fmt.Errorf("command unsupported on device type %d", devIdPage.Page.PheripherialDeviceType))
	}

	// logSenseImpl := scsi.LogSense{}
	logSense := scsi.RunLogSense(device)
	(&dto.CommandResponse{Data: logSense}).Print()
}
