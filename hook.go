package hookingo

import (
	"errors"
	"reflect"
	"sync"
	"unsafe"
)

// Hook interface
type Hook interface {
	Origin() interface{}
	Restore() error
}

type hook struct {
	// the modified instructions
	target []byte
	// the moved and jump back instructions
	jumper []byte
	// use to call the origin function
	origin interface{}
}

func (h *hook) Origin() interface{} {
	return h.origin
}

func (h *hook) Restore() error {
	return remove(slicePtr(h.target))
}

var (
	// hooks applied with target addresses as keys
	hooks map[uintptr]*hook
	// protect the hooks map
	lock sync.Mutex
)

var (
	// ErrDoubleHook means already hooked
	ErrDoubleHook = errors.New("double hook")
	// ErrHookNotFound means the hook not found
	ErrHookNotFound = errors.New("hook not found")
	// ErrDifferentType means from and to are of different types
	ErrDifferentType = errors.New("inputs are of different type")
	// ErrInputType means inputs are not func type
	ErrInputType = errors.New("inputs are not func type")
)

func init() {
	hooks = make(map[uintptr]*hook)
}

// Apply the hook
func Apply(from, to interface{}) (Hook, error) {
	vf := reflect.ValueOf(from)
	vt := reflect.ValueOf(to)
	if vf.Type() != vt.Type() {
		return nil, ErrDifferentType
	}
	if vf.Kind() != reflect.Func {
		return nil, ErrInputType
	}
	e := (*eface)(unsafe.Pointer(&from))
	return apply(vf.Pointer(), vt.Pointer(), e.typ)
}

func apply(from, to uintptr, typ unsafe.Pointer) (*hook, error) {
	lock.Lock()
	defer lock.Unlock()
	_, ok := hooks[from]
	if ok {
		return nil, ErrDoubleHook
	}
	h, err := applyHook(from, to)
	if err != nil {
		return nil, err
	}
	f := &funcval{fn: slicePtr(h.jumper)}
	e := (*eface)(unsafe.Pointer(&h.origin))
	e.data = unsafe.Pointer(f)
	e.typ = typ
	hooks[from] = h
	return h, nil
}

func remove(from uintptr) error {
	lock.Lock()
	defer lock.Unlock()
	h, ok := hooks[from]
	if ok {
		copy(h.target, h.jumper)
		freeJumper(h.jumper)
		h.jumper = nil
		h.target = nil
		h.origin = nil
		delete(hooks, from)
		return nil
	}
	return ErrHookNotFound
}
