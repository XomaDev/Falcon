package logic

import (
	"Falcon/ast/blockly"
)

type Bool struct {
	Value bool
}

func (b *Bool) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

func (b *Bool) Blockly() blockly.Block {
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
