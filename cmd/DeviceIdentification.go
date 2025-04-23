package cmd

import (
	"os"
	"strings"

	"github.com/viceo/tplibcmd/sg"
	"github.com/viceo/tplibcmd/util"
)

type DeviceIdentification struct {
	Device                 *os.File `json:"-"`
	DeviceName             string   `json:"device"`
	PheripherialQualifier  uint8    `json:"pheripherialQualifier"`
	PheripherialDeviceType uint8    `json:"pheripherialDeviceType"`
	PageCode               uint8    `json:"pageCode"`
	PageLength             uint8    `json:"pageLength"`
	CodeSet                uint8    `json:"codeSet"`
	IdentifierType         uint8    `json:"identifierType"`
	IdentifierLength       uint8    `json:"identifierLength"`
	VendorIdentification   string   `json:"vendorIdentification"`
	ProductIdentification  string   `json:"productIdentification"`
	UnitSerialNumber       string   `json:"unitSerialNumber"`
	SenseBuffer            string   `json:"senseBuffer"`
}

func RunDeviceIdentification(device *os.File) DeviceIdentification {
	cmd := sg.SgCmd{
		Cdb:            []byte{0x12, 0x00, 0x83, 0x00, 0xFF, 0x00},
		DataBuffer:     make([]byte, 96),
		SenseBuffer:    make([]byte, 32),
		DxferDirection: sg.SG_DXFER_FROM_DEV,
		Timeout:        uint32(30 * 1000), // 30 seconds
		Flags:          uint32(0),
	}

	syscallerr, scsierr := sg.ExecCmd(&cmd, device)
	util.PanicIfError(syscallerr)
	util.PanicIfError(scsierr)

	return newDeviceIdentification(&cmd, device)
}

func newDeviceIdentification(cmd *sg.SgCmd, device *os.File) DeviceIdentification {
	buffer := cmd.DataBuffer[0:42]
	return DeviceIdentification{
		Device:                 device,
		DeviceName:             device.Name(),
		PheripherialQualifier:  buffer[0] >> 4,
		PheripherialDeviceType: buffer[0] & 0x0F,
		PageCode:               buffer[1],
		PageLength:             buffer[3],
		CodeSet:                buffer[4] & 0x0F,
		IdentifierType:         buffer[5] & 0x0F,
		IdentifierLength:       buffer[7],
		VendorIdentification:   strings.TrimSpace(string(buffer[8:16])),
		ProductIdentification:  strings.TrimSpace(string(buffer[16:32])),
		UnitSerialNumber:       strings.TrimSpace(string(buffer[32:42])),
		SenseBuffer:            strings.TrimSpace(string(cmd.SenseBuffer)),
	}
}
