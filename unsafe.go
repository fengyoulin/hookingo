package hookingo

import (
	"unsafe"
)

type eface struct {
	typ  unsafe.Pointer
	data unsafe.Pointer
}

type funcval struct {
	fn uintptr
}

func makeSlice(addr, size uintptr) (bs []byte) {
	sh := (*[3]uintptr)(unsafe.Pointer(&bs))
	sh[0] = addr
	sh[1] = uintptr(size)
	sh[2] = uintptr(size)
	return
}

func slicePtr(bs []byte) uintptr {
	sh := (*[3]uintptr)(unsafe.Pointer(&bs))
	return sh[0]
}
