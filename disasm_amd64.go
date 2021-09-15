package hookingo

import (
	"golang.org/x/arch/x86/x86asm"
)

func analysis(src []byte) (inf info, err error) {
	inst, err := x86asm.Decode(src, 64)
	if err != nil {
		return
	}
	inf.length = inst.Len
	inf.relocatable = true
	for _, a := range inst.Args {
		if mem, ok := a.(x86asm.Mem); ok {
			if mem.Base == x86asm.RIP {
				inf.relocatable = false
				return
			}
		} else if _, ok := a.(x86asm.Rel); ok {
			inf.relocatable = false
			return
		}
	}
	return
}
