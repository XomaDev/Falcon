package list

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
)

type Expr struct {
	Elements []blockly.Expr
}

func (l *Expr) String() string {
	return sugar.Format("[%]", blockly.JoinExprs(", ", l.Elements))
}

func (l *Expr) Blockly() blockly.Block {
	return blockly.Block{
		Type:       "lists_create_with",
		Mutation:   &blockly.Mutation{ItemCount: len(l.Elements)},
		Values:     blockly.ValuesByPrefix("ADD", l.Elements),
		Consumable: true,
	}
}
