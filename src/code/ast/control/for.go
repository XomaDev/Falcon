package control

import (
	blockly2 "Falcon/code/ast/blockly"
	"Falcon/code/sugar"
)

type For struct {
	IName string
	From  blockly2.Expr
	To    blockly2.Expr
	By    blockly2.Expr
	Body  []blockly2.Expr
}

func (f *For) String() string {
	return sugar.Format("for %: % to % by % {\n%}",
		f.IName, f.From.String(), f.To.String(), f.By.String(), blockly2.PadBody(f.Body))
}

func (f *For) Blockly() blockly2.Block {
	return blockly2.Block{
		Type:       "controls_forRange",
		Fields:     []blockly2.Field{{Name: "VAR", Value: f.IName}},
		Values:     blockly2.MakeValues([]blockly2.Expr{f.From, f.To, f.By}, "START", "END", "STEP"),
		Statements: []blockly2.Statement{blockly2.CreateStatement("DO", f.Body)},
	}
}

func (f *For) Continuous() bool {
	return false
}

func (f *For) Consumable() bool {
	return false
}
