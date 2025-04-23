package cmd

import (
	"os"
	"strings"

	"github.com/viceo/tplibcmd/sg"
	"github.com/viceo/tplibcmd/util"
)

type DeviceIdentification struct {
	cmd                    sg.SgCmd
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
}

func Run(device *os.File) DeviceIdentification {
	cmd := sg.SgCmd{
		Cdb:            []byte{0x12, 0x01, 0x83, 0x00, 0xFF, 0x00},
		DataBuffer:     make([]byte, 96),
		SenseBuffer:    make([]byte, 32),
		DxferDirection: sg.SG_DXFER_FROM_DEV,
		Timeout:        uint32(30 * 1000), // 30 seconds
		Flags:          uint32(0),
	}

	syscallerr, scsierr := sg.ExecCmd(cmd, device)
	util.PanicIfError(syscallerr)
	util.PanicIfError(scsierr)

	return newDeviceIdentification(cmd, device)
}

func newDeviceIdentification(cmd sg.SgCmd, device *os.File) DeviceIdentification {
	return DeviceIdentification{
		Device:                 device,
		DeviceName:             device.Name(),
		PheripherialQualifier:  cmd.DataBuffer[0] >> 4,
		PheripherialDeviceType: cmd.DataBuffer[0] & 0x0F,
		PageCode:               cmd.DataBuffer[1],
		PageLength:             cmd.DataBuffer[3],
		CodeSet:                cmd.DataBuffer[4] & 0x0F,
		IdentifierType:         cmd.DataBuffer[5] & 0x0F,
		IdentifierLength:       cmd.DataBuffer[7],
		VendorIdentification:   strings.TrimSpace(string(cmd.DataBuffer[8:16])),
		ProductIdentification:  strings.TrimSpace(string(cmd.DataBuffer[16:32])),
		UnitSerialNumber:       strings.TrimSpace(string(cmd.DataBuffer[32:42])),
	}
}
