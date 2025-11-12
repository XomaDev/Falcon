package common

import (
	"Falcon/code/ast"
	"Falcon/code/ast/list"
	"Falcon/code/ast/variables"
	"Falcon/code/lex"
	"strconv"
)

type BinaryExpr struct {
	Where    *lex.Token
	Operands []ast.Expr
	Operator lex.Type
}

func (b *BinaryExpr) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (b *BinaryExpr) String() string {
	return ast.JoinExprs(" "+*b.Where.Content+" ", b.Operands)
}

// CanRepeat: return true if the binary expr can be optimized into one struct
//	without the need to create additional BinaryExpr struct for the same Operator.
//	This factor also depends on the type of Operator being used. (Some support, some don't)

func (b *BinaryExpr) CanRepeat(testOperator lex.Type) bool {
	if b.Operator != testOperator {
		return false
	}
	switch b.Operator {
	case lex.Power, lex.Dash, lex.Slash, lex.Colon:
		return false
	default:
		return true
	}
}

func (b *BinaryExpr) Blockly() ast.Block {
	switch b.Operator {
	case lex.BitwiseAnd, lex.BitwiseOr, lex.BitwiseXor:
		return b.bitwiseExpr()
	case lex.Equals, lex.NotEquals:
		return b.compareExpr()
	case lex.LogicAnd, lex.LogicOr:
		return b.boolExpr()
	case lex.Colon:
		return b.pairExpr()
	case lex.Plus, lex.Times:
		return b.addOrTimes()
	case lex.Dash, lex.Slash, lex.Power:
		return b.simpleMathExpr()
	case lex.Underscore:
		return b.textJoin()
	case lex.LessThan, lex.LessThanEqual, lex.GreatThan, lex.GreaterThanEqual:
		return b.relationalExpr()
	case lex.TextEquals, lex.TextNotEquals, lex.TextLessThan, lex.TextGreaterThan:
		return b.textCompare()
	case lex.Assign:
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
	return b.Operator != lex.Assign
}

func (b *BinaryExpr) Signature() []ast.Signature {
	switch b.Operator {
	case lex.BitwiseAnd, lex.BitwiseOr, lex.BitwiseXor:
		return []ast.Signature{ast.SignNumb}
	case lex.Equals, lex.NotEquals:
		return []ast.Signature{ast.SignBool}
	case lex.LogicAnd, lex.LogicOr:
		return []ast.Signature{ast.SignBool}
	case lex.Colon:
		return []ast.Signature{ast.SignList}
	case lex.Plus, lex.Times:
		return []ast.Signature{ast.SignNumb}
	case lex.Dash, lex.Slash, lex.Power:
		return []ast.Signature{ast.SignNumb}
	case lex.Underscore:
		return []ast.Signature{ast.SignText}
	case lex.LessThan, lex.LessThanEqual, lex.GreatThan, lex.GreaterThanEqual:
		return []ast.Signature{ast.SignBool}
	case lex.TextEquals, lex.TextNotEquals, lex.TextLessThan, lex.TextGreaterThan:
		return []ast.Signature{ast.SignBool}
	case lex.Assign:
		return []ast.Signature{ast.SignVoid}
	default:
		b.Where.Error("Unknown binary operator! " + b.Operator.String())
		panic("") // unreachable
	}
}

func (b *BinaryExpr) assignment() ast.Block {
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
		return ast.Block{
			Type:   "lexical_variable_set",
			Fields: []ast.Field{{Name: "VAR", Value: name}},
			Values: []ast.Value{{Name: "VALUE", Block: newValue.Blockly()}},
		}
	}
	panic("Unimplemented!")
}

func (b *BinaryExpr) textCompare() ast.Block {
	var fieldOp string
	switch b.Operator {
	case lex.TextEquals:
		fieldOp = "EQUAL"
	case lex.NotEquals:
		fieldOp = "NEQ"
	case lex.TextLessThan:
		fieldOp = "LT"
	case lex.TextGreaterThan:
		fieldOp = "GT"
	}
	return ast.Block{
		Type:   "text_compare",
		Fields: []ast.Field{{Name: "OP", Value: fieldOp}},
		Values: ast.MakeValues(b.Operands, "TEXT1", "TEXT2"),
	}
}

