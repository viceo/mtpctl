package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/viceo/tplibcmd/cmd"
	"github.com/viceo/tplibcmd/cmd/ibm"
	"github.com/viceo/tplibcmd/cmd/spectra"

	"github.com/viceo/tplibcmd/util"
)

func main() {
	defer errorHandler()

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

	response, err := json.Marshal(jsonResponse{
		MediaChangers: mediaChangers,
		ElementStatus: elementStatus,
	})

	util.PanicIfError(err)

	// Print JSON result
	fmt.Println(string(response))
}

// JSON response always expects no errors
type jsonResponse struct {
	HasError      bool                       `json:"hasError"`
	MediaChangers []cmd.DeviceIdentification `json:"mediaChangers"`
	ElementStatus cmd.ElementStatus          `json:"elementStatus"`
}

// Recover panics and log them in JSON format ending process.
// Manually JSON formatting because marshaling can return error
func errorHandler() {
	if r := recover(); r != nil {
		var stackTrace string
		stackTrace = string(debug.Stack())
		stackTrace = strings.ReplaceAll(stackTrace, "\n", "")
		stackTrace = strings.ReplaceAll(stackTrace, "\t", " ")
		stackTrace = strconv.Quote(stackTrace)
		fmt.Printf("{\"hasError\":true, \"error\": \"%s\", \"trace\": %s}\n", r, stackTrace)
	}
}
