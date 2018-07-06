package main

import (
	"github.com/gopherjs/vecty"
)

func main() {
	vecty.SetTitle("gox lang")
	p := &components.BodyComponent{}
	vecty.RenderBody(p)
	//js.Global.Get("console").Call("log", "dang")
}
