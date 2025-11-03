package control

import (
	ast2 "Falcon/code/ast"
	"Falcon/code/ast/fundamentals"
	"Falcon/code/ast/list"
	"Falcon/code/ast/variables"
	"Falcon/code/sugar"
)

type EachPair struct {
	KeyName   string
	ValueName string
	Iterable  ast2.Expr
	Body      []ast2.Expr
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
	yail += ast2.PadBodyYail(e.Body)
	yail += ") "
	yail += e.Iterable.Yail()
	yail += ") "
	return yail
}

func (e *EachPair) String() string {
	return sugar.Format("each %::% -> % {\n%}", e.KeyName, e.ValueName, e.Iterable.String(), ast2.PadBody(e.Body))
}

func (e *EachPair) Blockly() ast2.Block {
	return ast2.Block{
		Type: "controls_for_each_dict",
		Fields: []ast2.Field{
			{Name: "KEY", Value: e.KeyName},
			{Name: "VALUE", Value: e.ValueName},
		},
		Values:     []ast2.Value{{Name: "DICT", Block: e.Iterable.Blockly()}},
		Statements: []ast2.Statement{ast2.CreateStatement("DO", e.Body)},
	}
}

func (e *EachPair) Continuous() bool {
	return false
}

func (e *EachPair) Consumable() bool {
	return false
}
