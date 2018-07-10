package main

import (
	"github.com/gopherjs/vecty"
	"github.com/gowasm/vecty/elem"
	//"github.com/gowasm/vecty/props"
)

func getHTML() vecty.ComponentOrHTML {
	return elem.Div(vecty.Markup(style.Blue("blue")))
}

func main() {}
