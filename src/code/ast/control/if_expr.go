package control

import (
	blockly2 "Falcon/code/ast/blockly"
	"Falcon/code/sugar"
)

type SimpleIf struct {
	Condition blockly2.Expr
	Then      blockly2.Expr
	Else      blockly2.Expr
}

func (s *SimpleIf) String() string {
	return sugar.Format("if (%) % else %", s.Condition.String(), s.Then.String(), s.Else.String())
}

func (s *SimpleIf) Blockly() blockly2.Block {
	return blockly2.Block{
		Type:   "controls_choose",
		Values: blockly2.MakeValues([]blockly2.Expr{s.Condition, s.Then, s.Else}, "TEST", "THENRETURN", "ELSERETURN"),
	}
}

func (s *SimpleIf) Continuous() bool {
	return false
}

func (s *SimpleIf) Consumable() bool {
	return true
}
