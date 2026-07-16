package grammar

import (
	"github.com/q-uint/parser/ast"
	"github.com/q-uint/parser/op"
	"github.com/scim2/filter-parser/v2/internal/types"
)

func Path(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type:        typ.Path,
		TypeStrings: typ.Stringer,
		Value: op.Or{
			op.And{
				ValuePath,
				op.Optional(SubAttr),
			},
			AttrPath,
		},
	})
}
