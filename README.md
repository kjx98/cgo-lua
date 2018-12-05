# cgo-lua
A Go(lang) wrapper for LuaJit-2.1(Lua JIT Library) which is often used for embed script support.
Forked from [RyouZhang] (https://github.com/RyouZhang/go-lua)

[![Build Status](https://travis-ci.org/kjx98/cgo-lua.svg?branch=master)](https://travis-ci.org/kjx98/cgo-lua)
[![GoDoc](https://godoc.org/github.com/kjx98/cgo-lua?status.svg)](https://godoc.org/github.com/kjx98/cgo-lua)

To use the library you need LuaJit 2.1 installed.

## Example
```go
package main

import (
	"fmt"
	"github.com/kjx98/cgo-lua"
)


func main() {
	res, err := lua.Call("script.lua", "test_args", 69)
	if err != nil {
		fmt.Println("test_args", err)
	} else {
		fmt.Println(res)
	}
}
```

## Installing

Install the dependencies then run

```
$ go install github.com/kjx98/cgo-lua
```

### Dependencies

To use luajit for go, you need to have the
[Luajit-2.1](http://luajit.org) already installed:

##### Mac OS X

```
$ brew install luajit
```

##### Linux

Install from your package manager or install from source.

On Debian/Ubuntu/Fedora Linux `luajit-2.1` is available from the distro.

To compile first download [Luajit-2.1](http://luajit.org/download/LuaJIT-2.1.0-beta3.tar.gz) and:
```
$ untar and cd
$ make
$ sudo make install
```


## License
Copyright (c) 2018 Jesse Kuang

