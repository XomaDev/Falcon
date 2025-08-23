package control

import (
	blky "Falcon/ast/blockly"
	"Falcon/sugar"
)

type Do struct {
	Body   []blky.Expr
	Result blky.Expr
}

func (d *Do) String() string {
	return sugar.Format("do {\n%} -> %", blky.PadBody(d.Body), d.Result.String())
}

func (d *Do) Blockly() blky.Block {
	return blky.Block{
		Type:       "controls_do_then_return",
		Statements: []blky.Statement{blky.CreateStatement("STM", d.Body)},
		Values:     []blky.Value{{Name: "VALUE", Block: d.Result.Blockly()}},
	}
}

func (d *Do) Continuous() bool {
	return false
}

func (d *Do) Consumable() bool {
	return false
}
