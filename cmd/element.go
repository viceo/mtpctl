package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/viceo/tplibcmd/dto"
	"github.com/viceo/tplibcmd/scsi"
	"github.com/viceo/tplibcmd/scsi/ibm"
	"github.com/viceo/tplibcmd/scsi/spectra"
	"github.com/viceo/tplibcmd/util"
)

func ExecElementStatusCmd(_ *cobra.Command, args []string) {
	sgDevice := args[0]
	// Open SCSI device
	device, err := os.OpenFile(sgDevice, os.O_RDWR, 0)
	util.PanicIfError(err)
	defer device.Close()

	devIdPage := scsi.RunDeviceIdentification(device)
	if devIdPage.Page.PheripherialDeviceType != 8 {
		panic(fmt.Errorf("command unsupported on device type %d", devIdPage.Page.PheripherialDeviceType))
	}

	var elementStatusImpl scsi.IElementStatus
	switch strings.ToUpper(devIdPage.Page.VendorIdentification) {
	case "SPECTRA":
		elementStatusImpl = spectra.SPECTRA_TFINITY{}
	default:
		elementStatusImpl = ibm.IBM_TS4500{}
	}
	(&dto.CommandResponse{Data: scsi.RunElementStatus(elementStatusImpl, devIdPage.Page.Device)}).Print()
}
