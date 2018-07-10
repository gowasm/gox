package printer

import (
	"fmt"
	"strconv"
	"strings"

	"unicode"

	"github.com/gowasm/gox/ast"
	"github.com/gowasm/gox/token"
)

var elemNameMap = map[string]string{
	"a":          "Anchor",
	"abbr":       "Abbreviation",
	"b":          "Bold",
	"bdi":        "BidirectionalIsolation",
	"bdo":        "BidirectionalOverride",
	"blockquote": "BlockQuote",
	"br":         "Break",
	"cite":       "Citation",
	"col":        "Column",
	"colgroup":   "ColumnGroup",
	"datalist":   "DataList",
	"dd":         "Description",
	"del":        "DeletedText",
	"dfn":        "Definition",
	"dl":         "DescriptionList",
	"dt":         "DefinitionTerm",
	"em":         "Emphasis",
	"fieldset":   "FieldSet",
	"figcaption": "FigureCaption",
	"h1":         "Heading1",
	"h2":         "Heading2",
	"h3":         "Heading3",
	"h4":         "Heading4",
	"h5":         "Heading5",
	"h6":         "Heading6",
	"hgroup":     "HeadingsGroup",
	"hr":         "HorizontalRule",
	"i":          "Italic",
	"iframe":     "InlineFrame",
	"img":        "Image",
	"ins":        "InsertedText",
	"kbd":        "KeyboardInput",
	"li":         "ListItem",
	"menuitem":   "MenuItem",
	"nav":        "Navigation",
	"noframes":   "NoFrames",
	"noscript":   "NoScript",
	"ol":         "OrderedList",
	"optgroup":   "OptionsGroup",
	"p":          "Paragraph",
	"param":      "Parameter",
	"pre":        "Preformatted",
	"q":          "Quote",
	"rp":         "RubyParenthesis",
	"rt":         "RubyText",
	"rtc":        "RubyTextContainer",
	"s":          "Strikethrough",
	"samp":       "Sample",
	"sub":        "Subscript",
	"sup":        "Superscript",
	"tbody":      "TableBody",
	"textarea":   "TextArea",
	"td":         "TableData",
	"tfoot":      "TableFoot",
	"th":         "TableHeader",
	"thead":      "TableHead",
	"tr":         "TableRow",
	"u":          "Underline",
	"ul":         "UnorderedList",
	"var":        "Variable",
	"wbr":        "WordBreakOpportunity",
}

// Map html-style to actual js event names
var eventMap = map[string]string{
	"onAbort":          "abort",
	"onCancel":         "cancel",
	"onCanPlay":        "canplay",
	"onCanPlaythrough": "canplaythrough",
	"onChange":         "change",
	"onClick":          "click",
	"onCueChange":      "cuechange",
	"onDblClick":       "dblclick",
	"onDurationChange": "durationchange",
	"onEmptied":        "emptied",
	"onEnded":          "ended",
	"onInput":          "input",
	"onInvalid":        "invalid",
	"onKeyDown":        "keydown",
	"onKeyPress":       "keypress",
	"onKeyUp":          "keyup",
	"onLoadedData":     "loadeddata",
	"onLoadedMetadata": "loadedmetadata",
	"onLoadStart":      "loadstart",
	"onMouseDown":      "mousedown",
	"onMouseEnter":     "mouseenter",
	"onMouseleave":     "mouseleave",
	"onMouseMove":      "mousemove",
	"onMouseOut":       "mouseout",
	"onMouseOver":      "mouseover",
	"onMouseUp":        "mouseup",
	"onMouseWheel":     "mousewheel",
	"onPause":          "pause",
	"onPlay":           "play",
	"onPlaying":        "playing",
	"onProgress":       "progress",
	"onRateChange":     "ratechange",
	"onReset":          "reset",
	"onSeeked":         "seeked",
	"onSeeking":        "seeking",
	"onSelect":         "select",
	"onShow":           "show",
	"onStalled":        "stalled",
	"onSubmit":         "submit",
	"onSuspend":        "suspend",
	"onTimeUpdate":     "timeupdate",
	"onToggle":         "toggle",
	"onVolumeChange":   "volumechange",
	"onWaiting":        "waiting",
}

var attrMap = map[string]string{
	"autofocus": "autofocus",
	"checked":   "checked",
	//	"class":       "className",
	"for":         "htmlFor",
	"href":        "href",
	"id":          "id",
	"placeholder": "placeholder",
	"src":         "src",
	"type":        "type",
	"value":       "value",
}

