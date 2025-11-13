package control

import (
	"Falcon/code/ast"
	"Falcon/code/sugar"
)

type SimpleIf struct {
	Condition ast.Expr
	Then      ast.Expr
	Else      ast.Expr
}

func (s *SimpleIf) Yail() string {
	return "(if " + s.Condition.Yail() + " " + s.Then.Yail() + " " + s.Else.Yail() + ")"
}

func (s *SimpleIf) String() string {
	return sugar.Format("if (%) %\n\telse %", s.Condition.String(), s.Then.String(), s.Else.String())
}

func (s *SimpleIf) Blockly() ast.Block {
	return ast.Block{
		Type:   "controls_choose",
		Values: ast.MakeValues([]ast.Expr{s.Condition, s.Then, s.Else}, "TEST", "THENRETURN", "ELSERETURN"),
	}
}

func (s *SimpleIf) Continuous() bool {
	return false
}

func (s *SimpleIf) Consumable() bool {
	return true
}

func (s *SimpleIf) Signature() []ast.Signature {
	return ast.CombineSignatures(s.Then.Signature(), s.Else.Signature())
}
