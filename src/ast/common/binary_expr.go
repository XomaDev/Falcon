package common

import (
	"Falcon/ast/blockly"
	l "Falcon/lex"
	"strconv"
)

type BinaryExpr struct {
	Where    l.Token
	Operands []blockly.Expr
	Operator l.Type
}

func (b *BinaryExpr) String() string {
	return blockly.JoinExprs(" "+*b.Where.Content+" ", b.Operands)
}

func (b *BinaryExpr) Blockly() blockly.Block {
	switch b.Operator {
	case l.BitwiseAnd, l.BitwiseOr, l.BitwiseXor:
		return b.bitwiseExpr()
	case l.Equals, l.NotEquals:
		return b.compareExpr()
	case l.LogicAnd, l.LogicOr:
		return b.boolExpr()
	case l.Colon:
		return b.pairExpr()
	default:
		return b.mathExpr()
	}
}

func (b *BinaryExpr) pairExpr() blockly.Block {
	if len(b.Operands) != 2 {
		b.Where.Error("Pair operator ':' received more than two operands")
	}
	return blockly.Block{Type: "pair", Values: blockly.MakeValues(b.Operands, "KEY", "VALUE")}
}

func (b *BinaryExpr) boolExpr() blockly.Block {
	var fieldOp string
	if b.Operator == l.LogicAnd {
		fieldOp = "AND"
	} else {
		fieldOp = "OR"
	}
	values := []blockly.Value{
		{Name: "A", Block: b.Operands[0].Blockly()},
		{Name: "B", Block: b.Operands[1].Blockly()},
	}
	lenOperands := len(b.Operands)
	if lenOperands > 2 {
		for i := 2; i < lenOperands; i++ {
			values = append(values, blockly.Value{Name: "BOOL" + strconv.Itoa(i), Block: b.Operands[i].Blockly()})
		}
	}
	return blockly.Block{
		Type:     "logic_operation",
		Mutation: &blockly.Mutation{ItemCount: lenOperands},
		Values:   values,
		Fields:   []blockly.Field{{Name: "OP", Value: fieldOp}},
	}
}

func (b *BinaryExpr) compareExpr() blockly.Block {
	var fieldOp string
	if b.Operator == l.Equals {
		fieldOp = "EQ"
	} else {
		fieldOp = "NEQ"
	}
	return blockly.Block{
		Type:   "logic_compare",
		Values: blockly.MakeValues(b.Operands, "A", "B"),
		Fields: []blockly.Field{{Name: "OP", Value: fieldOp}},
	}
}

func (b *BinaryExpr) bitwiseExpr() blockly.Block {
	var fieldOp string
	switch b.Operator {
	case l.BitwiseAnd:
		fieldOp = "BITAND"
	case l.BitwiseOr:
		fieldOp = "BITIOR"
	case l.BitwiseXor:
		fieldOp = "BITXOR"
	}
	return blockly.Block{
		Type:     "math_bitwise",
		Values:   blockly.ToValues("NUM", b.Operands),
		Mutation: &blockly.Mutation{ItemCount: len(b.Operands)},
		Fields:   []blockly.Field{{Name: "OP", Value: fieldOp}},
	}
}

func (b *BinaryExpr) mathExpr() blockly.Block {
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
		b.Where.Error("Unknown math operator (%)", b.Operator.String())
	}
	return blockly.Block{
		Type:     blockType,
		Values:   blockly.ToValues("NUM", b.Operands),
		Mutation: &blockly.Mutation{ItemCount: len(b.Operands)},
	}
}
