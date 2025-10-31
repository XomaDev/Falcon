package control

import (
	"Falcon/code/ast"
	"Falcon/code/sugar"
)

type Do struct {
	Body   []ast.Expr
	Result ast.Expr
}

func (d *Do) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (d *Do) String() string {
	return sugar.Format("do {\n%} -> %", ast.PadBody(d.Body), d.Result.String())
}

func (d *Do) Blockly() ast.Block {
	return ast.Block{
		Type:       "controls_do_then_return",
		Statements: []ast.Statement{ast.CreateStatement("STM", d.Body)},
		Values:     []ast.Value{{Name: "VALUE", Block: d.Result.Blockly()}},
	}
}

func (d *Do) Continuous() bool {
	return false
}

func (d *Do) Consumable() bool {
	return false
}
