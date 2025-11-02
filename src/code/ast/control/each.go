package control

import (
	"Falcon/code/ast"
	"Falcon/code/sugar"
)

type Each struct {
	IName    string
	Iterable ast.Expr
	Body     []ast.Expr
}

func (e *Each) Yail() string {
	yail := "(foreach $"
	yail += e.IName
	yail += " (begin "
	yail += ast.PadBodyYail(e.Body)
	yail += ") "
	yail += e.Iterable.String()
	yail += ")"
	return yail
}

func (e *Each) String() string {
	return sugar.Format("each % -> % {\n%}", e.IName, e.Iterable.String(), ast.PadBody(e.Body))
}

func (e *Each) Blockly() ast.Block {
	return ast.Block{
		Type:       "controls_forEach",
		Fields:     []ast.Field{{Name: "VAR", Value: e.IName}},
		Values:     []ast.Value{{Name: "LIST", Block: e.Iterable.Blockly()}},
		Statements: []ast.Statement{ast.CreateStatement("DO", e.Body)},
	}
}

func (e *Each) Continuous() bool {
	return false
}

func (e *Each) Consumable() bool {
	return false
}
