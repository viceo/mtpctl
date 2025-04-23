package spectra

import (
	"fmt"
	"strings"

	"github.com/viceo/tplibcmd/cmd"
)

type SPECTRA_TFINITY struct{ cmd.ElementStatus }

func (x SPECTRA_TFINITY) NewDataTransferElementDescriptor(buffer []byte, page *cmd.ElementStatusPage) *cmd.DataTransferElementDescriptor {
	descriptor := cmd.ElementStatus{}.NewDataTransferElementDescriptor(buffer, page)
	descriptor.AdditionalSenseValue = x.ascmap(descriptor.AdditionalSenseCode, descriptor.AdditionalSenseCodeQualifier)
	return descriptor
}

func (SPECTRA_TFINITY) ascmap(asc string, ascq string) string {
	key := strings.ToUpper(fmt.Sprintf("%s%s", asc, ascq))
	x, ok := ascmap_TFINITY[key]
	if !ok {
		x = ""
	}
	return x
}
