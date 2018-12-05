package serve

import "github.com/gopherjs/gopherjs/js"

func ParseJSON(responseText string) *js.Object {
	return js.Global.Get("JSON").Call("parse", responseText)
}
