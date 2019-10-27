package hookingo

import (
	"golang.org/x/arch/arm/armasm"
)

func analysis(src []byte) (instLen int, err error) {
	inst, err := armasm.Decode(src, armasm.ModeARM)
	if err != nil {
		return 0, err
	}
	return inst.Len, nil
}
