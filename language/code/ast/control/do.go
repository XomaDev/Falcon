package control

import (
	"Falcon/code/ast"
)

type Do struct {
	Body   []ast.Expr
	Result ast.Expr
}

func (d *Do) Yail() string {
	yail := "(begin "
	yail += ast.PadBody(d.Body)
	yail += " "
	yail += d.Result.Yail()
	yail += ")"
	return yail
}

func (d *Do) String() string {
	return ast.JoinExprs("\n", d.Body) + "\n" + d.Result.String()
}

func (d *Do) Blockly(flags ...bool) ast.Block {
	return ast.Block{
		Type:       "controls_do_then_return",
		Statements: []ast.Statement{ast.CreateStatement("STM", d.Body)},
		Values:     []ast.Value{{Name: "VALUE", Block: d.Result.Blockly()}},
	}
}

func (d *Do) Continuous() bool {
	return false
}

func (d *Do) Consumable(flags ...bool) bool {
	return false
}

func (d *Do) Signature() []ast.Signature {
	return []ast.Signature{ast.SignVoid}
}
