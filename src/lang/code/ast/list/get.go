package list

import (
	ast2 "Falcon/lang/code/ast"
	"Falcon/lang/code/sugar"
)

type Get struct {
	List  ast2.Expr
	Index ast2.Expr
}

func (g *Get) Yail() string {
	args := []ast2.Expr{g.List, g.Index}
	return ast2.PrimitiveCall("yail-list-get-item", "setListItem", args, "list number")
}

func (g *Get) String() string {
	pFormat := "%[%]"
	if !g.List.Continuous() {
		pFormat = "(%)[%]"
	}
	return sugar.Format(pFormat, g.List.String(), g.Index.String())
}

func (g *Get) Blockly() ast2.Block {
	return ast2.Block{
		Type:   "lists_select_item",
		Values: ast2.MakeValues([]ast2.Expr{g.List, g.Index}, "LIST", "NUM"),
	}
}

func (g *Get) Continuous() bool {
	return true
}

func (g *Get) Consumable() bool {
	return true
}
