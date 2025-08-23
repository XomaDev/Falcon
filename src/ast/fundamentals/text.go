package fundamentals

import "Falcon/ast/blockly"

type Text struct {
	Content string
}

func (t *Text) String() string {
	return "\"" + t.Content + "\""
}

func (t *Text) Blockly() blockly.Block {
	return blockly.Block{
		Type:   "text",
		Fields: blockly.FieldsFromMap(map[string]string{"TEXT": t.Content}),
	}
}

func (t *Text) Continuous() bool {
	return true
}

func (t *Text) Consumable() bool {
	return true
}
