package spectra

import (
	"bytes"
	"encoding/binary"
	"fmt"

	_cmd "github.com/viceo/tplibcmd/cmd"
)

type SPECTRA_TFINITY struct{}

func (SPECTRA_TFINITY) NewElementStatusPage(buffer []byte) _cmd.ElementStatusPage {
	return _cmd.ElementStatusPage{
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

func (SPECTRA_TFINITY) NewDataTransferElementDescriptor(buffer []byte, page *_cmd.ElementStatusPage) _cmd.DataTransferElementDescriptor {
	var pvoltag, offsetindx = "", 12
	if page.PVolTag {
		pvoltag = string(buffer[offsetindx : offsetindx+36])
		// Adjust offset
		offsetindx += 36
	}
	identifierLength := uint8(buffer[offsetindx+3])
	deviceIdentifierStartPos := uint8(offsetindx + 4)
	return _cmd.DataTransferElementDescriptor{
		ElementAddress:               binary.BigEndian.Uint16(buffer[0:2]),
		Access:                       (buffer[2]&0x08)>>3 == 1,
		Except:                       (buffer[2]&0x04)>>2 == 1,
		Full:                         (buffer[2] & 0x01) == 1,
		AdditionalSenseCode:          fmt.Sprintf("%02x", buffer[4]),
		AdditionalSenseCodeQualifier: fmt.Sprintf("%02x", buffer[5]),
		SValid:                       (buffer[9]&0x80)>>7 == 1,
		Invert:                       (buffer[9]&0x40)>>6 == 1,
		SourceStorageElementAddress:  binary.BigEndian.Uint16(buffer[10:12]),
		PVolTag:                      pvoltag,
		AVolTag:                      string(""), /* DVCID = 1 doesn't containt AVoltag */
		CodeSet:                      buffer[offsetindx] & 0x0F,
		IdentifierType:               buffer[offsetindx+1] & 0x0F,
		IdentifierLength:             identifierLength,
		DeviceIdentifier:             string(bytes.Trim(buffer[deviceIdentifierStartPos:deviceIdentifierStartPos+identifierLength], "\x00")),
	}
}
