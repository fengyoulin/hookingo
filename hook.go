package hookingo

import (
	"errors"
	"reflect"
	"sync"
	"unsafe"
)

// Enabler is the interface the wraps the Enable method.
//
// Enable enables a hook which disabled by the Disable method of the Hook interface.
type Enabler interface {
	Enable()
}

// Hook represents a applied hook, it implements Origin, Disable and Restore. The
// Disable and Restore methods will change the code in the text segment, so are not
// concurrent safe, need special attention.
type Hook interface {
	// Origin returns the origin function, or a error if the origin function is not
	// usable after the hook applied.
	Origin() interface{}
	// Disable temporarily disables the hook and restores the origin function, the
	// hook can be enabled later using the returned Enabler.
	Disable() Enabler
	// Restore restores the origin function permanently, if you want to enable the
	// hook again, you should use the Apply function.
	Restore() error
}

type hook struct {
	// the modified instructions
	target []byte
	// the moved and jump back instructions
	jumper []byte
	// use to call the origin function
	origin interface{}
	// backup the jump to instructions when disabled
	backup []byte
}

func (h *hook) Origin() interface{} {
	return h.origin
}

func (h *hook) Restore() error {
	return remove(slicePtr(h.target))
}

func (h *hook) Disable() Enabler {
	disable(h)
	return &enabler{h:h}
}

type enabler struct {
	h *hook
}

func (e *enabler) Enable() {
	enable(e.h)
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
	// ErrRelativeAddr means cannot call the origin function
	ErrRelativeAddr = errors.New("relative address in instruction")
)

func init() {
	hooks = make(map[uintptr]*hook)
}

// Apply the hook, replace "from" with "to". This function will change the code in
// the text segment, so is not concurrent safe, need special attention.
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
	// early object allocation
	// we may hooking runtime.mallocgc
	// or may be runtime.newobject
	f := &funcval{}
	// early bucket allocation
	hooks[from] = nil
	h, err := applyHook(from, to)
	if err != nil {
		delete(hooks, from) // delete on failure
		return nil, err
	}
	if h.origin == nil {
		f.fn = slicePtr(h.jumper)
		e := (*eface)(unsafe.Pointer(&h.origin))
		e.data = unsafe.Pointer(f)
		e.typ = typ
	}
	// just set value here, should not alloc memory
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
		h.backup = nil
		delete(hooks, from)
		return nil
	}
	return ErrHookNotFound
}

func disable(h *hook) {
	if h.backup == nil {
		b := make([]byte, len(h.target))
		copy(b, h.target)
		h.backup = b
	}
	copy(h.target, h.jumper)
}

func enable(h *hook) {
	copy(h.target, h.backup)
}
