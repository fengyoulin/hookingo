# hookingo #

The name of this project reads like **hook in go** or **hooking go**, it is just an implementation of the "function hook" technology in golang.

### Usage
```shell script
$ go get github.com/fengyoulin/hookingo
```
Import the hooking package:
```go
import "github.com/fengyoulin/hookingo"
```
The document:
```go
var (
	// ErrDoubleHook means the function already hooked, you cannot hook it again.
	ErrDoubleHook = errors.New("double hook")
	// ErrHookNotFound means the hook not found in the applied hooks, maybe you
	// are trying to restore a corrupted hook.
	ErrHookNotFound = errors.New("hook not found")
	// ErrDifferentType means "from" and "to" are of different types, you should
	// only replace a function with one of the same type.
	ErrDifferentType = errors.New("inputs are of different type")
	// ErrInputType means either "from" or "to" are not func type, cannot apply.
	ErrInputType = errors.New("inputs are not func type")
	// ErrRelativeAddr means you cannot call the origin function with the hook
	// applied, try disable the hook, pay special attention to the concurrency.
	ErrRelativeAddr = errors.New("relative address in instruction")
)
```

#### type Enabler

```go
type Enabler interface {
	Enable()
}
```

Enabler is the interface the wraps the Enable method.

Enable enables a hook which disabled by the Disable method of the Hook
interface. This method will change the code in the text segment, so is
not concurrent safe, need special attention.

#### type Hook

```go
type Hook interface {
	// Origin returns the original function, or an error if the original function
	// is not usable after the hook applied.
	Origin() interface{}
	// Disable temporarily disables the hook and restores the original function, the
	// hook can be enabled later using the returned Enabler.
	Disable() Enabler
	// Restore restores the original function permanently, if you want to enable the
	// hook again, you should use the Apply function later.
	Restore() error
}
```

Hook represents an applied hook, it implements Origin, Disable and Restore. The
Disable and Restore methods will change the code in the text segment, so are not
concurrent safe, need special attention.

#### func Apply

```go
func Apply(from, to interface{}) (Hook, error)
```
Apply the hook, replace "from" with "to". This function will change the code in
the text segment, so is not concurrent safe, need special attention.  Some other
goroutines may executing the code when you replacing it.

### Example
```go
package main

import (
	"fmt"
	"github.com/fengyoulin/hookingo"
)

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

func main() {
	h, err := hookingo.Apply(f1, f3)
	if err != nil {
		panic(err)
	}
	s := f2("f")
	o := h.Origin()
	if f, ok := o.(func(string) string); ok {
		s += f("F")
	} else if e, ok := o.(error); ok {
		panic(e)
	}
	e := h.Disable()
	s += f2("f")
	e.Enable()
	s += f2("F")
	err = h.Restore()
	if err != nil {
		panic(err)
	}
	s += f2("f")
	if s != "ffffFffFFFFff" {
		panic(s)
	}
}
```
Use the `//go:noinline` directive or build the example with `-gcflags='-l'` to prevent inline optimization:
```shell script
go build -gcflags '-l'
```
This example should not panic.

### v0.2.0

Because inline hook modifies instructions at the entrypoint of the target function, it introduces lots of diffculties and limitations. So a new kind of hook was introduced in the new version, which makes things simpler.

#### type HookCaller

```go
type HookCaller interface {
	// Disable disables the hooks and restores the calls to the original
	// function, the hook can be enabled later using the returned Enabler.
	Disable() Enabler
	// Count returns the total number of modified call instructions. If the
	// hooks are disabled, it returns 0.
	Count() int
}
```

HookCaller applies a group of hooks in the caller, by changing the relative
address in some call instructions. It is not concurrent safe, need special
attention.

#### func Replace

```go
func Replace(caller, old, new interface{}, length int) (h HookCaller, err error)
```
Replace the calls to "old" with "new" in the first "length" bytes of caller,
without modify any instruction in "old". When "length" is negative, it will
return at the first return instruction in the caller.

### Example
```go
package main

import (
	"fmt"
	"github.com/fengyoulin/hookingo"
)

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

func main() {
	h, err := hookingo.Replace(f2, f1, f3, -1)
	if err != nil {
		panic(err)
	}
	s := f2("f")
	s += f1("F")
	e := h.Disable()
	s += f2("f")
	e.Enable()
	s += f2("F")
	_ = h.Disable()
	s += f2("f")
	if s != "ffffFffFFFFff" {
		panic(s)
	}
}
```
Use the `//go:noinline` directive or build the example with `-gcflags='-l'` to prevent inline optimization:
```shell script
go build -gcflags '-l'
```
This example should not panic.
