package logic

import (
	"Falcon/ast/blockly"
)

type BoolExpr struct {
	Value *string
}

func (b *BoolExpr) String() string {
	return *b.Value
}

func (b *BoolExpr) Blockly() blockly.Block {
	var bText string
	if *b.Value == "true" {
		bText = "TRUE"
	} else {
		bText = "FALSE"
	}
	return blockly.Block{
		Type:   "logic_boolean",
		Fields: blockly.ToFields(map[string]string{"BOOL": bText}),
	}
}
