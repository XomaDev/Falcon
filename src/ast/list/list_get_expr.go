package list

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
)

type Get struct {
	List  blockly.Expr
	Index blockly.Expr
}

func (g *Get) String() string {
	return sugar.Format("%[%]", g.List.String(), g.Index.String())
}

func (g *Get) Blockly() blockly.Block {
	return blockly.Block{
		Type:   "lists_select_item",
		Values: blockly.MakeValues([]blockly.Expr{g.List, g.Index}, "LIST", "NUM"),
	}
}
