package text

import (
	"Falcon/ast/blockly"
)

type Expr struct {
	Content string
}

func (t *Expr) String() string {
	return "\"" + t.Content + "\""
}

func (t *Expr) Blockly() blockly.Block {
	return blockly.Block{
		Type:       "text",
		Fields:     blockly.FieldsFromMap(map[string]string{"TEXT": t.Content}),
		Consumable: false,
	}
}
