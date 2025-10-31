package list

import (
	"Falcon/code/ast"
	"Falcon/code/sugar"
)

type Set struct {
	List  ast.Expr
	Index ast.Expr
	Value ast.Expr
}

func (s *Set) Yail() string {
	args := []ast.Expr{s.List, s.Index, s.Value}
	return ast.PrimitiveCall("yail-list-set-item!", "replaceListItem", args, "list number any")
}

func (s *Set) String() string {
	pFormat := "%[%] = %"
	if !s.List.Continuous() {
		pFormat = "(%)[%] = %"
	}
	return sugar.Format(pFormat, s.List.String(), s.Index.String(), s.Value.String())
}

func (s *Set) Blockly() ast.Block {
	return ast.Block{
		Type:   "lists_replace_item",
		Values: ast.MakeValues([]ast.Expr{s.List, s.Index, s.Value}, "LIST", "NUM", "ITEM"),
	}
}

func (s *Set) Continuous() bool {
	return false
}

func (s *Set) Consumable() bool {
	return false
}
