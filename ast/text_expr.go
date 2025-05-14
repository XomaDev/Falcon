package ast

type TextExpr struct {
	Content *string
}

func (t *TextExpr) String() string {
	return "\"" + *t.Content + "\""
}

func (t *TextExpr) Blockly() Block {
	return Block{
		Type:   "text",
		Fields: ToFields(map[string]string{"TEXT": *t.Content}),
	}
}
