package ast

import (
	"Falcon/sugar"
	"Falcon/types"
)

type PropExpr struct {
	Where types.Token
	On    Expr
	Name  *string
}

func (p *PropExpr) String() string {
	return sugar.Format("%->%v", p.On.String(), *p.Name)
}

func (p *PropExpr) Blockly() Block {
	var fieldOp string
	switch *p.Name {
	case "root":
		fieldOp = "ROOT"
	case "abs":
		fieldOp = "ABS"
	case "neg":
		fieldOp = "NEG"
	case "log":
		fieldOp = "LN"
	case "exp":
		fieldOp = "EXP"
	case "round":
		fieldOp = "ROUND"
	case "ceil":
		fieldOp = "CEILING"
	case "floor":
		fieldOp = "FLOOR"
	case "sin":
		fieldOp = "SIN"
	case "cos":
		fieldOp = "COS"
	case "tan":
		fieldOp = "TAN"
	case "asin":
		fieldOp = "ASIN"
	case "acos":
		fieldOp = "ACOS"
	case "atan":
		fieldOp = "ATAN"
	case "degrees":
		fieldOp = "RADIANS_TO_DEGREES"
	case "radians":
		fieldOp = "DEGREES_TO_RADIANS"
	default:
		p.Where.Error("Unknown property access ->%", *p.Name)
	}
	var blockType string
	if fieldOp == "SIN" || fieldOp == "COS" || fieldOp == "TAN" ||
		fieldOp == "ASIN" || fieldOp == "ACOS" || fieldOp == "ATAN" {
		blockType = "math_trig"
	} else if fieldOp == "RADIANS_TO_DEGREES" || fieldOp == "DEGREES_TO_RADIANS" {
		blockType = "math_convert_angles"
	} else {
		blockType = "math_single"
	}
	return Block{
		Type:   blockType,
		Fields: []Field{{Name: "OP", Value: fieldOp}},
		Values: []Value{{Name: "NUM", Block: p.On.Blockly()}},
	}
}
