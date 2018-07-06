package transformer

import (
	"fmt"
	"os"
	"strings"

	"github.com/gowasm/gox/ast"
	"github.com/gowasm/gox/parser"
	"github.com/gowasm/gox/printer"
	"github.com/gowasm/gox/token"
)

func rename() {
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, "../goxtests/args_and_more.gox", nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Can't parse file", err)
	}

	file.Name.Name = "what" // change package name

	r := &Renamer{"Foo", "Bar"}
	ast.Walk(r, file)
	cfg := &printer.Config{Mode: printer.GoxToGo}
	cfg.Fprint(os.Stdout, fs, file)
}

type Renamer struct {
	find    string
	replace string
}

func (r *Renamer) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.FuncDecl:
			if n.Recv != nil && n.Recv.List != nil && len(n.Recv.List) > 0 {
				field := n.Recv.List[0]
				typ := field.Type.(*ast.StarExpr).X.(*ast.Ident).Name
				if typ == r.find {
					field.Names[0].Name = strings.ToLower(r.replace[0:1])
				}
			}
		case *ast.Ident:
			if n.Name == r.find {
				n.Name = r.replace
			}
		}
	}
	return r
}
