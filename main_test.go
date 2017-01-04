package main

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type CheckRequestWithContextSuite struct{}

var _ = Suite(&CheckRequestWithContextSuite{})

func ParseAndReturnWarnings(c *C, file string) []RequestWithContextWarning {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "file.go", file, 0)
	c.Assert(err, IsNil)

	conf := types.Config{Importer: importer.Default()}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}
	_, err = conf.Check("bad", fset, []*ast.File{f}, info)
	c.Assert(err, IsNil)

	return CheckRequestWithContext(fset, info, []*ast.File{f})
}

func (*CheckRequestWithContextSuite) TestCheckRequestWithContextBad(c *C) {
	var bad = `
package main

import (
	"bytes"
	"net/http"
)

func bad() {
	url := "http://example.com"
	payload := "this should be flagged by the context checker"
	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(payload))
	req.WithContext(nil)
}
`
	warnings := ParseAndReturnWarnings(c, bad)
	c.Assert(warnings, NotNil)
	c.Assert(warnings, DeepEquals,
		[]RequestWithContextWarning{{
			Pos: token.Position{
				Filename: "file.go",
				Offset:   223,
				Line:     13,
				Column:   2,
			},
			Name: "",
		}})
}

func (*CheckRequestWithContextSuite) TestCheckRequestWithContextGood(c *C) {
	var good = `
package main

import (
	"bytes"
	"net/http"
)

func good() {
	url := "http://example.com"
	payload := "this should be flagged by the context checker"
	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(payload))
	req = req.WithContext(nil)
}
`
	warnings := ParseAndReturnWarnings(c, good)
	c.Assert(warnings, NotNil)
	c.Assert(warnings, HasLen, 0)
}
