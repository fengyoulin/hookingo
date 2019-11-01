package hookingo

const (
	jumperSize = 32
)

func applyHook(from, to uintptr) (*hook, error) {
	src := makeSlice(from, 32)
	l, err := ensureLength(src, 7)
	if err != nil {
		return nil, err
	}
	err = protectPages(from, uintptr(l))
	if err != nil {
		return nil, err
	}
	dst, err := allocJumper()
	if err != nil {
		return nil, err
	}
	src = makeSlice(from, uintptr(l))
	copy(dst, src)
	addr := from + uintptr(l)
	inst := []byte{
		0x50,                               // PUSH EAX
		0x50,                               // PUSH EAX
		0xb8,                               // MOV EAX, addr
		byte(addr), byte(addr >> 8),        // .
		byte(addr >> 16), byte(addr >> 24), // .
		0x89, 0x44, 0x24, 0x04,             // MOV [ESP+4], EAX
		0x58,                               // POP EAX
		0xc3,                               // RET
	}
	ret := makeSlice(slicePtr(dst)+uintptr(l), uintptr(len(dst)-l))
	copy(ret, inst)
	for i := l + len(inst); i < len(dst); i++ {
		dst[i] = 0xcc
	}
	addr = to
	inst = []byte{
		0xb8,                               // MOV EAX, addr
		byte(addr), byte(addr >> 8),        // .
		byte(addr >> 16), byte(addr >> 24), // .
		0xff, 0xe0,                         // JMP EAX
	}
	copy(src, inst)
	for i := len(inst); i < len(src); i++ {
		src[i] = 0xcc
	}
	hk := &hook{
		target: src,
		jumper: dst,
	}
	return hk, nil
}
