package control

import (
	ast2 "Falcon/lang/code/ast"
	"Falcon/lang/code/sugar"
)

type While struct {
	Condition ast2.Expr
	Body      []ast2.Expr
}

func (w *While) Yail() string {
	yail := "(while "
	yail += w.Condition.Yail()
	yail += " (begin "
	yail += ast2.PadBodyYail(w.Body)
	yail += "))"
	return yail
}

func (w *While) String() string {
	return sugar.Format("while % {\n%}", w.Condition.String(), ast2.PadBody(w.Body))
}

func (w *While) Blockly() ast2.Block {
	return ast2.Block{
		Type:       "controls_while",
		Values:     []ast2.Value{{Name: "TEST", Block: w.Condition.Blockly()}},
		Statements: []ast2.Statement{ast2.CreateStatement("DO", w.Body)},
	}
}

func (w *While) Continuous() bool {
	return false
}

func (w *While) Consumable() bool {
	return false
}
