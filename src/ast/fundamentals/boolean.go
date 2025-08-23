package fundamentals

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
)

type Boolean struct {
	Value bool
}

func (b *Boolean) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

func (b *Boolean) Blockly() blockly.Block {
	var bText string
	if b.Value {
		bText = "TRUE"
	} else {
		bText = "FALSE"
	}
	return blockly.Block{
		Type:   "logic_boolean",
		Fields: blockly.FieldsFromMap(map[string]string{"BOOL": bText}),
	}
}

func (b *Boolean) Continuous() bool {
	return true
}

func (b *Boolean) Consumable() bool {
	return true
}

type Not struct {
	Expr blockly.Expr
}

func (n *Not) String() string {
	return sugar.Format("!%", n.Expr.String())
}

func (n *Not) Blockly() blockly.Block {
	return blockly.Block{
		Type:   "logic_negate",
		Values: []blockly.Value{{Name: "BOOL", Block: n.Expr.Blockly()}},
	}
}

func (n *Not) Continuous() bool {
	return false
}

func (n *Not) Consumable() bool {
	return true
}
