package control

import (
	"Falcon/lang/code/ast"
	"Falcon/lang/code/ast/fundamentals"
	"Falcon/lang/code/ast/list"
	"Falcon/lang/code/ast/variables"
	"Falcon/lang/code/sugar"
)

type EachPair struct {
	KeyName   string
	ValueName string
	Iterable  ast.Expr
	Body      []ast.Expr
}

func (e *EachPair) Yail() string {
	getKey := list.Get{
		List:  &variables.Get{Global: false, Name: "item"},
		Index: &fundamentals.Number{Content: "1"},
	}
	getValue := list.Get{
		List:  &variables.Get{Global: false, Name: "item"},
		Index: &fundamentals.Number{Content: "2"},
	}
	setKey := "( " + e.KeyName + " " + getKey.Yail() + ")"
	setValue := "( " + e.ValueName + " " + getValue.Yail() + ")"

	yail := "(foreach $item "
	yail += "(let ( " + setKey + " " + setValue + ")"
	yail += ast.PadBodyYail(e.Body)
	yail += ") "
	yail += e.Iterable.Yail()
	yail += ") "
	return yail
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
