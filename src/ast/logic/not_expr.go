package logic

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
)

type Not struct {
	Expr blockly.Expr
}

func (n *Not) String() string {
	return sugar.Format("!%", n.Expr.String())
}

func (n *Not) Blockly() blockly.Block {
	return blockly.Block{
		Type:   "logic_negate",
		Values: []blockly.Value{{Name: "BOOL", Block: n.Expr.Blockly()}},
	}
}
