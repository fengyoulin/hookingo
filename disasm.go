package hookingo

// info of one instruction
type info struct {
	length      int
	relocatable bool
}

func ensureLength(src []byte, size int) (info, error) {
	var inf info
	inf.relocatable = true
	for inf.length < size {
		i, err := analysis(src)
		if err != nil {
			return inf, err
		}
		inf.relocatable = inf.relocatable && i.relocatable
		inf.length += i.length
		src = src[i.length:]
	}
	return inf, nil
}
