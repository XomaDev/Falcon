package control

import (
	"Falcon/ast/blockly"
	"Falcon/sugar"
)

type While struct {
	Condition blockly.Expr
	Body      []blockly.Expr
}

func (w *While) String() string {
	return sugar.Format("while % {\n%}", w.Condition.String(), blockly.PadBody(w.Body))
}

func (w *While) Blockly() blockly.Block {
	return blockly.Block{
		Type:       "controls_while",
		Values:     []blockly.Value{{Name: "TEST", Block: w.Condition.Blockly()}},
		Statements: []blockly.Statement{blockly.CreateStatement("DO", w.Body)},
	}
}

func (w *While) Continuous() bool {
	return false
}

func (w *While) Consumable() bool {
	return false
}
