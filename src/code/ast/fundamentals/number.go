package fundamentals

import (
	"Falcon/code/ast"
)

type Number struct {
	Content string
}

func (n *Number) Yail() string {
	return n.Content
}

func (n *Number) String() string {
	return n.Content
}

func (n *Number) Blockly() ast.Block {
	return ast.Block{
		Type:   "math_number",
		Fields: ast.FieldsFromMap(map[string]string{"NUM": n.Content}),
	}
}

func (n *Number) Continuous() bool {
	return true
}

func (n *Number) Consumable() bool {
	return true
}
