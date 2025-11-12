package fundamentals

import (
	"Falcon/code/ast"
	"Falcon/code/sugar"
	"strings"
)

type List struct {
	Elements []ast.Expr
}

func (l *List) Yail() string {
	return ast.PrimitiveCall("make-yail-list", "make a list", l.Elements, strings.Repeat("any ", len(l.Elements)))
}

func (l *List) String() string {
	return sugar.Format("[%]", ast.JoinExprs(", ", l.Elements))
}

func (l *List) Blockly() ast.Block {
	return ast.Block{
		Type:     "lists_create_with",
		Mutation: &ast.Mutation{ItemCount: len(l.Elements)},
		Values:   ast.ValuesByPrefix("ADD", l.Elements),
	}
}

func (l *List) Continuous() bool {
	return true
}

func (l *List) Consumable() bool {
	return true
}

func (l *List) Signature() []ast.Signature {
	return []ast.Signature{ast.SignList}
}
