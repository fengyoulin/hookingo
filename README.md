# hookingo #

**hook in go** or **hooking go**

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
	o, e := hookingo.Apply(say1, say2)
	if e != nil {
		fmt.Println(e)
		return
	}
	s := "Golang"
	say1(s)
	if f, ok := o.Origin().(func(string)); ok {
		f(s)
	} else if e, ok := o.Origin().(error); ok {
		fmt.Println(e)
	}
}
```
Build:
```shell script
go build -gcflags '-N -l'
```