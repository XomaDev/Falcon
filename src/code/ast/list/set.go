package list

import (
	blockly2 "Falcon/code/ast/blockly"
	"Falcon/code/sugar"
)

type Set struct {
	List  blockly2.Expr
	Index blockly2.Expr
	Value blockly2.Expr
}

func (s *Set) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (s *Set) String() string {
	pFormat := "%[%] = %"
	if !s.List.Continuous() {
		pFormat = "(%)[%] = %"
	}
	return sugar.Format(pFormat, s.List.String(), s.Index.String(), s.Value.String())
}

func (s *Set) Blockly() blockly2.Block {
	return blockly2.Block{
		Type:   "lists_replace_item",
		Values: blockly2.MakeValues([]blockly2.Expr{s.List, s.Index, s.Value}, "LIST", "NUM", "ITEM"),
	}
}

func (s *Set) Continuous() bool {
	return false
}

func (s *Set) Consumable() bool {
	return false
}
