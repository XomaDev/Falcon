package control

import (
	"Falcon/code/ast/blockly"
	"Falcon/code/sugar"
)

type Do struct {
	Body   []blockly.Expr
	Result blockly.Expr
}

func (d *Do) String() string {
	return sugar.Format("do {\n%} -> %", blockly.PadBody(d.Body), d.Result.String())
}

func (d *Do) Blockly() blockly.Block {
	return blockly.Block{
		Type:       "controls_do_then_return",
		Statements: []blockly.Statement{blockly.CreateStatement("STM", d.Body)},
		Values:     []blockly.Value{{Name: "VALUE", Block: d.Result.Blockly()}},
	}
}

func (d *Do) Continuous() bool {
	return false
}

func (d *Do) Consumable() bool {
	return false
}
