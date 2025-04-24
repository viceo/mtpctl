package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/viceo/tplibcmd/cmd"
	"github.com/viceo/tplibcmd/cmd/ibm"
	"github.com/viceo/tplibcmd/cmd/spectra"
	"github.com/viceo/tplibcmd/util"
)

func ExecElementStatusCmd(c *cobra.Command, args []string) {
	paths, err := filepath.Glob("/dev/sg*")
	util.PanicIfError(err)

	var mediaChangers []cmd.DeviceIdentification

	for _, p := range paths {
		// Open SCSI device
		device, err := os.OpenFile(p, os.O_RDWR, 0)
		util.PanicIfError(err)
		defer device.Close()

		// Inquiry device
		devIdPage := cmd.RunDeviceIdentification(device)
		if devIdPage.Page.PheripherialDeviceType == 8 {
			mediaChangers = append(mediaChangers, devIdPage)
		}

	}

	var elementStatusImpl cmd.IElementStatus
	switch strings.ToUpper(mediaChangers[0].Page.VendorIdentification) {
	case "SPECTRA":
		elementStatusImpl = spectra.SPECTRA_TFINITY{}
	default:
		elementStatusImpl = ibm.IBM_TS4500{}
	}
	elementStatus := cmd.RunElementStatus(elementStatusImpl, mediaChangers[0].Page.Device)

	response, err := json.Marshal(elementStatus)
	util.PanicIfError(err)

	// Print JSON result
	fmt.Println(string(response))
}
