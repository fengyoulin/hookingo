package hookingo

import (
	"golang.org/x/arch/x86/x86asm"
)

func analysis(src []byte) (instLen int, err error) {
	inst, err := x86asm.Decode(src, 64)
	if err != nil {
		return 0, err
	}
	for _, a := range inst.Args {
		if mem, ok := a.(x86asm.Mem); ok {
			if mem.Base == x86asm.RIP {
				return 0, ErrRelativeAddr
			}
		}
	}
	return inst.Len, nil
}
