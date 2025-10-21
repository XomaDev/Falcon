package fundamentals

import (
	blockly2 "Falcon/code/ast/blockly"
)

type Text struct {
	Content string
}

func (t *Text) String() string {
	return "\"" + t.Content + "\""
}

func (t *Text) Blockly() blockly2.Block {
	return blockly2.Block{
		Type:   "text",
		Fields: blockly2.FieldsFromMap(map[string]string{"TEXT": t.Content}),
	}
}

func (t *Text) Continuous() bool {
	return true
}

func (t *Text) Consumable() bool {
	return true
}
