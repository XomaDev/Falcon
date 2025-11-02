package common

import (
	ast2 "Falcon/code/ast"
	list2 "Falcon/code/ast/list"
	"Falcon/code/ast/variables"
	lex2 "Falcon/code/lex"
	"strconv"
)

type BinaryExpr struct {
	Where    *lex2.Token
	Operands []ast2.Expr
	Operator lex2.Type
}

func (b *BinaryExpr) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (b *BinaryExpr) String() string {
	return ast2.JoinExprs(" "+*b.Where.Content+" ", b.Operands)
}

// CanRepeat: return true if the binary expr can be optimized into one struct
//	without the need to create additional BinaryExpr struct for the same Operator.
//	This factor also depends on the type of Operator being used. (Some support, some don't)

func (b *BinaryExpr) CanRepeat(testOperator lex2.Type) bool {
	if b.Operator != testOperator {
		return false
	}
	switch b.Operator {
	case lex2.Power, lex2.Dash, lex2.Slash, lex2.Colon:
		return false
	default:
		return true
	}
}

func (b *BinaryExpr) Blockly() ast2.Block {
	switch b.Operator {
	case lex2.BitwiseAnd, lex2.BitwiseOr, lex2.BitwiseXor:
		return b.bitwiseExpr()
	case lex2.Equals, lex2.NotEquals:
		return b.compareExpr()
	case lex2.LogicAnd, lex2.LogicOr:
		return b.boolExpr()
	case lex2.Colon:
		return b.pairExpr()
	case lex2.Plus, lex2.Times:
		return b.addOrTimes()
	case lex2.Dash, lex2.Slash, lex2.Power:
		return b.simpleMathExpr()
	case lex2.Underscore:
		return b.textJoin()
	case lex2.LessThan, lex2.LessThanEqual, lex2.GreatThan, lex2.GreaterThanEqual:
		return b.relationalExpr()
	case lex2.TextEquals, lex2.TextNotEquals, lex2.TextLessThan, lex2.TextGreaterThan:
		return b.textCompare()
	case lex2.Assign:
		return b.assignment()
	default:
		b.Where.Error("Unknown binary operator! " + b.Operator.String())
		panic("") // unreachable
	}
}

func (b *BinaryExpr) Continuous() bool {
	return false
}

func (b *BinaryExpr) Consumable() bool {
	return b.Operator != lex2.Assign
}

func (b *BinaryExpr) assignment() ast2.Block {
	if len(b.Operands) != 2 {
		b.Where.Error("Assignment '=' received more than two operands")
	}
	settable := b.Operands[0]
	newValue := b.Operands[1]

	if listGet, ok := settable.(*list2.Get); ok {
		listSet := list2.Set{List: listGet.List, Index: listGet.Index, Value: newValue}
		return listSet.Blockly()
	} else if varGet, ok := settable.(*variables.Get); ok {
		var name string
		if varGet.Global {
			name = "global " + varGet.Name
		} else {
			name = varGet.Name
		}
		return ast2.Block{
			Type:   "lexical_variable_set",
			Fields: []ast2.Field{{Name: "VAR", Value: name}},
			Values: []ast2.Value{{Name: "VALUE", Block: newValue.Blockly()}},
		}
	}
	panic("Unimplemented!")
}

func (b *BinaryExpr) textCompare() ast2.Block {
	var fieldOp string
	switch b.Operator {
	case lex2.TextEquals:
		fieldOp = "EQUAL"
	case lex2.NotEquals:
		fieldOp = "NEQ"
	case lex2.TextLessThan:
		fieldOp = "LT"
	case lex2.TextGreaterThan:
		fieldOp = "GT"
	}
	return ast2.Block{
		Type:   "text_compare",
		Fields: []ast2.Field{{Name: "OP", Value: fieldOp}},
		Values: ast2.MakeValues(b.Operands, "TEXT1", "TEXT2"),
	}
}

