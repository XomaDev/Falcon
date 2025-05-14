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
	default:
		q.Where.Error("Unknown question ? %", q.Question)
	}
	return blockly.Block{
		Type:   "math_is_a_number",
		Fields: []blockly.Field{{Name: "OP", Value: fieldOp}},
		Values: []blockly.Value{{Name: "NUM", Block: q.On.Blockly()}},
	}
}
