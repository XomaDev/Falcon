package math

import (
	"Falcon/ast/blockly"
	"Falcon/types"
)

type MathExpr struct {
	Operands []blockly.Expr
	Operator types.Token
}

func (b *MathExpr) String() string {
	return blockly.JoinExprs(*b.Operator.Content, b.Operands)
}

func (b *MathExpr) Blockly() blockly.Block {
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
	var fields []blockly.Field
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
		fields = append(fields, blockly.Field{Name: "OP", Value: fieldOp})
	}
	return blockly.Block{
		Type:     blockType,
		Values:   blockly.ToValues("NUM", b.Operands),
		Mutation: &blockly.Mutation{ItemCount: len(b.Operands)},
		Fields:   fields,
	}
}
