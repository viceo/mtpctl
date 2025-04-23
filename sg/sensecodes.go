package sg

import (
	"strings"
)

func parseSenseCode(senseKey string, asc string, ascq string) string {
	senseKey = strings.ToUpper(senseKey)
	asc = strings.ToUpper(asc)
	ascq = strings.ToUpper(ascq)
	value, ok := spectraSenseCodesMap[senseKey][asc][ascq]
	if ok {
		return value
	} else {
		return ""
	}
}

var spectraSenseCodesMap = map[string]map[string]map[string]string{
	"02": {
		"04": {
			"00": "Unit not ready",
			"01": "Unit is becoming ready",
			"03": "Unit NOT Ready. Manual Intervention Required",
			"83": "Door is open. Robot disabled",
		},
		"3A": {
			"00": "Tape is not loaded and threaded",
		},
	},
	"04": {
		"05": {
			"00": "Logical unit does not respond. Device is permanently inaccessible",
		},
		"2E": {
			"01": "Third party device failure. There are problems communicating with the third party device",
			"02": "Copy target device is unreachable. The target device specified is invalid",
			"04": "Copy target device data under‐run",
			"05": "Copy target device data over-run",
		},
		"40": {
			"D1": "Import/export door could not be extended",
			"D2": "Import/export door could not be retracted",
		},
		"4C": {
			"00": "Unit failed initialization",
		},
		"81": {
			"01": "Drive failed to unload",
			"02": "Tape failed load; move marked successful",
			"04": "Drive failed to come ready",
			"05": "ADI is enabled and iADT failed",
		},
		"85": {
			"01": "Move failed; tape left in picker",
			"02": "Move failed; tape left in source",
			"03": "Move failed; picker will reset",
			"04": "Long axis motor blocked",
			"05": "Gripper motor blocked",
			"06": "Rotary motor blocked",
			"07": "Medium axis motor blocked",
			"08": "Short axis motor blocked",
			"09": "Parameter block is corrupted",
			"0A": "Picker failed to park",
			"0B": "Picker failed initialization. Cannot communicate with barcode scanner",
			"0E": "Cartridge stuck in slot. The robotic picker was unable to pull the tape from the source",
			"25": "Cartridge stuck in drive. The robot was unable to remove the tape from the drive. (The tape is usually left in the mouth of the drive)",
			"90": "No mechanical picker version defined",
			"91": "Calibration block not found",
			"92": "No rack version defined",
			"99": "General robotics failure",
		},
		"86": {
			"00": "Fibre failed initialization",
		},
		"87": {
			"00": "Invalid FPROM / Invalid ID bits",
			"01": "FPROM ERASE operation failed",
			"02": "FPROM WRITE operation failed",
		},
		"88": {
			"00": "General picker definition error",
			"01": "Invalid picker type",
			"02": "Invalid rack type",
			"03": "Invalid library size",
			"04": "Invalid chassis type",
			"05": "Invalid IE door type",
		},
		"90": {
			"00": "Internal SCSI error unknown",
			"01": "Internal SCSI command failed",
			"02": "SCSI command timed out",
			"03": "Internal SCSI command was aborted by host",
			"04": "Initiator detected Error Message Received",
			"05": "Internal SCSI command reselect timeout",
		},
		// WIP... missing all from here to the end of 04
		// "91": {
		// 	"XX": ?????
		// }
	},
	"05": {
		"00": {
			"16": "The initiator is trying to initiate an additional command to the target device before the first command is complete",
		},
		"1A": {
			"00": "Parameter list length error",
		},
		"20": {
			"00": "Invalid command code",
		},
		"21": {
			"01": "Invalid element address",
		},
		"24": {
			"00": "Invalid field in CDB",
		},
		"25": {
			"00": "LUN not supported",
		},
		"26": {
			"00": "Invalid field in parameter list",
			"01": "Parameter not supported",
			"02": "Parameter value invalid",
			"06": "Too many target descriptors",
			"07": "Unsupported target descriptor type code",
			"08": "Too many segment descriptors",
			"09": "Unsupported segment descriptor type code",
			"0A": "Unexpected inexact segment",
			"0B": "Inline data length exceeded",
			"0C": " Invalid operation for copy source or destination",
		},
		"2C": {
			"00": "Command sequence error. FLASH code download to tape drive via Write Buffer failed",
		},
		"2E": {
			"03": "Incorrect copy target device type",
		},
		"39": {
			"00": "An attempt was made to save an emulated mode page. Saving parameters not supported",
		},
		"3B": {
			"0D": "Medium destination is full",
			"0E": "Medium source element empty",
			"11": "Media magazine not accessible",
		},
		"3D": {
			"00": "Identify message error",
			"80": "Disconnects must be allowed",
		},
		"3E": {
			"00": " Could not get wrap information from drive",
		},
		"49": {
			"00": "An attempt was made to issue an ACA Queue command",
		},
		"53": {
			"02": "Media removal prevented",
		},
		"80": {
			"00": "Generic invalid move",
			"01": "Picker not empty",
			"03": "Source magazine not available",
			"04": "Destination magazine not available",
			"05": "Source drive is not available",
			"06": "Destination drive is not available",
			"07": "The medium source is invalid. No barcode label was detected by the robotic picker",
			"18": "Element is reserved for front panel",
		},
		"81": {
			"00": "Duplicate SCSI ID on this bus",
			"01": "TAP is Exit Only",
			"02": "Library is full of tapes. No more tapes may be loaded",
			"03": "Cannot move tape from drive to TAP with Queued Unloads enabled",
			"10": "Unable to insert tape. The TAP is Exit Only",
		},
		"89": {
			"24": "Library is full of tapes",
		},
	},
	"06": {
		"28": {
			"00": "Inventory possibly altered",
			"01": "Door element accessed",
		},
		"29": {
			"00": "A reset has occurred",
			"80": "Drive failed power‐on self test (POST) or user issued diagnostic test. This is a Sony‐unique error code",
		},
		"2A": {
			"01": "Mode parameters have changed",
		},
		"2F": {
			"00": " Commands aborted; cleared by another initiator",
		},
		"3F": {
			"01": "New firmware was loaded successfully",
			"03": "The device’s inquiry data has changed. A new tape library inquiry data was configured for the device",
		},
		"83": {
			"00": "Barcode label is unread",
			"01": "Problem reading barcode label",
			"02": "Tape is queued for unload",
		},
		"84": {
			"01": "No response from SCSI target",
			"02": "Check unexpected condition from target",
		},
	},
}
