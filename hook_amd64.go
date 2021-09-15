package hookingo

const (
	jumperSize = 64
)

func applyHook(from, to uintptr) (*hook, error) {
	src := makeSlice(from, 32)
	inf, err := ensureLength(src, 19)
	if err != nil {
		return nil, err
	}
	err = protectPages(from, uintptr(inf.length))
	if err != nil {
		return nil, err
	}
	dst, err := allocJumper()
	if err != nil {
		return nil, err
	}
	// early object allocation
	hk := &hook{}
	src = makeSlice(from, uintptr(inf.length))
	copy(dst, src)
	addr := from + uintptr(inf.length)
	inst := []byte{
		0x50,                               // PUSH RAX
		0x50,                               // PUSH RAX
		0x48, 0xb8,                         // MOV RAX, addr
		byte(addr), byte(addr >> 8),        // .
		byte(addr >> 16), byte(addr >> 24), // .
		byte(addr >> 32), byte(addr >> 40), // .
		byte(addr >> 48), byte(addr >> 56), // .
		0x48, 0x89, 0x44, 0x24, 0x08,       // MOV [RSP+8], RAX
		0x58,                               // POP RAX
		0xc3,                               // RET
	}
	ret := makeSlice(slicePtr(dst)+uintptr(inf.length), uintptr(len(dst)-inf.length))
	copy(ret, inst)
	for i := inf.length + len(inst); i < len(dst); i++ {
		dst[i] = 0xcc
	}
	addr = to
	inst = []byte{
		0x50,                               // PUSH RAX
		0x50,                               // PUSH RAX
		0x48, 0xb8,                         // MOV RAX, addr
		byte(addr), byte(addr >> 8),        // .
		byte(addr >> 16), byte(addr >> 24), // .
		byte(addr >> 32), byte(addr >> 40), // .
		byte(addr >> 48), byte(addr >> 56), // .
		0x48, 0x89, 0x44, 0x24, 0x08,       // MOV [RSP+8], RAX
		0x58,                               // POP RAX
		0xc3,                               // RET
	}
	copy(src, inst)
	for i := len(inst); i < len(src); i++ {
		src[i] = 0xcc
	}
	hk.target = src
	hk.jumper = dst
	if !inf.relocatable {
		hk.origin = ErrRelativeAddr
	}
	return hk, nil
}
