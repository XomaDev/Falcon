package control

import (
	ast2 "Falcon/code/ast"
	"Falcon/code/sugar"
)

type Each struct {
	IName    string
	Iterable ast2.Expr
	Body     []ast2.Expr
}

func (e *Each) Yail() string {
	yail := "(foreach $"
	yail += e.IName
	yail += " (begin "
	yail += ast2.PadBodyYail(e.Body)
	yail += ") "
	yail += e.Iterable.String()
	yail += ")"
	return yail
}

func (e *Each) String() string {
	return sugar.Format("each % -> % {\n%}", e.IName, e.Iterable.String(), ast2.PadBody(e.Body))
}

func (e *Each) Blockly() ast2.Block {
	return ast2.Block{
		Type:       "controls_forEach",
		Fields:     []ast2.Field{{Name: "VAR", Value: e.IName}},
		Values:     []ast2.Value{{Name: "LIST", Block: e.Iterable.Blockly()}},
		Statements: []ast2.Statement{ast2.CreateStatement("DO", e.Body)},
	}
}

func (e *Each) Continuous() bool {
	return false
}

func (e *Each) Consumable() bool {
	return false
}
