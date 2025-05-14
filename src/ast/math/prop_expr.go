package math

import (
	"Falcon/ast/blockly"
	"Falcon/lex"
	"Falcon/sugar"
)

type Prop struct {
	Where lex.Token
	On    blockly.Expr
	Name  string
}

func (p *Prop) String() string {
	return sugar.Format("%->%v", p.On.String(), p.Name)
}

func (p *Prop) Blockly() blockly.Block {
	var fieldOp string
	switch p.Name {
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
	case "hex":
		fieldOp = "DEC_TO_HEX"
	case "bin":
		fieldOp = "DEC_TO_BIN"
	case "fromHex":
		fieldOp = "HEX_TO_DEC"
	case "fromBin":
		fieldOp = "BIN_TO_DEC"
	default:
		p.Where.Error("Unknown property access ->%", p.Name)
	}
	var blockType string

	switch fieldOp {
	case "SIN", "COS", "TAN", "ASIN", "ACOS", "ATAN":
		blockType = "math_trig"
	case "RADIANS_TO_DEGREES", "DEGREES_TO_RADIANS":
		blockType = "math_convert_angles"
	case "DEC_TO_HEX", "HEX_TO_DEC", "DEC_TO_BIN", "BIN_TO_DEC":
		blockType = "math_convert_number"
	default:
		blockType = "math_single"
	}

	return blockly.Block{
		Type:   blockType,
		Fields: []blockly.Field{{Name: "OP", Value: fieldOp}},
		Values: []blockly.Value{{Name: "NUM", Block: p.On.Blockly()}},
	}
}