func goxToVecty(gox *ast.GoxExpr) ast.Expr {
	fmt.Println("GOX:", gox.TagName, gox.Ctag, gox.Otag)
	isComponent := unicode.IsUpper(rune(gox.TagName.Name[0]))

	if isComponent {
		return newComponent(gox)
	} else {
		args := []ast.Expr{
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote(gox.TagName.Name),
			}}

		// Add the attributes
		args = append(args, mapProps(gox.Attrs)...)

		// Add the contents of the tag
		for _, expr := range gox.X {
			switch expr := expr.(type) {
			// just throw the Go stuff right in there
			case *ast.GoExpr:
				e := newCallExpr(
					newSelectorExpr("vecty", "Text"),
					[]ast.Expr{expr},
				)
				args = append(args, e)

			case *ast.BareWordsExpr:
				if len(strings.TrimSpace(expr.Value)) == 0 {
					continue
				}
				e := newCallExpr(
					newSelectorExpr("vecty", "Text"),
					[]ast.Expr{expr},
				)
				args = append(args, e)
			default:
				args = append(args, expr)
			}
		}
		fmt.Println(gox.TagName, gox.Ctag, gox.Otag)
		for _, arg := range args {
			fmt.Println(arg.Pos(), arg.End())
		}
		fmt.Println(args)

		var elemname string
		if tg, ok := elemNameMap[gox.TagName.Name]; ok {
			elemname = tg
		} else {
			elemname = strings.Title(gox.TagName.Name)
		}
		return newCallExpr(
			// check the elemMap here TODO
			newSelectorExpr("elem", elemname),
			args[1:], // don't include the tag itself as an arg (or you get vecty.Div("div"))
			//nil,
		)
	}
}

func newSelectorExpr(x, sel string) *ast.SelectorExpr {
	return &ast.SelectorExpr{
		X:   ast.NewIdent(x),
		Sel: ast.NewIdent(sel)}
}

func newCallExpr(fun ast.Expr, args []ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{
		Fun:      fun,
		Args:     args,
		Ellipsis: token.NoPos, Lparen: token.NoPos, Rparen: token.NoPos}
}

func newComponent(gox *ast.GoxExpr) *ast.UnaryExpr {
	var args []ast.Expr
	for _, attr := range gox.Attrs {
		if attr.Rhs == nil { // default to true like JSX
			attr.Rhs = ast.NewIdent("true")
		}
		expr := &ast.KeyValueExpr{
			Key:   ast.NewIdent(attr.Lhs.Name),
			Colon: token.NoPos,
			Value: attr.Rhs,
		}

		args = append(args, expr)
	}

	return &ast.UnaryExpr{
		OpPos: token.NoPos,
		Op:    token.AND,
		X: &ast.CompositeLit{
			Type:   ast.NewIdent(gox.TagName.Name),
			Lbrace: token.NoPos,
			Elts:   args,
			Rbrace: token.NoPos,
		},
	}
}

func mapProps(goxAttrs []*ast.GoxAttrStmt) []ast.Expr {
	var mapped = []ast.Expr{}
	for _, attr := range goxAttrs {
		// set default of Rhs to true if none provided
		if attr.Rhs == nil { // default to true like JSX
			attr.Rhs = ast.NewIdent("true")
		}

		var expr ast.Expr
		fmt.Println("LHS attr", attr.Lhs.Name)
		// if prop is an event listener (e.g. "onClick")
		if _, ok := eventMap[attr.Lhs.Name]; ok {
			expr = newEventListener(attr)
		} else if attr.Lhs.Name == "class" {
			// if it's a class statement
			expr = newCallExpr(
				newSelectorExpr("vecty", "Markup"),
				[]ast.Expr{
					newCallExpr(
						newSelectorExpr("vecty", "Class"),
						[]ast.Expr{
							attr.Rhs,
						}),
				},
			)

		} else if mappedName, ok := attrMap[attr.Lhs.Name]; ok {
			// if it's a vecty controlled prop
			expr = newCallExpr(
				newSelectorExpr("prop", mappedName),
				[]ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: strconv.Quote(mappedName)},
					attr.Rhs,
				},
			)
		} else {
			// if prop is a normal attribute
			expr = newCallExpr(
				newSelectorExpr("vecty", "Attribute"),
				[]ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: strconv.Quote(attr.Lhs.Name)},
					attr.Rhs,
				},
			)
		}

		mapped = append(mapped, expr)
	}

	return mapped
}

func newEventListener(goxAttr *ast.GoxAttrStmt) ast.Expr {
	return &ast.UnaryExpr{
		OpPos: token.NoPos,
		Op:    token.AND,
		X: &ast.CompositeLit{
			Type:   newSelectorExpr("vecty", "EventListener"),
			Lbrace: token.NoPos,
			Elts: []ast.Expr{
				&ast.KeyValueExpr{
					Key: ast.NewIdent("Name"),
					Value: &ast.BasicLit{
						Kind:  token.STRING,
						Value: strconv.Quote(eventMap[goxAttr.Lhs.Name]),
					},
				},
				&ast.KeyValueExpr{
					Key:   ast.NewIdent("Listener"),
					Value: goxAttr.Rhs,
				},
			},
			Rbrace: token.NoPos,
		},
	}
}
