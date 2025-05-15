package control

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
)

type SimpleIf struct {
	Condition blockly.Expr
	Then      blockly.Expr
	Else      blockly.Expr
}

func (s *SimpleIf) String() string {
	return sugar.Format("if (%) % else %", s.Condition.String(), s.Then.String(), s.Else.String())
}

func (s *SimpleIf) Blockly() blockly.Block {
	return blockly.Block{
		Type:   "controls_choose",
		Values: blockly.MakeValues([]blockly.Expr{s.Condition, s.Then, s.Else}, "TEST", "THENRETURN", "ELSERETURN"),
	}
}
