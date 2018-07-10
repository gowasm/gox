package main

import (
	"github.com/gowasm/vecty"

	"github.com/gowasm/vecty/elem"
)

func getHTML2() vecty.ComponentOrHTML {
	return elem.Span(&MyComponent{Parameter1: "Hello World"})
}

type MyComponent struct {
	vecty.Core
	Parameter1 string
}

func (c *MyComponent) Render() vecty.ComponentOrHTML {
	return elem.Div(vecty.Text(c.Parameter1))
}

//func main()	{}
