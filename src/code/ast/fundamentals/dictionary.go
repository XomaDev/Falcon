package fundamentals

import (
	blockly2 "Falcon/code/ast/blockly"
	"Falcon/code/sugar"
)

type Dictionary struct {
	Elements []blockly2.Expr
}

func (d *Dictionary) String() string {
	return sugar.Format("{ % }", blockly2.JoinExprs(", ", d.Elements))
}

func (d *Dictionary) Blockly() blockly2.Block {
	return blockly2.Block{
		Type:     "dictionaries_create_with",
		Mutation: &blockly2.Mutation{ItemCount: len(d.Elements)},
		Values:   blockly2.ValuesByPrefix("ADD", d.Elements),
	}
}

func (d *Dictionary) Continuous() bool {
	return true
}

func (d *Dictionary) Consumable() bool {
	return true
}

type WalkAll struct {
}

func (w *WalkAll) String() string {
	return "walkAll"
}

func (w *WalkAll) Blockly() blockly2.Block {
	return blockly2.Block{Type: "dictionaries_walk_all"}
}

func (w *WalkAll) Continuous() bool {
	return true
}

func (w *WalkAll) Consumable() bool {
	return true
}
