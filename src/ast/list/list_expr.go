package list

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
)

type ListExpr struct {
	Elements []blockly.Expr
}

func (l *ListExpr) String() string {
	return sugar.Format("[%]", blockly.JoinExprs(", ", l.Elements))
}

func (l *ListExpr) Blockly() blockly.Block {
	return blockly.Block{
		Type:     "lists_create_with",
		Mutation: &blockly.Mutation{ItemCount: len(l.Elements)},
		Values:   blockly.ToValues("ADD", l.Elements),
	}
}
