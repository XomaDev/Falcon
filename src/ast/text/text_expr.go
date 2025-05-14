package text

import (
	"Falcon/ast/blockly"
)

type TextExpr struct {
	Content *string
}

func (t *TextExpr) String() string {
	return "\"" + *t.Content + "\""
}

func (t *TextExpr) Blockly() blockly.Block {
	return blockly.Block{
		Type:   "text",
		Fields: blockly.ToFields(map[string]string{"TEXT": *t.Content}),
	}
}
