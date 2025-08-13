package scsi

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/viceo/tplibcmd/sg"
	"github.com/viceo/tplibcmd/util"
)

type TapeAlertFlag struct {
	Name     string `json:"name"`
	HexCode  string `json:"hexCode"`
	HexValue uint16 `json:"hexValue"`
	Length   uint8  `json:"length"`
	Value    uint8  `json:"value"`
}

type LogSense struct {
	PageCode   uint8           `json:"pageCode"`
	PageLength uint16          `json:"pageLength"`
	Flags      []TapeAlertFlag `json:"flags"`
}

func RunLogSense(device *os.File) LogSense {
	cmd := sg.SgCmd{
		Cdb: []byte{
			0x4D,
			0x00,
			0x2E,
			0x00,
			0x00,
			0x00,
			0x00,
			0xFF,
			0xFF,
			0x00,
		},
		DataBuffer:     make([]byte, 2*1000),
		SenseBuffer:    make([]byte, 16),
		DxferDirection: sg.SG_DXFER_FROM_DEV,
		Timeout:        uint32(30 * 1000),
		Flags:          uint32(0),
	}

	syscallerr := sg.ExecCmd(&cmd, device)
	util.PanicIfError(syscallerr)
	// fmt.Printf("%+v", cmd)

	pageLength := binary.BigEndian.Uint16(cmd.DataBuffer[2:4])
	logSense := LogSense{
		PageCode:   cmd.DataBuffer[0],
		PageLength: pageLength,
		Flags:      getFlags(cmd.DataBuffer[4 : 4+pageLength]),
	}

	return logSense
}

func getFlags(buffer []byte) []TapeAlertFlag {
	flags := []TapeAlertFlag{}
	iterator := 0
	buffLen := len(buffer)

	for iterator+4 <= buffLen { // at least header exists
		flagLength := int(buffer[iterator+3])

		// Safety check: don't go out of bounds
		if iterator+4+flagLength > buffLen {
			break
		}

		flag := TapeAlertFlag{
			HexValue: binary.BigEndian.Uint16(buffer[iterator : iterator+2]),
			Length:   uint8(flagLength),
		}

		// String Hex Code Value
		flag.HexCode = fmt.Sprintf("0x%04X", flag.HexValue)
		flag.Name = TapeAlertFlagNames[flag.HexValue]

		// Read only first byte of data (for TapeAlert, it's usually 0 or 1)
		if flagLength > 0 {
			flag.Value = buffer[iterator+4]
		}

		flags = append(flags, flag)

		iterator += 4 + flagLength // move to next entry
	}

	return flags
}

var TapeAlertFlagNames = map[uint16]string{
	0x01: "Read warning",
	0x02: "Write warning",
	0x03: "Hard error",
	0x04: "Media",
	0x05: "Read failure",
	0x06: "Write failure",
	0x07: "Media life",
	0x08: "Not data grade",
	0x09: "Write protect",
	0x0A: "No removal",
	0x0B: "Cleaning media",
	0x0C: "Unsupported format",
	0x0E: "Unrecoverable snapped tape",
	0x0F: "Cartridge memory chip failure",
	0x10: "Forced eject",
	0x11: "Read-only format",
	0x12: "Tape directory corrupted in cartridge memory",
	0x13: "Nearing media life",
	0x14: "Clean now",
	0x15: "Clean periodic",
	0x16: "Expired clean",
	0x17: "Invalid cleaning tape",
	0x19: "Interface",
	0x1A: "Cooling fan failure",
	0x1B: "Power supply",
	0x1E: "Hardware A",
	0x1F: "Hardware B",
	0x20: "Interface",
	0x21: "Eject media",
	0x22: "Download fail",
	0x23: "Drive humidity",
	0x24: "Drive temperature",
	0x25: "Drive voltage",
	0x26: "Predictive failure",
	0x27: "Failure",
	0x31: "Diminished native capacity",
	0x33: "Tape directory invalid at unload",
	0x34: "Tape system area write failure",
	0x35: "Tape system area read failure",
	0x37: "Load failure",
	0x38: "Unrecoverable unload failure",
	0x3B: "WORM Medium – Integrity Check Failed",
	0x3C: "WORM Medium – Overwrite Attempted",
	0x3D: "Encryption policy violation",
}
