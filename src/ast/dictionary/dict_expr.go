package dictionary

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
)

type Dictionary struct {
	Elements []blockly.Expr
}

func (d *Dictionary) String() string {
	return sugar.Format("{ % }", blockly.JoinExprs(", ", d.Elements))
}

func (d *Dictionary) Blockly() blockly.Block {
	return blockly.Block{
		Type:       "dictionaries_create_with",
		Mutation:   &blockly.Mutation{ItemCount: len(d.Elements)},
		Values:     blockly.ValuesByPrefix("ADD", d.Elements),
		Consumable: true,
	}
}
