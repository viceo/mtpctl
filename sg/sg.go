package sg

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// SCSI Generic (sg) ioctl constants
const (
	SG_IO = 0x2285 // ioctl command for SCSI generic I/O
)

// Data transfer direction constants
const (
	SG_DXFER_NONE        = -1
	SG_DXFER_TO_DEV      = -2
	SG_DXFER_FROM_DEV    = -3
	SG_DXFER_TO_FROM_DEV = -4
)

// Interface IDs constants
const (
	SG_INTERFACE_0   = 0
	SG_INTERFACE_V3  = 'S'
	SG_INTERFACE_V4  = 'Q'
	SG_INTERFACE_USB = 'U'
)

type SgIoHdr struct {
	InterfaceID    int32
	DxferDirection int32
	CmdLen         uint8
	MxSbLen        uint8
	IovecCount     uint16
	DxferLen       uint32
	Dxferp         uintptr
	Cmdp           uintptr
	Sbp            uintptr
	Timeout        uint32
	Flags          uint32
	PackID         int32
	UsrPtr         uintptr
	Status         uint8
	MaskedStatus   uint8
	MsgStatus      uint8
	SbLenWr        uint8
	HostStatus     uint16
	DriverStatus   uint16
	Resid          int32
	Duration       uint32
	Info           uint32
}

type SgCmd struct {
	Cdb            []byte
	DataBuffer     []byte
	SenseBuffer    []byte
	DxferDirection int32
	Timeout        uint32
	Flags          uint32
	Status         uint8
}

type SgSenseData struct {
	SenseLength uint8  `json:"senseLength"`
	SenseKey    string `json:"senseKey"`
	Asc         string `json:"asc"`
	Ascq        string `json:"ascq"`
}

func (x SgCmd) GetDataSlice(from uint16, to uint16) (clone []byte) {
	clone = make([]byte, len(x.DataBuffer[from:to]))
	copy(clone, x.DataBuffer[from:to])
	return clone
}

func (x *SgCmd) GetSenseData() SgSenseData {
	return SgSenseData{
		SenseLength: uint8(len(x.SenseBuffer)),
		SenseKey:    fmt.Sprintf("%02x", x.SenseBuffer[2]&0x0F),
		Asc:         fmt.Sprintf("%02x", x.SenseBuffer[12]),
		Ascq:        fmt.Sprintf("%02x", x.SenseBuffer[13]),
	}
}

func ExecCmd(cmd *SgCmd, device *os.File) (syscallerr error) {
	// Setup sg_io_hdr structure
	hdr := SgIoHdr{
		InterfaceID:    SG_INTERFACE_V3,
		DxferDirection: cmd.DxferDirection,
		CmdLen:         uint8(len(cmd.Cdb)),         // Command length (cdb)
		MxSbLen:        uint8(len(cmd.SenseBuffer)), // Maximum Sensebuffer length
		DxferLen:       uint32(len(cmd.DataBuffer)), // Data Transfer length
		Dxferp:         uintptr(unsafe.Pointer(&cmd.DataBuffer[0])),
		Cmdp:           uintptr(unsafe.Pointer(&cmd.Cdb[0])),
		Sbp:            uintptr(unsafe.Pointer(&cmd.SenseBuffer[0])),
		Timeout:        cmd.Timeout, // 30 seconds in milliseconds
		Flags:          cmd.Flags,
	}

	// Execute the SCSI command via ioctl
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(device.Fd()),
		uintptr(SG_IO),
		uintptr(unsafe.Pointer(&hdr)),
	)

	if errno != 0 {
		syscallerr = fmt.Errorf("IOCTL status:%v", errno)
	}

	return syscallerr
}
