package hookingo

var (
	pageSize uintptr
	freePool [][]byte
)

func allocJumper() ([]byte, error) {
	if len(freePool) == 0 {
		page, err := allocPage()
		if err != nil {
			return nil, err
		}
		for addr := page; addr < page+pageSize; addr += jumperSize {
			freePool = append(freePool, makeSlice(addr, jumperSize))
		}
	}
	l := len(freePool)
	j := freePool[l-1]
	freePool = freePool[:l-1]
	return j, nil
}

func freeJumper(j []byte) {
	freePool = append(freePool, j)
}
