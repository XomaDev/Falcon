package common

import (
	"Falcon/ast/blockly"
	"Falcon/lex"
	"Falcon/sugar"
)

type Question struct {
	Where    lex.Token
	On       blockly.Expr
	Question string
}

func (q *Question) String() string {
	return sugar.Format("% ? %", q.On.String(), q.Question)
}

func (q *Question) Blockly() blockly.Block {
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

func (q *Question) listIsEmpty() blockly.Block {
	return blockly.Block{
		Type:       "lists_is_empty",
		Values:     []blockly.Value{{Name: "LIST", Block: q.On.Blockly()}},
		Consumable: true,
	}
}

func (q *Question) textIsEmpty() blockly.Block {
	return blockly.Block{
		Type:       "text_isEmpty",
		Values:     []blockly.Value{{Name: "VALUE", Block: q.On.Blockly()}},
		Consumable: true,
	}
}

func (q *Question) dictQuestion() blockly.Block {
	return blockly.Block{
		Type:       "dictionaries_is_dict",
		Values:     []blockly.Value{{Name: "THING", Block: q.On.Blockly()}},
		Consumable: true,
	}
}

func (q *Question) listQuestion() blockly.Block {
	return blockly.Block{
		Type:       "lists_is_list",
		Values:     []blockly.Value{{Name: "ITEM", Block: q.On.Blockly()}},
		Consumable: true,
	}
}

func (q *Question) textQuestion() blockly.Block {
	return blockly.Block{
		Type:       "text_is_string",
		Values:     []blockly.Value{{Name: "ITEM", Block: q.On.Blockly()}},
		Consumable: true,
	}
}

func (q *Question) mathQuestion() blockly.Block {
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
	return blockly.Block{
		Type:       "math_is_a_number",
		Fields:     []blockly.Field{{Name: "OP", Value: fieldOp}},
		Values:     []blockly.Value{{Name: "NUM", Block: q.On.Blockly()}},
		Consumable: true,
	}
}
