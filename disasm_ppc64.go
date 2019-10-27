package hookingo

import (
	"encoding/binary"
	"golang.org/x/arch/ppc64/ppc64asm"
)

func analysis(src []byte) (instLen int, err error) {
	inst, err := ppc64asm.Decode(src, binary.BigEndian)
	if err != nil {
		return 0, err
	}
	return inst.Len, nil
}
