package fundamentals

import "Falcon/code/ast/blockly"

type Number struct {
	Content string
}

func (n *Number) Yail() string {
	return n.Content
}

func (n *Number) String() string {
	return n.Content
}

func (n *Number) Blockly() blockly.Block {
	return blockly.Block{
		Type:   "math_number",
		Fields: blockly.FieldsFromMap(map[string]string{"NUM": n.Content}),
	}
}

func (n *Number) Continuous() bool {
	return true
}

func (n *Number) Consumable() bool {
	return true
}
