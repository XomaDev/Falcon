package common

import (
	ast2 "Falcon/lang/code/ast"
	"Falcon/lang/code/lex"
	"Falcon/lang/code/sugar"
)

type Question struct {
	Where    *lex.Token
	On       ast2.Expr
	Question string
}

func (q *Question) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (q *Question) String() string {
	pFormat := "% ? %"
	if !q.On.Continuous() {
		pFormat = "(%) ? %"
	}
	return sugar.Format(pFormat, q.On.String(), q.Question)
}

func (q *Question) Blockly() ast2.Block {
	switch q.Question {
	case "number", "base10", "hexa", "bin":
		return q.mathQuestion()
	case "text":
		return q.textQuestion()
	case "list":
		return q.listQuestion()
	case "dict":
		return q.dictQuestion()
	case "emptyText":
		return q.textIsEmpty()
	case "emptyList":
		return q.listIsEmpty()
	default:
		q.Where.Error("Unknown question ? %", q.Question)
	}
	panic("Unreachable")
}

func (q *Question) Continuous() bool {
	return false
}

func (q *Question) Consumable() bool {
	return true
}

func (q *Question) listIsEmpty() ast2.Block {
	return ast2.Block{
		Type:   "lists_is_empty",
		Values: []ast2.Value{{Name: "LIST", Block: q.On.Blockly()}},
	}
}

func (q *Question) textIsEmpty() ast2.Block {
	return ast2.Block{
		Type:   "text_isEmpty",
		Values: []ast2.Value{{Name: "VALUE", Block: q.On.Blockly()}},
	}
}

func (q *Question) dictQuestion() ast2.Block {
	return ast2.Block{
		Type:   "dictionaries_is_dict",
		Values: []ast2.Value{{Name: "THING", Block: q.On.Blockly()}},
	}
}

func (q *Question) listQuestion() ast2.Block {
	return ast2.Block{
		Type:   "lists_is_list",
		Values: []ast2.Value{{Name: "ITEM", Block: q.On.Blockly()}},
	}
}

func (q *Question) textQuestion() ast2.Block {
	return ast2.Block{
		Type:   "text_is_string",
		Values: []ast2.Value{{Name: "ITEM", Block: q.On.Blockly()}},
	}
}

func (q *Question) mathQuestion() ast2.Block {
	var fieldOp string
	switch q.Question {
	case "number":
		fieldOp = "NUMBER"
	case "base10":
		fieldOp = "BASE10"
	case "hexa":
		fieldOp = "HEXADECIMAL"
	case "bin":
		fieldOp = "BINARY"
	}
	return ast2.Block{
		Type:   "math_is_a_number",
		Fields: []ast2.Field{{Name: "OP", Value: fieldOp}},
		Values: []ast2.Value{{Name: "NUM", Block: q.On.Blockly()}},
	}
}
