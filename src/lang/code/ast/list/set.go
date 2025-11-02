package list

import (
	ast2 "Falcon/lang/code/ast"
	"Falcon/lang/code/sugar"
)

type Set struct {
	List  ast2.Expr
	Index ast2.Expr
	Value ast2.Expr
}

func (s *Set) Yail() string {
	args := []ast2.Expr{s.List, s.Index, s.Value}
	return ast2.PrimitiveCall("yail-list-set-item!", "replaceListItem", args, "list number any")
}

func (s *Set) String() string {
	pFormat := "%[%] = %"
	if !s.List.Continuous() {
		pFormat = "(%)[%] = %"
	}
	return sugar.Format(pFormat, s.List.String(), s.Index.String(), s.Value.String())
}

func (s *Set) Blockly() ast2.Block {
	return ast2.Block{
		Type:   "lists_replace_item",
		Values: ast2.MakeValues([]ast2.Expr{s.List, s.Index, s.Value}, "LIST", "NUM", "ITEM"),
	}
}

func (s *Set) Continuous() bool {
	return false
}

func (s *Set) Consumable() bool {
	return false
}
