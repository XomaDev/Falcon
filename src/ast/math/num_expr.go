package math

import (
	"Falcon/ast/blockly"
)

type Num struct {
	Content string
}

func (n *Num) String() string {
	return n.Content
}

func (n *Num) Blockly() blockly.Block {
	return blockly.Block{
		Type:   "math_number",
		Fields: blockly.FieldsFromMap(map[string]string{"NUM": n.Content}),
	}
}
