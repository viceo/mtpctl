package ibm

import (
	"fmt"
	"strings"

	"github.com/viceo/tplibcmd/scsi"
)

type IBM_TS4500 struct{ scsi.ElementStatus }

func (x IBM_TS4500) NewDataTransferElementDescriptor(buffer []byte, page *scsi.ElementStatusPage) *scsi.DataTransferElementDescriptor {
	descriptor := scsi.ElementStatus{}.NewDataTransferElementDescriptor(buffer, page)
	descriptor.AdditionalSenseValue = x.ascmap(descriptor.AdditionalSenseCode, descriptor.AdditionalSenseCodeQualifier)
	return descriptor
}

func (IBM_TS4500) ascmap(asc string, ascq string) string {
	key := strings.ToUpper(fmt.Sprintf("%s%s", asc, ascq))
	x, ok := ascmap_TS4500[key]
	if !ok {
		x = ""
	}
	return x
}
