package grammar

import (
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/op"
	typ "github.com/scim2/filter-parser/v2/types"
)

// A boolean has no case sensitivity or uniqueness.
// More info: https://tools.ietf.org/html/rfc7643#section-2.3.2

func False(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(
		ast.Capture{
			Type:  typ.False,
			Value: op.And{
				op.Or{'F', 'f'},
				op.Or{'A', 'a'},
				op.Or{'L', 'l'},
				op.Or{'S', 's'},
				op.Or{'E', 'e'},
			},
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
			Value: op.And{
				op.Or{'T', 't'},
				op.Or{'R', 'r'},
				op.Or{'U', 'u'},
				op.Or{'E', 'e'},
			},
		},
	)
}
