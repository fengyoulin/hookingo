//go:build windows
// +build windows

package hookingo

import (
	"golang.org/x/sys/windows"
)

func init() {
	pageSize = uintptr(windows.Getpagesize())
}

// allocate a new readable, writable and executable page
func allocPage() (uintptr, error) {
	return windows.VirtualAlloc(0, pageSize, windows.MEM_COMMIT, windows.PAGE_EXECUTE_READWRITE)
}

// make the pages readable, writable and executable
func protectPages(addr, size uintptr) error {
	start := pageSize * (addr / pageSize)
	length := pageSize * ((addr + size + pageSize - 1 - start) / pageSize)
	var old uint32
	for i := uintptr(0); i < length; i += pageSize {
		err := windows.VirtualProtect(start+i, pageSize, windows.PAGE_EXECUTE_READWRITE, &old)
		if err != nil {
			return err
		}
	}
	return nil
}
