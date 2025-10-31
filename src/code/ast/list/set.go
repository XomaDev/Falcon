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
	yail := "(call-yail-primitive yail-list-set-item! (*list-for-runtime* "
	yail += s.List.Yail()
	yail += " "
	yail += s.Index.Yail()
	yail += " "
	yail += s.Value.Yail()
	yail += ") '(list number any) \"replace list item\")"
	return yail
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
