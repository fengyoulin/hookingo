package hookingo

import (
	"encoding/binary"
	"golang.org/x/arch/ppc64/ppc64asm"
)

func analysis(src []byte) (inf info, err error) {
	inst, err := ppc64asm.Decode(src, binary.LittleEndian)
	if err != nil {
		return
	}
	inf.length = inst.Len
	return
}
