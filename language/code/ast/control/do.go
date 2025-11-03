package control

import (
	ast2 "Falcon/code/ast"
	"Falcon/code/sugar"
)

type Do struct {
	Body   []ast2.Expr
	Result ast2.Expr
}

func (d *Do) Yail() string {
	yail := "(begin "
	yail += ast2.PadBody(d.Body)
	yail += " "
	yail += d.Result.Yail()
	yail += ")"
	return yail
}

func (d *Do) String() string {
	return sugar.Format("do {\n%} -> %", ast2.PadBody(d.Body), d.Result.String())
}

func (d *Do) Blockly() ast2.Block {
	return ast2.Block{
		Type:       "controls_do_then_return",
		Statements: []ast2.Statement{ast2.CreateStatement("STM", d.Body)},
		Values:     []ast2.Value{{Name: "VALUE", Block: d.Result.Blockly()}},
	}
}

func (d *Do) Continuous() bool {
	return false
}

func (d *Do) Consumable() bool {
	return false
}
