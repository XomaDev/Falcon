package fundamentals

import (
	"Falcon/code/ast/blockly"
	"Falcon/code/sugar"
	"strings"
)

type Dictionary struct {
	Elements []blockly.Expr
}

func (d *Dictionary) Yail() string {
	yail := "(call-yail-primitive make-yail-dictionary (*list-for-runtime* "
	yail += blockly.JoinYailExprs(" ", d.Elements)
	yail += ") '("
	yail += strings.Repeat("pair ", len(d.Elements))
	yail += ") \"make a dictionary\")"
	return yail
}

func (d *Dictionary) String() string {
	return sugar.Format("{ % }", blockly.JoinExprs(", ", d.Elements))
}

func (d *Dictionary) Blockly() blockly.Block {
	return blockly.Block{
		Type:     "dictionaries_create_with",
		Mutation: &blockly.Mutation{ItemCount: len(d.Elements)},
		Values:   blockly.ValuesByPrefix("ADD", d.Elements),
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

func (w *WalkAll) Yail() string {
	return "(static-field com.google.appinventor.components.runtime.util.YailDictionary 'ALL)"
}

func (w *WalkAll) String() string {
	return "walkAll"
}

func (w *WalkAll) Blockly() blockly.Block {
	return blockly.Block{Type: "dictionaries_walk_all"}
}

func (w *WalkAll) Continuous() bool {
	return true
}

func (w *WalkAll) Consumable() bool {
	return true
}
