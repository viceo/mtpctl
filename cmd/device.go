package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/viceo/tplibcmd/dto"
	"github.com/viceo/tplibcmd/scsi"
	"github.com/viceo/tplibcmd/util"
)

func DeviceIdentificationCmd(_ *cobra.Command, args []string) {
	sgDevice := args[0]
	device, err := os.OpenFile(sgDevice, os.O_RDWR, 0)
	util.PanicIfError(err)
	defer device.Close()
	(&dto.CommandResponse{Data: scsi.RunDeviceIdentification(device)}).Print()
}
