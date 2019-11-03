package hookingo

import (
	"golang.org/x/arch/arm64/arm64asm"
)

func analysis(src []byte) (inf info, err error) {
	_, err = arm64asm.Decode(src)
	if err != nil {
		return
	}
	inf.length = 4
	return
}
