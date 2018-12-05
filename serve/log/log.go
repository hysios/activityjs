// +build !debug

package log

import "github.com/gopherjs/gopherjs/js"

func Log(pat string, args ...interface{}) {
	if js.Global != nil {
		js.Global.Get("console").Call("log", args)
	}
}
