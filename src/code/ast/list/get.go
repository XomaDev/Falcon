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
	yail := "(call-yail-primitive yail-list-get-item (*list-for-runtime* "
	yail += g.List.Yail() // TODO handled differently if list is empty
	yail += " " + g.Index.Yail()
	yail += ") '(list number) \"select list item\")"
	return yail
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
