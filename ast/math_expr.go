package ast

import (
	"Falcon/types"
)

type MathExpr struct {
	Operands []Expr
	Operator types.Token
}

func (b *MathExpr) String() string {
	return JoinExprs(b.Operator.String(), b.Operands)
}

func (b *MathExpr) Blockly() Block {
	var blockType string
	switch *b.Operator.Content {
	case "+":
		blockType = "math_add"
	case "-":
		blockType = "math_subtract"
	case "*":
		blockType = "math_multiply"
	case "/":
		blockType = "math_division"
	default:
		b.Operator.Error("Unknown binary operator (%v)", *b.Operator.Content)
	}
	return Block{
		Type:     blockType,
		Values:   ToValues("NUM", b.Operands),
		Mutation: &Mutation{ItemCount: len(b.Operands)},
	}
}
