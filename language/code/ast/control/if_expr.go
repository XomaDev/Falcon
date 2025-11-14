package control

import (
	"Falcon/code/ast"
	"Falcon/code/ast/fundamentals"
	"Falcon/code/sugar"
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
	format := "if (%) "
	var stringThen string
	if len(s.normalThen) == 1 {
		format += "%"
		stringThen = s.normalThen[0].String()
	} else {
		format += "{\n%}"
		stringThen = ast.PadBody(s.normalThen)
	}
	format += " else "
	var stringElse string
	if len(s.normalElse) == 1 {
		format += "%"
		stringElse = ast.JoinExprs("\n", s.normalElse)
	} else {
		format += "{\n%}"
		stringElse = ast.PadBody(s.normalElse)
	}
	return sugar.Format(format, s.condition.String(), stringThen, stringElse)
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

func (s *SimpleIf) Consumable(flags ...bool) bool {
	return !(len(flags) > 0 && flags[0])
}

func (s *SimpleIf) Signature() []ast.Signature {
	return ast.CombineSignatures(s.smartThen.Signature(), s.smartElse.Signature())
}