func (b *BinaryExpr) relationalExpr() ast2.Block {
	var fieldOp string
	switch b.Operator {
	case lex2.LessThan:
		fieldOp = "LT"
	case lex2.LessThanEqual:
		fieldOp = "LT"
	case lex2.GreatThan:
		fieldOp = "GT"
	case lex2.GreaterThanEqual:
		fieldOp = "GTE"
	}
	return ast2.Block{
		Type:   "math_compare",
		Fields: []ast2.Field{{Name: "OP", Value: fieldOp}},
		Values: ast2.MakeValues(b.Operands, "A", "B"),
	}
}

func (b *BinaryExpr) textJoin() ast2.Block {
	return ast2.Block{
		Type:     "text_join",
		Mutation: &ast2.Mutation{ItemCount: len(b.Operands)},
		Values:   ast2.ValuesByPrefix("ADD", b.Operands),
	}
}

func (b *BinaryExpr) pairExpr() ast2.Block {
	if len(b.Operands) != 2 {
		b.Where.Error("Pair operator ':' received more than two operands")
	}
	return ast2.Block{
		Type:   "pair",
		Values: ast2.MakeValues(b.Operands, "KEY", "VALUE"),
	}
}

func (b *BinaryExpr) boolExpr() ast2.Block {
	var fieldOp string
	if b.Operator == lex2.LogicAnd {
		fieldOp = "AND"
	} else {
		fieldOp = "OR"
	}
	values := []ast2.Value{
		{Name: "A", Block: b.Operands[0].Blockly()},
		{Name: "B", Block: b.Operands[1].Blockly()},
	}
	lenOperands := len(b.Operands)
	if lenOperands > 2 {
		for i := 2; i < lenOperands; i++ {
			values = append(values, ast2.Value{Name: "BOOL" + strconv.Itoa(i), Block: b.Operands[i].Blockly()})
		}
	}
	return ast2.Block{
		Type:     "logic_operation",
		Mutation: &ast2.Mutation{ItemCount: lenOperands},
		Values:   values,
		Fields:   []ast2.Field{{Name: "OP", Value: fieldOp}},
	}
}

func (b *BinaryExpr) compareExpr() ast2.Block {
	var fieldOp string
	if b.Operator == lex2.Equals {
		fieldOp = "EQ"
	} else {
		fieldOp = "NEQ"
	}
	return ast2.Block{
		Type:   "logic_compare",
		Values: ast2.MakeValues(b.Operands, "A", "B"),
		Fields: []ast2.Field{{Name: "OP", Value: fieldOp}},
	}
}

func (b *BinaryExpr) bitwiseExpr() ast2.Block {
	var fieldOp string
	switch b.Operator {
	case lex2.BitwiseAnd:
		fieldOp = "BITAND"
	case lex2.BitwiseOr:
		fieldOp = "BITIOR"
	case lex2.BitwiseXor:
		fieldOp = "BITXOR"
	}
	return ast2.Block{
		Type:     "math_bitwise",
		Values:   ast2.ValuesByPrefix("NUM", b.Operands),
		Mutation: &ast2.Mutation{ItemCount: len(b.Operands)},
		Fields:   []ast2.Field{{Name: "OP", Value: fieldOp}},
	}
}

func (b *BinaryExpr) simpleMathExpr() ast2.Block {
	var blockType string
	switch b.Operator {
	case lex2.Dash:
		blockType = "math_subtract"
	case lex2.Slash:
		blockType = "math_division"
	case lex2.Power:
		blockType = "math_power"
	}
	return ast2.Block{
		Type:   blockType,
		Values: ast2.MakeValues(b.Operands, "A", "B"),
	}
}

func (b *BinaryExpr) addOrTimes() ast2.Block {
	var blockType string
	if b.Operator == lex2.Plus {
		blockType = "math_add"
	} else {
		blockType = "math_multiply"
	}
	return ast2.Block{
		Type:     blockType,
		Values:   ast2.ValuesByPrefix("NUM", b.Operands),
		Mutation: &ast2.Mutation{ItemCount: len(b.Operands)},
	}
}
