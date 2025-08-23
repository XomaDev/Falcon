package common

import "Falcon/ast/blockly"

type EmptySocket struct{}

func (e *EmptySocket) String() string {
	return "undefined"
}

func (e *EmptySocket) Blockly() blockly.Block {
	return blockly.Block{
		Type:   "math_number",
		Fields: blockly.FieldsFromMap(map[string]string{"NUM": "0"}),
	}
}

func (e *EmptySocket) Continuous() bool {
	return true
}

func (e *EmptySocket) Consumable() bool {
	return false
}
