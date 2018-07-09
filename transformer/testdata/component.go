package main

import (
	"github.com/gopherjs/vecty"
)

func getHTML2() vecty.HTML {
	return elem.span("span", &MyComponent{Parameter1: "Hello World"})
}

type MyComponent struct {
	vecty.Core
	Parameter1	string
}

func (c *MyComponent) Render() *vecty.HTML {
	return elem.div("div", c.Parameter1)
}

func main()	{}
