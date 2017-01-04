package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/loader"
)

type RequestWithContextWarning struct {
	Pos  token.Position
	Name string
}

func (warning RequestWithContextWarning) String() string {
	return fmt.Sprintf("net/http request.WithContext() called without lvalue at %s: %q\n", warning.Pos, warning.Name)
}

func main() {
	tests := flag.Bool("t", false, "also check dependencies of test files")
	flag.Parse()
	var conf loader.Config
	for _, p := range flag.Args() {
		if *tests {
			conf.ImportWithTests(p)
		} else {
			conf.Import(p)
		}
	}
	lprog, err := conf.Load()
	if err != nil {
		log.Fatalf("Load error: %v", err)
	}
	for _, info := range lprog.InitialPackages() {
		warnings := CheckRequestWithContext(lprog.Fset, &info.Info, info.Files)
		for _, warning := range warnings {
			fmt.Printf("%s\n", warning.String())
		}
	}
}

func CheckRequestWithContext(fset *token.FileSet, info *types.Info, files []*ast.File) []RequestWithContextWarning {
	warnings := make([]RequestWithContextWarning, 0)
	for _, file := range files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.CallExpr:
				switch selType := n.Fun.(type) {
				case *ast.SelectorExpr:
					receiver := selType.X
					receiverType := info.Types[receiver].Type
					selector := selType.Sel
					// XXX: use real ast.Objs here -- https://godoc.org/go/types#LookupFieldOrMethod
					if selector.String() == "WithContext" && receiverType != nil && receiverType.String() == "*net/http.Request" {
						// XXX: better way to get the parent node?
						// or perhaps there is a better way to determine if this CallExpr is part of an assignment
						// info.ObjectOf(selector).Parent() is a *Scope, maybe we should iterate over that
						inAssignment := false
						path, _ := astutil.PathEnclosingInterval(file, n.Pos(), n.End())
						for _, part := range path {
							switch parentType := part.(type) {
							case *ast.AssignStmt:
								for _, expr := range parentType.Rhs {
									if expr == n {
										inAssignment = true
									}
								}
							}
						}
						if !inAssignment {
							// XXX: can we print out the entire line of the file?
							warnings = append(warnings, RequestWithContextWarning{Pos: fset.Position(n.Pos()), Name: ""})
						}
					}
				}
			}
			return true
		})
	}
	return warnings
}