func (b *BinaryExpr) relationalExpr() ast.Block {
	var fieldOp string
	switch b.Operator {
	case lex.LessThan:
		fieldOp = "LT"
	case lex.LessThanEqual:
		fieldOp = "LT"
	case lex.GreatThan:
		fieldOp = "GT"
	case lex.GreaterThanEqual:
		fieldOp = "GTE"
	}
	return ast.Block{
		Type:   "math_compare",
		Fields: []ast.Field{{Name: "OP", Value: fieldOp}},
		Values: ast.MakeValues(b.Operands, "A", "B"),
	}
}

func (b *BinaryExpr) textJoin() ast.Block {
	return ast.Block{
		Type:     "text_join",
		Mutation: &ast.Mutation{ItemCount: len(b.Operands)},
		Values:   ast.ValuesByPrefix("ADD", b.Operands),
	}
}

func (b *BinaryExpr) pairExpr() ast.Block {
	if len(b.Operands) != 2 {
		b.Where.Error("Pair operator ':' received more than two operands")
	}
	return ast.Block{
		Type:   "pair",
		Values: ast.MakeValues(b.Operands, "KEY", "VALUE"),
	}
}

func (b *BinaryExpr) boolExpr() ast.Block {
	var fieldOp string
	if b.Operator == lex.LogicAnd {
		fieldOp = "AND"
	} else {
		fieldOp = "OR"
	}
	values := []ast.Value{
		{Name: "A", Block: b.Operands[0].Blockly()},
		{Name: "B", Block: b.Operands[1].Blockly()},
	}
	lenOperands := len(b.Operands)
	if lenOperands > 2 {
		for i := 2; i < lenOperands; i++ {
			values = append(values, ast.Value{Name: "BOOL" + strconv.Itoa(i), Block: b.Operands[i].Blockly()})
		}
	}
	return ast.Block{
		Type:     "logic_operation",
		Mutation: &ast.Mutation{ItemCount: lenOperands},
		Values:   values,
		Fields:   []ast.Field{{Name: "OP", Value: fieldOp}},
	}
}

func (b *BinaryExpr) compareExpr() ast.Block {
	var fieldOp string
	if b.Operator == lex.Equals {
		fieldOp = "EQ"
	} else {
		fieldOp = "NEQ"
	}
	return ast.Block{
		Type:   "logic_compare",
		Values: ast.MakeValues(b.Operands, "A", "B"),
		Fields: []ast.Field{{Name: "OP", Value: fieldOp}},
	}
}

func (b *BinaryExpr) bitwiseExpr() ast.Block {
	var fieldOp string
	switch b.Operator {
	case lex.BitwiseAnd:
		fieldOp = "BITAND"
	case lex.BitwiseOr:
		fieldOp = "BITIOR"
	case lex.BitwiseXor:
		fieldOp = "BITXOR"
	}
	return ast.Block{
		Type:     "math_bitwise",
		Values:   ast.ValuesByPrefix("NUM", b.Operands),
		Mutation: &ast.Mutation{ItemCount: len(b.Operands)},
		Fields:   []ast.Field{{Name: "OP", Value: fieldOp}},
	}
}

func (b *BinaryExpr) simpleMathExpr() ast.Block {
	var blockType string
	switch b.Operator {
	case lex.Dash:
		blockType = "math_subtract"
	case lex.Slash:
		blockType = "math_division"
	case lex.Power:
		blockType = "math_power"
	}
	return ast.Block{
		Type:   blockType,
		Values: ast.MakeValues(b.Operands, "A", "B"),
	}
}

func (b *BinaryExpr) addOrTimes() ast.Block {
	var blockType string
	if b.Operator == lex.Plus {
		blockType = "math_add"
	} else {
		blockType = "math_multiply"
	}
	return ast.Block{
		Type:     blockType,
		Values:   ast.ValuesByPrefix("NUM", b.Operands),
		Mutation: &ast.Mutation{ItemCount: len(b.Operands)},
	}
}
