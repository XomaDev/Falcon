package logic

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
)

type NotExpr struct {
	Expr blockly.Expr
}

func (n *NotExpr) String() string {
	return sugar.Format("!%", n.Expr.String())
}

func (n *NotExpr) Blockly() blockly.Block {
	return blockly.Block{
		Type:   "logic_negate",
		Values: []blockly.Value{{Name: "BOOL", Block: n.Expr.Blockly()}},
	}
}
