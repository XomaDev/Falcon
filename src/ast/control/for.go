package control

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
)

type For struct {
	IName string
	From  blockly.Expr
	To    blockly.Expr
	By    blockly.Expr
	Body  []blockly.Expr
}

func (f *For) String() string {
	return sugar.Format("for %: % to % by % {\n%}",
		f.IName, f.From.String(), f.To.String(), f.By.String(), blockly.PadBody(f.Body))
}

func (f *For) Blockly() blockly.Block {
	return blockly.Block{
		Type:       "controls_forRange",
		Fields:     []blockly.Field{{Name: "VAR", Value: f.IName}},
		Values:     blockly.MakeValues([]blockly.Expr{f.From, f.To, f.By}, "START", "END", "STEP"),
		Statements: []blockly.Statement{blockly.CreateStatement("DO", f.Body)},
	}
}

func (f *For) Continuous() bool {
	return false
}

func (f *For) Consumable() bool {
	return false
}
