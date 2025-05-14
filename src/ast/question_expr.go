package ast

import (
	"Falcon/sugar"
	"Falcon/types"
)

type QuestionExp struct {
	Where    types.Token
	On       Expr
	Question *string
}

func (q *QuestionExp) String() string {
	return sugar.Format("% ? %", q.On.String(), *q.Question)
}

func (q *QuestionExp) Blockly() Block {
	var fieldOp string
	switch *q.Question {
	case "number":
		fieldOp = "NUMBER"
	case "base10":
		fieldOp = "BASE10"
	case "hexa":
		fieldOp = "HEXADECIMAL"
	case "bin":
		fieldOp = "BINARY"
	default:
		q.Where.Error("Unknown question ? %", *q.Question)
	}
	return Block{
		Type:   "math_is_a_number",
		Fields: []Field{{Name: "OP", Value: fieldOp}},
		Values: []Value{{Name: "NUM", Block: q.On.Blockly()}},
	}
}
