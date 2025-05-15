package control

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
)

type EachPair struct {
	KeyName   string
	ValueName string
	Iterable  blockly.Expr
	Body      []blockly.Expr
}

func (e *EachPair) String() string {
	return sugar.Format("each %::% -> % {\n%}", e.KeyName, e.ValueName, e.Iterable.String(), blockly.PadBody(e.Body))
}

func (e *EachPair) Blockly() blockly.Block {
	return blockly.Block{
		Type: "controls_for_each_dict",
		Fields: []blockly.Field{
			{Name: "KEY", Value: e.KeyName},
			{Name: "VALUE", Value: e.ValueName},
		},
		Values:     []blockly.Value{{Name: "DICT", Block: e.Iterable.Blockly()}},
		Statements: []blockly.Statement{blockly.CreateStatement("DO", e.Body)},
	}
}
