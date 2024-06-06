package hookingo_test

import (
	"fmt"
	"github.com/fengyoulin/hookingo"
	"testing"
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
		s := f2()
		o := h.Origin()
		if f, ok := o.(func() string); ok {
			s += f()
		} else if e, ok := o.(error); ok {
			return "", e
		}
		e := h.Disable()
		s += f2()
		e.Enable()
		s += f2()
		err = h.Restore()
		if err != nil {
			return "", err
		}
		h = nil
		s += f2()
		return s, nil
	}()
	if err != nil {
		t.Error(err)
	} else if s != "f2f3f1f2f1f2f3f2f1" {
		t.Error(s)
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
		s := f2()
		s += f1()
		e := h.Disable()
		s += f2()
		e.Enable()
		s += f2()
		e = h.Disable()
		h = nil
		s += f2()
		return s, nil
	}()
	if err != nil {
		t.Error(err)
	} else if s != "f2f3f1f2f1f2f3f2f1" {
		t.Error(s)
	}
}

func f1() string {
	s := "f1"
	fmt.Print(s)
	return s
}

func f2() string {
	s := "f2"
	fmt.Print(s)
	return s + f1()
}

func f3() string {
	s := "f3"
	fmt.Print(s)
	return s
}
