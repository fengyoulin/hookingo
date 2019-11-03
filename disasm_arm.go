package hookingo

import (
	"golang.org/x/arch/arm/armasm"
)

func analysis(src []byte) (inf info, err error) {
	inst, err := armasm.Decode(src, armasm.ModeARM)
	if err != nil {
		return
	}
	inf.length = inst.Len
	return
}
