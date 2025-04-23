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
		idPage := cmd.RunDeviceIdentification(device)
		if idPage.PheripherialDeviceType == 8 {
			mediaChangers = append(mediaChangers, idPage)
		}

	}

	elementStatus := cmd.RunElementStatus(mediaChangers[0].Device)

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
