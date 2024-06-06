package hookingo

import (
	"reflect"
)

// HookCaller applies a group of hooks in the caller, by changing the relative
// address in some call instructions. It is not concurrent safe, need special
// attention.
type HookCaller interface {
	// Disable disables the hooks and restores the calls to the original function,
	// the hook can be enabled later using the returned Enabler.
	Disable() Enabler
	// Count returns the total number of modified call instructions. If the hooks
	// are disabled, it returns 0.
	Count() int
}

type hookCaller struct {
	old uintptr
	new uintptr
	pos []uintptr
	dis bool
}

func (h *hookCaller) Disable() Enabler {
	setCall(h.pos, h.old)
	h.dis = true
	return &enableCaller{h: h}
}

func (h *hookCaller) Count() int {
	if h.dis {
		return 0
	}
	return len(h.pos)
}

type enableCaller struct {
	h *hookCaller
}

func (e *enableCaller) Enable() {
	setCall(e.h.pos, e.h.new)
	e.h.dis = false
}

// Replace the calls to "old" with "new" in caller, without modify any instructions in "old"
func Replace(caller, old, new interface{}, limit int) (h HookCaller, err error) {
	vf := reflect.ValueOf(old)
	vt := reflect.ValueOf(new)
	if vf.Type() != vt.Type() {
		return nil, ErrDifferentType
	}
	if vf.Kind() != reflect.Func {
		return nil, ErrInputType
	}
	vc := reflect.ValueOf(caller)
	if vc.Kind() != reflect.Func {
		return nil, ErrInputType
	}
	as, err := findCall(vc.Pointer(), vf.Pointer(), limit)
	if err != nil {
		return nil, err
	}
	if err = batchProtect(as); err != nil {
		return nil, err
	}
	hk := &hookCaller{}
	setCall(as, vt.Pointer())
	hk.pos = as
	hk.old = vf.Pointer()
	hk.new = vt.Pointer()
	return hk, nil
}
