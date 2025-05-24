package common

import (
	"Falcon/ast/blockly"
	"Falcon/lex"
	"Falcon/sugar"
)

var Conversions = map[string][]string{
	"root":     {"math", "ROOT"},
	"abs":      {"math", "ABS"},
	"neg":      {"math", "NEG"},
	"log":      {"math", "LN"},
	"exp":      {"math", "EXP"},
	"round":    {"math", "ROUND"},
	"ceil":     {"math", "CEILING"},
	"floor":    {"math", "FLOOR"},
	"sin":      {"math", "SIN"},
	"cos":      {"math", "COS"},
	"tan":      {"math", "TAN"},
	"asin":     {"math", "ASIN"},
	"acos":     {"math", "ACOS"},
	"atan":     {"math", "ATAN"},
	"degrees":  {"math", "RADIANS_TO_DEGREES"},
	"radians":  {"math", "DEGREES_TO_RADIANS"},
	"hex":      {"math", "DEC_TO_HEX"},
	"bin":      {"math", "DEC_TO_BIN"},
	"parseHex": {"math", "HEX_TO_DEC"},
	"parseBin": {"math", "BIN_TO_DEC"},
}

type Convert struct {
	Where lex.Token
	On    blockly.Expr
	Name  string
}

func (p *Convert) String() string {
	return sugar.Format("%->%", p.On.String(), p.Name)
}

func (p *Convert) Blockly() blockly.Block {
	tags, ok := Conversions[p.Name]
	if !ok {
		p.Where.Error("Unknown conversion access ->%", p.Name)
	}
	switch tags[0] {
	case "math":
		return p.mathConversion(tags[1])
	default:
		panic("Unknown undefined module " + tags[0])
	}
}

func (p *Convert) mathConversion(fieldOp string) blockly.Block {
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
		Type:       blockType,
		Fields:     []blockly.Field{{Name: "OP", Value: fieldOp}},
		Values:     []blockly.Value{{Name: "NUM", Block: p.On.Blockly()}},
		Consumable: true,
	}
}
