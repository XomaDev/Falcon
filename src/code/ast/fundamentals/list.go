package fundamentals

import (
	blockly2 "Falcon/code/ast/blockly"
	"Falcon/code/sugar"
)

type List struct {
	Elements []blockly2.Expr
}

func (l *List) String() string {
	return sugar.Format("[%]", blockly2.JoinExprs(", ", l.Elements))
}

func (l *List) Blockly() blockly2.Block {
	return blockly2.Block{
		Type:     "lists_create_with",
		Mutation: &blockly2.Mutation{ItemCount: len(l.Elements)},
		Values:   blockly2.ValuesByPrefix("ADD", l.Elements),
	}
}

func (l *List) Continuous() bool {
	return true
}

func (l *List) Consumable() bool {
	return true
}
