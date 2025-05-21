package common

import (
	"Falcon/ast/blockly"
	"Falcon/ast/list"
	"Falcon/ast/variables"
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

// CanRepeat: return true if the binary expr can be optimized into one struct
//
//	without the need to create additional BinaryExpr struct for the same Operator.
//	This factor also depends on the type of Operator being used. (Some support, some don't)
func (b *BinaryExpr) CanRepeat(testOperator l.Type) bool {
	if b.Operator != testOperator {
		return false
	}
	switch b.Operator {
	case l.Power, l.Dash, l.Slash, l.Colon:
		return false
	default:
		return true
	}
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
	case l.Plus, l.Times:
		return b.addOrTimes()
	case l.Dash, l.Slash, l.Power:
		return b.simpleMathExpr()
	case l.Underscore:
		return b.textJoin()
	case l.LessThan, l.LessThanEqual, l.GreatThan, l.GreaterThanEqual:
		return b.relationalExpr()
	case l.TextEquals, l.TextNotEquals, l.TextLessThan, l.TextGreaterThan:
		return b.textCompare()
	case l.Assign:
		return b.assignment()
	default:
		b.Where.Error("Unknown binary operator! " + b.Operator.String())
		panic("") // unreachable
	}
}

func (b *BinaryExpr) assignment() blockly.Block {
	if len(b.Operands) != 2 {
		b.Where.Error("Assignment '=' received more than two operands")
	}
	settable := b.Operands[0]
	newValue := b.Operands[1]

	if listGet, ok := settable.(*list.Get); ok {
		listSet := list.Set{List: listGet.List, Index: listGet.Index, Value: newValue}
		return listSet.Blockly()
	} else if varGet, ok := settable.(*variables.Get); ok {
		var name string
		if varGet.Global {
			name = "global " + varGet.Name
		} else {
			name = varGet.Name
		}
		return blockly.Block{
			Type:   "lexical_variable_set",
			Fields: []blockly.Field{{Name: "VAR", Value: name}},
			Values: []blockly.Value{{Name: "VALUE", Block: newValue.Blockly()}},
		}
	}
	panic("Unimplemented!")
}

func (b *BinaryExpr) textCompare() blockly.Block {
	var fieldOp string
	switch b.Operator {
	case l.TextEquals:
		fieldOp = "EQUAL"
	case l.NotEquals:
		fieldOp = "NEQ"
	case l.TextLessThan:
		fieldOp = "LT"
	case l.TextGreaterThan:
		fieldOp = "GT"
	}
	return blockly.Block{
		Type:   "text_compare",
		Fields: []blockly.Field{{Name: "OP", Value: fieldOp}},
		Values: blockly.MakeValues(b.Operands, "TEXT1", "TEXT2"),
	}
}

func (b *BinaryExpr) relationalExpr() blockly.Block {
	var fieldOp string
	switch b.Operator {
	case l.LessThan:
		fieldOp = "LT"
	case l.LessThanEqual:
		fieldOp = "LT"
	case l.GreatThan:
		fieldOp = "GT"
	case l.GreaterThanEqual:
		fieldOp = "GTE"
	}
	return blockly.Block{
		Type:   "math_compare",
		Fields: []blockly.Field{{Name: "OP", Value: fieldOp}},
		Values: blockly.MakeValues(b.Operands, "A", "B"),
	}
}

func (b *BinaryExpr) textJoin() blockly.Block {
	return blockly.Block{
		Type:     "text_join",
		Mutation: &blockly.Mutation{ItemCount: len(b.Operands)},
		Values:   blockly.ValuesByPrefix("ADD", b.Operands),
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
		Values:   blockly.ValuesByPrefix("NUM", b.Operands),
		Mutation: &blockly.Mutation{ItemCount: len(b.Operands)},
		Fields:   []blockly.Field{{Name: "OP", Value: fieldOp}},
	}
}

func (b *BinaryExpr) simpleMathExpr() blockly.Block {
	var blockType string
	switch b.Operator {
	case l.Dash:
		blockType = "math_subtract"
	case l.Slash:
		blockType = "math_division"
	case l.Power:
		blockType = "math_power"
	}
	return blockly.Block{Type: blockType, Values: blockly.MakeValues(b.Operands, "A", "B")}
}

func (b *BinaryExpr) addOrTimes() blockly.Block {
	var blockType string
	if b.Operator == l.Plus {
		blockType = "math_add"
	} else {
		blockType = "math_multiply"
	}
	return blockly.Block{
		Type:     blockType,
		Values:   blockly.ValuesByPrefix("NUM", b.Operands),
		Mutation: &blockly.Mutation{ItemCount: len(b.Operands)},
	}
}
