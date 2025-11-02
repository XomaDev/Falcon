package fundamentals

import (
	ast2 "Falcon/lang/code/ast"
	"Falcon/lang/code/sugar"
)

type Boolean struct {
	Value bool
}

func (b *Boolean) Yail() string {
	if b.Value {
		return "#t"
	}
	return "#f"
}

func (b *Boolean) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

func (b *Boolean) Blockly() ast2.Block {
	var bText string
	if b.Value {
		bText = "TRUE"
	} else {
		bText = "FALSE"
	}
	return ast2.Block{
		Type:   "logic_boolean",
		Fields: ast2.FieldsFromMap(map[string]string{"BOOL": bText}),
	}
}

func (b *Boolean) Continuous() bool {
	return true
}

func (b *Boolean) Consumable() bool {
	return true
}

type Not struct {
	Expr ast2.Expr
}

func (n *Not) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (n *Not) String() string {
	return sugar.Format("!%", n.Expr.String())
}

func (n *Not) Blockly() ast2.Block {
	return ast2.Block{
		Type:   "logic_negate",
		Values: []ast2.Value{{Name: "BOOL", Block: n.Expr.Blockly()}},
	}
}

func (n *Not) Continuous() bool {
	return false
}

func (n *Not) Consumable() bool {
	return true
}
