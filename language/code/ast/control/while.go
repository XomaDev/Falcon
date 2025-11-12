package control

import (
	"Falcon/code/ast"
	"Falcon/code/sugar"
)

type While struct {
	Condition ast.Expr
	Body      []ast.Expr
}

func (w *While) Yail() string {
	yail := "(while "
	yail += w.Condition.Yail()
	yail += " (begin "
	yail += ast.PadBodyYail(w.Body)
	yail += "))"
	return yail
}

func (w *While) String() string {
	return sugar.Format("while % {\n%}", w.Condition.String(), ast.PadBody(w.Body))
}

func (w *While) Blockly() ast.Block {
	return ast.Block{
		Type:       "controls_while",
		Values:     []ast.Value{{Name: "TEST", Block: w.Condition.Blockly()}},
		Statements: []ast.Statement{ast.CreateStatement("DO", w.Body)},
	}
}

func (w *While) Continuous() bool {
	return false
}

func (w *While) Consumable() bool {
	return false
}

func (w *While) Signature() []ast.Signature {
	return []ast.Signature{ast.SignVoid}
}
