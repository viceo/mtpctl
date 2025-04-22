package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	cmd0x12 "github.com/viceo/tplibcmd/cmd0x12"
	"github.com/viceo/tplibcmd/util"
)

func main() {
	defer errorHandler()

	paths, err := filepath.Glob("/dev/sg*")
	util.PanicIfError(err)

	var mediaChangers []cmd0x12.DeviceIdentificationPage

	for _, p := range paths {
		// Open SCSI device
		device, err := os.OpenFile(p, os.O_RDWR, 0)
		util.PanicIfError(err)
		defer device.Close()

		// Inquiry device
		idPage := cmd0x12.NewDeviceIdentificationPage(device)
		if idPage.PheripherialDeviceType == 8 {
			mediaChangers = append(mediaChangers, idPage)
		}

	}

	response, err := json.Marshal(jsonResponse{
		MediaChangers: mediaChangers,
	})
	util.PanicIfError(err)

	// Print JSON result
	fmt.Println(string(response))
}

// JSON response always expects no errors
type jsonResponse struct {
	HasError      bool                               `json:"hasError"`
	MediaChangers []cmd0x12.DeviceIdentificationPage `json:"mediaChangers"`
}

// Recover panics and log them in JSON format ending process.
// Manually JSON formatting because marshaling can return error
func errorHandler() {
	if r := recover(); r != nil {
		fmt.Printf("{\"hasError\":true, \"error\": \"%s\"}\n", r)
	}
}
