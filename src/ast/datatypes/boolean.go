package datatypes

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
		Type:       "logic_boolean",
		Fields:     blockly.FieldsFromMap(map[string]string{"BOOL": bText}),
		Consumable: true,
	}
}

type Not struct {
	Expr blockly.Expr
}

func (n *Not) String() string {
	return sugar.Format("!%", n.Expr.String())
}

func (n *Not) Blockly() blockly.Block {
	return blockly.Block{
		Type:       "logic_negate",
		Values:     []blockly.Value{{Name: "BOOL", Block: n.Expr.Blockly()}},
		Consumable: true,
	}
}
