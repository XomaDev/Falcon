package math

import (
	"Falcon/ast/blockly"
	l "Falcon/lex"
)

type Expr struct {
	Where    l.Token
	Operands []blockly.Expr
	Operator l.Type
}

func (b *Expr) String() string {
	return blockly.JoinExprs(" "+*b.Where.Content+" ", b.Operands)
}

func (b *Expr) Blockly() blockly.Block {
	var blockType string

	switch b.Operator {
	case l.Plus:
		blockType = "math_add"
	case l.Dash:
		blockType = "math_subtract"
	case l.Times:
		blockType = "math_multiply"
	case l.Slash:
		blockType = "math_division"
	case l.Power:
		blockType = "math_power"
	case l.BitwiseAnd, l.BitwiseOr, l.BitwiseXor:
		blockType = "math_bitwise"
	default:
		b.Where.Error("Unknown binary operator (%v)", b.Operator.String())
	}
	var fields []blockly.Field
	if blockType == "math_bitwise" {
		var fieldOp string
		switch b.Operator {
		case l.BitwiseAnd:
			fieldOp = "BITAND"
		case l.BitwiseOr:
			fieldOp = "BITIOR"
		case l.BitwiseXor:
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
