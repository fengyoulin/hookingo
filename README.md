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
```

#### type Enabler

```go
type Enabler interface {
	Enable()
}
```

Enabler is the interface the wraps the Enable method.

Enable enables a hook which disabled by the Disable method of the Hook
interface.

#### type Hook

```go
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
```

Hook represents a applied hook, it implements Origin, Disable and Restore. The
Disable and Restore methods will change the code in the text segment, so are not
concurrent safe, need special attention.

#### func  Apply

```go
func Apply(from, to interface{}) (Hook, error)
```
Apply the hook, replace "from" with "to". This function will change the code in
the text segment, so is not concurrent safe, need special attention.


### Example
```go
package main

import (
	"fmt"
	"github.com/fengyoulin/hookingo"
)

func say1(n string) {
	fmt.Printf("Hello, %s!\n", n)
}

func say2(n string) {
	fmt.Printf("%s，你好！\n", n)
}

func disable(s string, h hookingo.Hook) {
	defer h.Disable().Enable()
	say1(s)
}

func main() {
	s := "Golang"
	// replace say1 with say2
	h, e := hookingo.Apply(say1, say2)
	if e != nil {
		panic(e)
	}
	say1(s) // 1st, hooked
	if f, ok := h.Origin().(func(string)); ok {
		f(s) // 2nd, try to call original say1
	} else if e, ok := h.Origin().(error); ok {
		panic(e)
	}
	disable(s, h) // 3rd, temporary disable hook
	say1(s) // 4th, enabled again
	// restore say1
	e = h.Restore()
	if e != nil {
		panic(e)
	}
	say1(s) // 5th, restored
}
```
Build the example with gcflags to prevent inline optimization:
```shell script
go build -gcflags '-l'
```
The example should output:
```shell script
Golang，你好！
Hello, Golang!
Hello, Golang!
Golang，你好！
Hello, Golang!
```