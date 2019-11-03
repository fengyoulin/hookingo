# hookingo #

The name of this project reads like **hook in go** or **hooking go**, it is just an implementation of the "function hook" technology in golang.

Usage:
```shell script
go get github.com/fengyoulin/hookingo
```
Example:
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

func main() {
	// replace say1 with say2
	o, e := hookingo.Apply(say1, say2)
	if e != nil {
		fmt.Println(e)
		return
	}
	s := "Golang"
	// call say1
	say1(s)
	// try to call original say1
	if f, ok := o.Origin().(func(string)); ok {
		f(s)
	} else if e, ok := o.Origin().(error); ok {
		fmt.Println(e)
	}
	// restore say1
	e = o.Restore()
	say1(s)
}
```
Build the example with gcflags to prevent inline optimization:
```shell script
go build -gcflags '-N -l'
```
The example should output:
```shell script
Golang，你好！
Hello, Golang!
Hello, Golang!
```