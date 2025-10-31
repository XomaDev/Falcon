package common

import (
	blockly2 "Falcon/code/ast/blockly"
)

type EmptySocket struct{}

func (e *EmptySocket) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (e *EmptySocket) String() string {
	return "undefined"
}

func (e *EmptySocket) Blockly() blockly2.Block {
	return blockly2.Block{
		Type:   "math_number",
		Fields: blockly2.FieldsFromMap(map[string]string{"NUM": "0"}),
	}
}

func (e *EmptySocket) Continuous() bool {
	return true
}

func (e *EmptySocket) Consumable() bool {
	return false
}
