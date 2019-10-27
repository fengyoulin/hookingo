package hookingo

import (
	"fmt"
	"testing"
)

func TestHook(t *testing.T) {
	s, err := f4()
	if err != nil {
		t.Error(err)
	}
	if s != "f4f2f3f1f2f1" {
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

func f4() (string, error) {
	s := "f4"
	fmt.Print(s)
	h, err := Apply(f1, f3)
	if err != nil {
		return "", err
	}
	s += f2()
	o := h.Origin()
	f, ok := o.(func() string)
	if ok {
		s += f()
	}
	err = h.Restore()
	if err != nil {
		return "", err
	}
	s += f2()
	return s, nil
}
