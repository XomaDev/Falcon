package control

import (
	ast2 "Falcon/code/ast"
	"Falcon/code/sugar"
)

type For struct {
	IName string
	From  ast2.Expr
	To    ast2.Expr
	By    ast2.Expr
	Body  []ast2.Expr
}

func (f *For) Yail() string {
	yail := "(forrange "
	yail += f.IName
	yail += " (begin "
	yail += ast2.PadBodyYail(f.Body)
	yail += ") "
	yail += f.From.Yail()
	yail += " "
	yail += f.To.Yail()
	yail += " "
	yail += f.By.Yail()
	yail += ")"
	return yail
}

func (f *For) String() string {
	return sugar.Format("for %: % to % by % {\n%}",
		f.IName, f.From.String(), f.To.String(), f.By.String(), ast2.PadBody(f.Body))
}

func (f *For) Blockly() ast2.Block {
	return ast2.Block{
		Type:       "controls_forRange",
		Fields:     []ast2.Field{{Name: "VAR", Value: f.IName}},
		Values:     ast2.MakeValues([]ast2.Expr{f.From, f.To, f.By}, "START", "END", "STEP"),
		Statements: []ast2.Statement{ast2.CreateStatement("DO", f.Body)},
	}
}

func (f *For) Continuous() bool {
	return false
}

func (f *For) Consumable() bool {
	return false
}
