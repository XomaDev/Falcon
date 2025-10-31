package list

import (
	blockly2 "Falcon/code/ast/blockly"
	"Falcon/code/sugar"
)

type Get struct {
	List  blockly2.Expr
	Index blockly2.Expr
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

func (g *Get) Blockly() blockly2.Block {
	return blockly2.Block{
		Type:   "lists_select_item",
		Values: blockly2.MakeValues([]blockly2.Expr{g.List, g.Index}, "LIST", "NUM"),
	}
}

func (g *Get) Continuous() bool {
	return true
}

func (g *Get) Consumable() bool {
	return true
}
