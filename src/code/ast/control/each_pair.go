package control

import (
	blockly2 "Falcon/code/ast/blockly"
	"Falcon/code/sugar"
)

type EachPair struct {
	KeyName   string
	ValueName string
	Iterable  blockly2.Expr
	Body      []blockly2.Expr
}

func (e *EachPair) String() string {
	return sugar.Format("each %::% -> % {\n%}", e.KeyName, e.ValueName, e.Iterable.String(), blockly2.PadBody(e.Body))
}

func (e *EachPair) Blockly() blockly2.Block {
	return blockly2.Block{
		Type: "controls_for_each_dict",
		Fields: []blockly2.Field{
			{Name: "KEY", Value: e.KeyName},
			{Name: "VALUE", Value: e.ValueName},
		},
		Values:     []blockly2.Value{{Name: "DICT", Block: e.Iterable.Blockly()}},
		Statements: []blockly2.Statement{blockly2.CreateStatement("DO", e.Body)},
	}
}

func (e *EachPair) Continuous() bool {
	return false
}

func (e *EachPair) Consumable() bool {
	return false
}
