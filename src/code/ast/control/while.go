package control

import (
	blockly2 "Falcon/code/ast/blockly"
	"Falcon/code/sugar"
)

type While struct {
	Condition blockly2.Expr
	Body      []blockly2.Expr
}

func (w *While) String() string {
	return sugar.Format("while % {\n%}", w.Condition.String(), blockly2.PadBody(w.Body))
}

func (w *While) Blockly() blockly2.Block {
	return blockly2.Block{
		Type:       "controls_while",
		Values:     []blockly2.Value{{Name: "TEST", Block: w.Condition.Blockly()}},
		Statements: []blockly2.Statement{blockly2.CreateStatement("DO", w.Body)},
	}
}

func (w *While) Continuous() bool {
	return false
}

func (w *While) Consumable() bool {
	return false
}
