package fundamentals

import (
	"Falcon/code/ast/blockly"
	"Falcon/code/sugar"
	"strings"
)

type List struct {
	Elements []blockly.Expr
}

func (l *List) Yail() string {
	yail := "(call-yail-primitive make-yail-list (*list-for-runtime* "
	yail += blockly.JoinYailExprs(" ", l.Elements)
	yail += ") '("
	yail += strings.Repeat("any ", len(l.Elements))
	yail += ") \"make a list\")"
	return yail
}

func (l *List) String() string {
	return sugar.Format("[%]", blockly.JoinExprs(", ", l.Elements))
}

func (l *List) Blockly() blockly.Block {
	return blockly.Block{
		Type:     "lists_create_with",
		Mutation: &blockly.Mutation{ItemCount: len(l.Elements)},
		Values:   blockly.ValuesByPrefix("ADD", l.Elements),
	}
}

func (l *List) Continuous() bool {
	return true
}

func (l *List) Consumable() bool {
	return true
}
