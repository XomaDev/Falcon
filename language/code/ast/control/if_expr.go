package control

import (
	"Falcon/code/ast"
	"Falcon/code/ast/fundamentals"
	"Falcon/code/sugar"
	"strings"
)

type SimpleIf struct {
	condition ast.Expr

	smartThen ast.Expr
	smartElse ast.Expr

	normalThen []ast.Expr
	normalElse []ast.Expr
}

func MakeSimpleIf(condition ast.Expr, then []ast.Expr, elze []ast.Expr) *SimpleIf {
	return &SimpleIf{
		condition:  condition,
		smartThen:  &fundamentals.SmartBody{Body: then},
		smartElse:  &fundamentals.SmartBody{Body: elze},
		normalThen: then,
		normalElse: elze,
	}
}

func (s *SimpleIf) Yail() string {
	return "(if " + s.condition.Yail() + " " + s.smartThen.Yail() + " " + s.smartElse.Yail() + ")"
}

func (s *SimpleIf) String() string {
	format := "if (%) %\nelse %"
	if strings.Contains(s.smartThen.String(), "\n") {
		format = "if (%) % else %"
	}
	return sugar.Format(format, s.condition.String(), s.smartThen.String(), s.smartElse.String())
}

func (s *SimpleIf) Blockly(flags ...bool) ast.Block {
	if len(flags) > 0 && flags[0] {
		// Blockly expects a statement here, we have to mutate
		fullIf := If{
			Conditions: []ast.Expr{s.condition},
			Bodies:     [][]ast.Expr{s.normalThen},
			ElseBody:   s.normalElse,
		}
		return fullIf.Blockly()
	}
	return ast.Block{
		Type:   "controls_choose",
		Values: ast.MakeValues([]ast.Expr{s.condition, s.smartThen, s.smartElse}, "TEST", "THENRETURN", "ELSERETURN"),
	}
}

func (s *SimpleIf) Continuous() bool {
	return false
}

func (s *SimpleIf) Consumable() bool {
	return true
}

func (s *SimpleIf) Signature() []ast.Signature {
	return ast.CombineSignatures(s.smartThen.Signature(), s.smartElse.Signature())
}
