package cmd

import (
	"os"
	"strings"

	"github.com/viceo/tplibcmd/sg"
	"github.com/viceo/tplibcmd/util"
)

type DeviceIdentificationPage struct {
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

type DeviceIdentification struct {
	Page      DeviceIdentificationPage `json:"page"`
	SenseData sg.SgSenseData           `json:"senseData"`
}

type IDeviceIdentification interface {
	Init()
	GetCDB() []byte
}

func RunDeviceIdentification(device *os.File) DeviceIdentification {
	cmd := sg.SgCmd{
		Cdb: []byte{
			0x12,
			0x01,
			0x83,
			0x00,
			0xFF,
			0x00,
		},
		DataBuffer:     make([]byte, 128),
		SenseBuffer:    make([]byte, 16),
		DxferDirection: sg.SG_DXFER_FROM_DEV,
		Timeout:        uint32(30 * 1000), // 30 seconds
		Flags:          uint32(0),
	}

	syscallerr := sg.ExecCmd(&cmd, device)
	util.PanicIfError(syscallerr)

	return DeviceIdentification{
		Page:      newDeviceIdentificationPage(&cmd, device),
		SenseData: cmd.GetSenseData(),
	}
}

func newDeviceIdentificationPage(cmd *sg.SgCmd, device *os.File) DeviceIdentificationPage {
	buffer := cmd.DataBuffer[0:42]
	return DeviceIdentificationPage{
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
	}
}
