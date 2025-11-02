package fundamentals

import (
	ast2 "Falcon/lang/code/ast"
	"Falcon/lang/code/sugar"
	"strings"
)

type List struct {
	Elements []ast2.Expr
}

func (l *List) Yail() string {
	return ast2.PrimitiveCall("make-yail-list", "make a list", l.Elements, strings.Repeat("any ", len(l.Elements)))
}

func (l *List) String() string {
	return sugar.Format("[%]", ast2.JoinExprs(", ", l.Elements))
}

func (l *List) Blockly() ast2.Block {
	return ast2.Block{
		Type:     "lists_create_with",
		Mutation: &ast2.Mutation{ItemCount: len(l.Elements)},
		Values:   ast2.ValuesByPrefix("ADD", l.Elements),
	}
}

func (l *List) Continuous() bool {
	return true
}

func (l *List) Consumable() bool {
	return true
}
