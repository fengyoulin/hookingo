package hookingo

import (
	"golang.org/x/arch/arm64/arm64asm"
)

func analysis(src []byte) (instLen int, err error) {
	_, err = arm64asm.Decode(src)
	if err != nil {
		return 0, err
	}
	return 4, nil
}
