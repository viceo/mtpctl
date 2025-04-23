package cmd

import (
	"encoding/binary"
	"os"

	"github.com/viceo/tplibcmd/sg"
	"github.com/viceo/tplibcmd/util"
)

type ElementStatusPages struct {
	ElementTypeCode             uint8  `json:"elementTypeCode"`
	PVolTag                     bool   `json:"pvoltag"`
	AVolTag                     bool   `json:"avoltag"`
	ElementDescriptorsLength    uint16 `json:"elementDescriptorsLength"`
	ElementDescriptorsByteCount uint32 `json:"elementDescriptorsByteCount"`
}

type ElementStatusHeader struct {
	FirstElementAddressReported uint16 `json:"firstElementAddressReported"`
	NumberOfElementsReported    uint16 `json:"numberOfElementsReported"`
	ElementStatusPagesByteCount uint32 `json:"elementStatusPagesByteCount"`
}

type ElementStatus struct {
	Header    ElementStatusHeader `json:"header"`
	Pages     ElementStatusPages  `json:"pages"`
	SenseData sg.SgSenseData      `json:"senseData"`
}

func RunElementStatus(device *os.File) ElementStatus {
	cmd := sg.SgCmd{
		Cdb:            []byte{0xB8, 0x04, 0x00, 0x00, 0xFF, 0xFF, 0x01, 0x00, 0xFF, 0xFF, 0x00, 0x00},
		DataBuffer:     make([]byte, 64*1000),
		SenseBuffer:    make([]byte, 16),
		DxferDirection: sg.SG_DXFER_FROM_DEV,
		Timeout:        uint32(30 * 1000), // 30 seconds
		Flags:          uint32(0),
	}

	syscallerr := sg.ExecCmd(&cmd, device)
	util.PanicIfError(syscallerr)

	return ElementStatus{
		Header:    newElementStatusHeader(&cmd),
		Pages:     newElementStatusPages(&cmd),
		SenseData: cmd.GetSenseData(),
	}
}

func newElementStatusHeader(cmd *sg.SgCmd) ElementStatusHeader {
	buffer := cmd.DataBuffer[0:8]
	return ElementStatusHeader{
		FirstElementAddressReported: binary.BigEndian.Uint16(buffer[0:2]),
		NumberOfElementsReported:    binary.BigEndian.Uint16(buffer[2:4]),
		// Pack the three bytes in a uint32 (4 bytes)
		ElementStatusPagesByteCount: uint32(buffer[5])<<16 |
			uint32(buffer[6])<<8 |
			uint32(buffer[7]),
	}
}

func newElementStatusPages(cmd *sg.SgCmd) ElementStatusPages {
	buffer := cmd.DataBuffer[8:16]
	return ElementStatusPages{
		ElementTypeCode:          buffer[0],
		PVolTag:                  (buffer[1]&0x80)>>7 == 1,
		AVolTag:                  (buffer[1]&0x40)>>6 == 1,
		ElementDescriptorsLength: binary.BigEndian.Uint16(buffer[2:4]),
		// Pack the three bytes in a uint32 (4 bytes)
		ElementDescriptorsByteCount: uint32(buffer[5])<<16 |
			uint32(buffer[6])<<8 |
			uint32(buffer[7]),
	}
}
