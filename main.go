package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	inquiry "github.com/viceo/tplibcmd/cmd/0x12"
	readElementStatus "github.com/viceo/tplibcmd/cmd/0xb8"
	"github.com/viceo/tplibcmd/util"
)

func main() {
	defer errorHandler()

	paths, err := filepath.Glob("/dev/sg*")
	util.PanicIfError(err)

	var mediaChangers []inquiry.DeviceIdentification

	for _, p := range paths {
		// Open SCSI device
		device, err := os.OpenFile(p, os.O_RDWR, 0)
		util.PanicIfError(err)
		defer device.Close()

		// Inquiry device
		idPage := inquiry.Run(device)
		if idPage.PheripherialDeviceType == 8 {
			mediaChangers = append(mediaChangers, idPage)
		}

	}

	elementStatus := readElementStatus.Run(mediaChangers[0].Device)

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
	HasError      bool                                  `json:"hasError"`
	MediaChangers []inquiry.DeviceIdentification        `json:"mediaChangers"`
	ElementStatus readElementStatus.ElementStatusHeader `json:"elementStatus"`
}

// Recover panics and log them in JSON format ending process.
// Manually JSON formatting because marshaling can return error
func errorHandler() {
	if r := recover(); r != nil {
		fmt.Printf("{\"hasError\":true, \"error\": \"%s\"}\n", r)
		fmt.Println(string(debug.Stack()))
	}
}
