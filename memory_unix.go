// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package hookingo

import (
	"golang.org/x/sys/unix"
)

func init() {
	pageSize = uintptr(unix.Getpagesize())
}

func allocPage() (uintptr, error) {
	data, err := unix.Mmap(-1, 0, int(pageSize), unix.PROT_EXEC|unix.PROT_READ|unix.PROT_WRITE, unix.MAP_ANONYMOUS)
	if err != nil {
		return 0, err
	}
	return slicePtr(data), nil
}

func protectPages(addr, size uintptr) error {
	start := pageSize * (addr / pageSize)
	length := pageSize * ((addr + size + pageSize - 1 - start) / pageSize)
	for i := uintptr(0); i < length; i += pageSize {
		data := makeSlice(start+1, pageSize)
		err := unix.Mprotect(data, unix.PROT_EXEC|unix.PROT_READ|unix.PROT_WRITE)
		if err != nil {
			return err
		}
	}
	return nil
}
