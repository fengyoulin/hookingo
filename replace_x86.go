//go:build 386 || amd64
// +build 386 amd64

package hookingo

import (
	"golang.org/x/arch/x86/x86asm"
	"unsafe"
)

func findCall(from, target uintptr, limit int) (addrs []uintptr, err error) {
	lmt := limit
	var off uintptr
	for {
		src := makeSlice(from+off, 32)
		var inst x86asm.Inst
		inst, err = x86asm.Decode(src, 64)
		if err != nil {
			return
		}
		if inst.Op == x86asm.CALL && inst.Len == 5 {
			if 0xe8 == *(*byte)(unsafe.Pointer(from + off)) {
				ta := from + off + 5 + uintptr(*(*int32)(unsafe.Pointer(from + off + 1)))
				if ta == target {
					addrs = append(addrs, from+off+1)
					lmt--
				}
			}
		} else if inst.Op == x86asm.RET {
			break
		}
		off += uintptr(inst.Len)
		if lmt == 0 {
			break
		}
	}
	return
}

func batchProtect(addrs []uintptr) (err error) {
	for i := 0; i < len(addrs); i++ {
		if err = protectPages(addrs[i], 4); err != nil {
			return
		}
	}
	return nil
}

func setCall(addrs []uintptr, target uintptr) {
	for i := 0; i < len(addrs); i++ {
		addr := addrs[i]
		*(*int32)(unsafe.Pointer(addr)) = int32(target - (addr + 4))
	}
}
