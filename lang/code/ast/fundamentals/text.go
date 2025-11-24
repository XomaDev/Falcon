package fundamentals

import (
	"Falcon/code/ast"
	"strings"
)

type Text struct {
	Content string
}

func (t *Text) String() string {
	return "\"" + strings.Replace(t.Content, "\"", "\\\"", -1) + "\""
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
