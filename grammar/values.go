package grammar

import (
	"github.com/di-wu/parser/ast"
	typ "github.com/scim2/filter-parser/types"
)

func False(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(
		ast.Capture{
			Type:  typ.False,
			Value: "false",
		},
	)
}

func Null(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(
		ast.Capture{
			Type:  typ.Null,
			Value: "null",
		},
	)
}

func True(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(
		ast.Capture{
			Type:  typ.True,
			Value: "true",
		},
	)
}
