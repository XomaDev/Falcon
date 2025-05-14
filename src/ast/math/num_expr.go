package math

import (
	"Falcon/ast/blockly"
)

type NumExpr struct {
	Content *string
}

func (n *NumExpr) String() string {
	return *n.Content
}

func (n *NumExpr) Blockly() blockly.Block {
	return blockly.Block{
		Type:   "math_number",
		Fields: blockly.ToFields(map[string]string{"NUM": *n.Content}),
	}
}
