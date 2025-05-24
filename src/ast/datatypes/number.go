package datatypes

import "Falcon/ast/blockly"

type Number struct {
	Content string
}

func (n *Number) String() string {
	return n.Content
}

func (n *Number) Blockly() blockly.Block {
	return blockly.Block{
		Type:       "math_number",
		Fields:     blockly.FieldsFromMap(map[string]string{"NUM": n.Content}),
		Consumable: true,
	}
}
