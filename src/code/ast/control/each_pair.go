package control

import (
	"Falcon/code/ast"
	"Falcon/code/sugar"
)

type EachPair struct {
	KeyName   string
	ValueName string
	Iterable  ast.Expr
	Body      []ast.Expr
}

func (e *EachPair) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (e *EachPair) String() string {
	return sugar.Format("each %::% -> % {\n%}", e.KeyName, e.ValueName, e.Iterable.String(), ast.PadBody(e.Body))
}

func (e *EachPair) Blockly() ast.Block {
	return ast.Block{
		Type: "controls_for_each_dict",
		Fields: []ast.Field{
			{Name: "KEY", Value: e.KeyName},
			{Name: "VALUE", Value: e.ValueName},
		},
		Values:     []ast.Value{{Name: "DICT", Block: e.Iterable.Blockly()}},
		Statements: []ast.Statement{ast.CreateStatement("DO", e.Body)},
	}
}

func (e *EachPair) Continuous() bool {
	return false
}

func (e *EachPair) Consumable() bool {
	return false
}
