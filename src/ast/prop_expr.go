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
	default:
		p.Where.Error("Unknown property access ->%", *p.Name)
	}
	return Block{
		Type:   "math_single",
		Fields: []Field{{Name: "OP", Value: fieldOp}},
		Values: []Value{{Name: "NUM", Block: p.On.Blockly()}},
	}
}
