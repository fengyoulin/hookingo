package hookingo

import (
	"reflect"
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
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	sh.Data = addr
	sh.Len = int(size)
	sh.Cap = int(size)
	return
}

func slicePtr(bs []byte) uintptr {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	return sh.Data
}
