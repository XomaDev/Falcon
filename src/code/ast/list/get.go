package list

import (
	"Falcon/code/ast"
	"Falcon/code/sugar"
)

type Get struct {
	List  ast.Expr
	Index ast.Expr
}

func (g *Get) Yail() string {
	args := []ast.Expr{g.List, g.Index}
	return ast.PrimitiveCall("yail-list-get-item", "setListItem", args, "list number")
}

func (g *Get) String() string {
	pFormat := "%[%]"
	if !g.List.Continuous() {
		pFormat = "(%)[%]"
	}
	return sugar.Format(pFormat, g.List.String(), g.Index.String())
}

func (g *Get) Blockly() ast.Block {
	return ast.Block{
		Type:   "lists_select_item",
		Values: ast.MakeValues([]ast.Expr{g.List, g.Index}, "LIST", "NUM"),
	}
}

func (g *Get) Continuous() bool {
	return true
}

func (g *Get) Consumable() bool {
	return true
}
