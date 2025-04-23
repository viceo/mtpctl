package cmd

import (
	"encoding/binary"
	"os"

	"github.com/viceo/tplibcmd/sg"
	"github.com/viceo/tplibcmd/util"
)

type ElementStatusHeader struct {
	FirstElementAddressReported uint16 `json:"firstElementAddressReported"`
	NumberOfElementsReported    uint16 `json:"numberOfElementsReported"`
	ElementStatusPagesByteCount uint32 `json:"elementStatusPagesByteCount"`
}

type ElementStatusPage struct {
	ElementTypeCode             uint8  `json:"elementTypeCode"`
	PVolTag                     bool   `json:"pvoltag"`
	AVolTag                     bool   `json:"avoltag"`
	ElementDescriptorsLength    uint16 `json:"elementDescriptorsLength"`
	ElementDescriptorsByteCount uint32 `json:"elementDescriptorsByteCount"`
}

type DataTransferElementDescriptor struct {
	ElementAddress               uint16 `json:"elementAddress"`
	Access                       bool   `json:"access"`
	Except                       bool   `json:"except"`
	Full                         bool   `json:"full"`
	AdditionalSenseCode          string `json:"asc"`
	AdditionalSenseCodeQualifier string `json:"ascq"`
	SValid                       bool   `json:"svalid"`
	Invert                       bool   `json:"invert"`
	SourceStorageElementAddress  uint16 `json:"sourceStorageElementAddress"`
	PVolTag                      string `json:"pvoltag"`
	AVolTag                      string `json:"avoltag"`
	CodeSet                      uint8  `json:"codeset"`
	IdentifierType               uint8  `json:"identifierType"`
	IdentifierLength             uint8  `json:"identifierLength"`
	DeviceIdentifier             string `json:"deviceIdentifier"`
}

type ElementStatus struct {
	Header      ElementStatusHeader             `json:"header"`
	Page        ElementStatusPage               `json:"pages"`
	Descriptors []DataTransferElementDescriptor `json:"descriptors"`
	SenseData   sg.SgSenseData                  `json:"senseData"`
}

type IElementStatus interface {
	NewElementStatusPage(buffer []byte) ElementStatusPage
	NewDataTransferElementDescriptor(buffer []byte, page *ElementStatusPage) DataTransferElementDescriptor
}

func RunElementStatus[T IElementStatus](impl T, device *os.File) ElementStatus {
	cmd := sg.SgCmd{
		Cdb: []byte{
			0xB8, /* Operation Code */
			0x04, /* bit4: VolTag, bit3-0: Element Type Code */
			0x00, /* Starting Element Address */
			0x00, /* Starting Element Address */
			0xFF, /* Number of elements */
			0xFF, /* Number of elements */
			0x01, /* bit1: CurData, bit0: DVCID */
			0x00, /* Allocation length */
			0xFF, /* Allocation length */
			0xFF, /* Allocation length */
			0x00,
			0x00,
		},
		DataBuffer:     make([]byte, 64*1000),
		SenseBuffer:    make([]byte, 16),
		DxferDirection: sg.SG_DXFER_FROM_DEV,
		Timeout:        uint32(30 * 1000), // 30 seconds
		Flags:          uint32(0),
	}

	syscallerr := sg.ExecCmd(&cmd, device)
	util.PanicIfError(syscallerr)

	// Generic behavior
	header := newElementStatusHeader(&cmd)

	// Specific behavior
	page := impl.NewElementStatusPage(cmd.GetDataSlice(8, 16))

	var dataTransferElementDescriptorList []DataTransferElementDescriptor
	var indx uint16
	for indx = range header.NumberOfElementsReported {
		dataTransferElementDescriptorList = append(
			dataTransferElementDescriptorList,
			impl.NewDataTransferElementDescriptor(cmd.GetDataSlice(
				(indx*page.ElementDescriptorsLength)+16,
				((indx+1)*page.ElementDescriptorsLength)+16,
			), &page),
		)
	}

	return ElementStatus{
		Header:      header,
		Page:        page,
		Descriptors: dataTransferElementDescriptorList,
		SenseData:   cmd.GetSenseData(),
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
