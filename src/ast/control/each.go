package control

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
)

type Each struct {
	IName    string
	Iterable blockly.Expr
	Body     []blockly.Expr
}

func (e *Each) String() string {
	return sugar.Format("each % -> % {\n%}", e.IName, e.Iterable.String(), blockly.PadBody(e.Body))
}

func (e *Each) Blockly() blockly.Block {
	return blockly.Block{
		Type:       "controls_forEach",
		Fields:     []blockly.Field{{Name: "VAR", Value: e.IName}},
		Values:     []blockly.Value{{Name: "LIST", Block: e.Iterable.Blockly()}},
		Statements: []blockly.Statement{blockly.CreateStatement("DO", e.Body)},
	}
}

func (e *Each) Continuous() bool {
	return false
}

func (e *Each) Consumable() bool {
	return false
}
