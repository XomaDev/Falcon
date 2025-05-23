package list

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
)

type Set struct {
	List  blockly.Expr
	Index blockly.Expr
	Value blockly.Expr
}

func (s *Set) String() string {
	return sugar.Format("%[%] = %", s.List.String(), s.Index.String(), s.Value.String())
}

func (s *Set) Blockly() blockly.Block {
	return blockly.Block{
		Type:       "lists_replace_item",
		Values:     blockly.MakeValues([]blockly.Expr{s.List, s.Index, s.Value}, "LIST", "NUM", "ITEM"),
		Consumable: false,
	}
}
