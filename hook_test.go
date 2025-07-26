package hookingo_test

import (
	"testing"

	"github.com/fengyoulin/hookingo"
)

func TestApply(t *testing.T) {
	s, err := func() (string, error) {
		h, err := hookingo.Apply(f1, f3)
		if err != nil {
			return "", err
		}
		defer func() {
			if h != nil {
				_ = h.Restore()
			}
		}()
		s := f2("f")
		o := h.Origin()
		if f, ok := o.(func(string) string); ok {
			s += f("F")
		} else if e, ok := o.(error); ok {
			return "", e
		}
		e := h.Disable()
		s += f2("f")
		e.Enable()
		s += f2("F")
		err = h.Restore()
		if err != nil {
			return "", err
		}
		h = nil
		s += f2("f")
		return s, nil
	}()
	if err != nil {
		t.Error(err)
	} else if x := "ffffFffFFFFff"; s != x {
		t.Errorf("%s != %s", s, x)
	}
}

func TestReplace(t *testing.T) {
	s, err := func() (string, error) {
		h, err := hookingo.Replace(f2, f1, f3, -1)
		if err != nil {
			return "", err
		}
		defer func() {
			if h != nil {
				h.Disable()
			}
		}()
		s := f2("f")
		s += f1("F")
		e := h.Disable()
		s += f2("f")
		e.Enable()
		s += f2("F")
		_ = h.Disable()
		h = nil
		s += f2("f")
		return s, nil
	}()
	if err != nil {
		t.Error(err)
	} else if x := "ffffFffFFFFff"; s != x {
		t.Errorf("%s != %s", s, x)
	}
}

//go:noinline
func f1(s string) string {
	return s
}

//go:noinline
func f2(s string) string {
	return s + f1(s)
}

//go:noinline
func f3(s string) string {
	return s + s + s
}
