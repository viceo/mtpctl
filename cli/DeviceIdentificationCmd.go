package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/viceo/tplibcmd/cmd"
	"github.com/viceo/tplibcmd/util"
)

func DeviceIdentificationCmd(c *cobra.Command, args []string) {
	sgDevice := args[0]
	device, err := os.OpenFile(sgDevice, os.O_RDWR, 0)
	util.PanicIfError(err)
	defer device.Close()
	deviceIdentification := cmd.RunDeviceIdentification(device)
	response, err := json.Marshal(deviceIdentification)
	util.PanicIfError(err)
	fmt.Println(string(response))
}
