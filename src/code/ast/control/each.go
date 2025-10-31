package control

import (
	blockly2 "Falcon/code/ast/blockly"
	"Falcon/code/sugar"
)

type Each struct {
	IName    string
	Iterable blockly2.Expr
	Body     []blockly2.Expr
}

func (e *Each) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (e *Each) String() string {
	return sugar.Format("each % -> % {\n%}", e.IName, e.Iterable.String(), blockly2.PadBody(e.Body))
}

func (e *Each) Blockly() blockly2.Block {
	return blockly2.Block{
		Type:       "controls_forEach",
		Fields:     []blockly2.Field{{Name: "VAR", Value: e.IName}},
		Values:     []blockly2.Value{{Name: "LIST", Block: e.Iterable.Blockly()}},
		Statements: []blockly2.Statement{blockly2.CreateStatement("DO", e.Body)},
	}
}

func (e *Each) Continuous() bool {
	return false
}

func (e *Each) Consumable() bool {
	return false
}
