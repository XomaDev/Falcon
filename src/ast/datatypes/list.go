package datatypes

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
)

type List struct {
	Elements []blockly.Expr
}

func (l *List) String() string {
	return sugar.Format("[%]", blockly.JoinExprs(", ", l.Elements))
}

func (l *List) Blockly() blockly.Block {
	return blockly.Block{
		Type:       "lists_create_with",
		Mutation:   &blockly.Mutation{ItemCount: len(l.Elements)},
		Values:     blockly.ValuesByPrefix("ADD", l.Elements),
		Consumable: true,
	}
}
