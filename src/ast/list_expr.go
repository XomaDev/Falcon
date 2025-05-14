package ast

import "Falcon/sugar"

type ListExpr struct {
	Elements []Expr
}

func (l *ListExpr) String() string {
	return sugar.Format("[%]", JoinExprs(", ", l.Elements))
}

func (l *ListExpr) Blockly() Block {
	return Block{
		Type:     "lists_create_with",
		Mutation: &Mutation{ItemCount: len(l.Elements)},
		Values:   ToValues("ADD", l.Elements),
	}
}
