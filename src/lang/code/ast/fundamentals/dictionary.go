package fundamentals

import (
	ast2 "Falcon/lang/code/ast"
	"Falcon/lang/code/sugar"
	"strings"
)

type Dictionary struct {
	Elements []ast2.Expr
}

func (d *Dictionary) Yail() string {
	return ast2.PrimitiveCall(
		"make-yail-dictionary",
		"make a dictionary",
		d.Elements,
		strings.Repeat("pair ", len(d.Elements)),
	)
}

func (d *Dictionary) String() string {
	return sugar.Format("{ % }", ast2.JoinExprs(", ", d.Elements))
}

func (d *Dictionary) Blockly() ast2.Block {
	return ast2.Block{
		Type:     "dictionaries_create_with",
		Mutation: &ast2.Mutation{ItemCount: len(d.Elements)},
		Values:   ast2.ValuesByPrefix("ADD", d.Elements),
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

func (w *WalkAll) Blockly() ast2.Block {
	return ast2.Block{Type: "dictionaries_walk_all"}
}

func (w *WalkAll) Continuous() bool {
	return true
}

func (w *WalkAll) Consumable() bool {
	return true
}
