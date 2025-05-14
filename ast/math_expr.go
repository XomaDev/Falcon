package ast

import (
	"Falcon/types"
)

type MathExpr struct {
	Operands []Expr
	Operator types.Token
}

func (b *MathExpr) String() string {
	return JoinExprs(*b.Operator.Content, b.Operands)
}

func (b *MathExpr) Blockly() Block {
	operator := *b.Operator.Content
	var blockType string

	switch operator {
	case "+":
		blockType = "math_add"
	case "-":
		blockType = "math_subtract"
	case "*":
		blockType = "math_multiply"
	case "/":
		blockType = "math_division"
	case "^":
		blockType = "math_power"
	case "&", "|", "~":
		blockType = "math_bitwise"
	default:
		b.Operator.Error("Unknown binary operator (%v)", *b.Operator.Content)
	}
	var fields []Field
	if blockType == "math_bitwise" {
		var fieldOp string
		switch operator {
		case "&":
			fieldOp = "BITAND"
		case "|":
			fieldOp = "BITIOR"
		case "~":
			fieldOp = "BITXOR"
		}
		fields = append(fields, Field{Name: "OP", Value: fieldOp})
	}
	return Block{
		Type:     blockType,
		Values:   ToValues("NUM", b.Operands),
		Mutation: &Mutation{ItemCount: len(b.Operands)},
		Fields:   fields,
	}
}
