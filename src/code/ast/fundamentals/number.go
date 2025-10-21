package fundamentals

import (
	blockly2 "Falcon/code/ast/blockly"
)

type Number struct {
	Content string
}

func (n *Number) String() string {
	return n.Content
}

func (n *Number) Blockly() blockly2.Block {
	return blockly2.Block{
		Type:   "math_number",
		Fields: blockly2.FieldsFromMap(map[string]string{"NUM": n.Content}),
	}
}

func (n *Number) Continuous() bool {
	return true
}

func (n *Number) Consumable() bool {
	return true
}
