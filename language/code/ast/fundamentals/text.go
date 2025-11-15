package fundamentals

import (
	"Falcon/code/ast"
)

type Text struct {
	Content string
}

// Yail TODO fix escaping problem
func (t *Text) Yail() string {
	return "\"" + t.Content + "\""
}

// String TODO fix escaping problem
func (t *Text) String() string {
	return "\"" + t.Content + "\""
}

func (t *Text) Blockly(flags ...bool) ast.Block {
	return ast.Block{
		Type:   "text",
		Fields: ast.FieldsFromMap(map[string]string{"TEXT": t.Content}),
	}
}

func (t *Text) Continuous() bool {
	return true
}

func (t *Text) Consumable(flags ...bool) bool {
	return true
}

func (t *Text) Signature() []ast.Signature {
	return []ast.Signature{ast.SignText}
}
