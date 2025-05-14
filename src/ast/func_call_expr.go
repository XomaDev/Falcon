package ast

import (
	"Falcon/sugar"
	"Falcon/types"
	"strconv"
)

type FuncCall struct {
	Where types.Token
	Name  *string
	Args  []Expr
}

func (f *FuncCall) String() string {
	return sugar.Format("%(%)", *f.Name, JoinExprs(", ", f.Args))
}

func (f *FuncCall) Blockly() Block {
	switch *f.Name {
	case "bin", "octal", "hexa":
		return f.mathRadix()
	case "randInt":
		return f.randInt()
	case "randFloat":
		return f.randFloat()
	case "setRandSeed":
		return f.setRandSeed()
	case "min", "max":
		return f.minOrMax()
	case "avgOf", "maxOf", "minOf", "geoMeanOf", "stdDevOf", "stdErrOf":
		return f.mathOnList()
	case "nodeOf":
		return f.modeOf()
	case "mod", "rem", "quot":
		return f.mathDivide()
	case "aTan2":
		return f.atan2()
	case "formatDecimal":
		return f.formatDecimal()
	default:
		panic("Unimplemented")
	}
}

func (f *FuncCall) formatDecimal() Block {
	f.assertArgLen(2)
	return Block{
		Type:   "math_format_as_decimal",
		Values: MakeValues(f.Args, "NUM", "PLACES"),
	}
}

func (f *FuncCall) atan2() Block {
	f.assertArgLen(2)
	return Block{
		Type:   "math_atan2",
		Values: MakeValues(f.Args, "Y", "X"),
	}
}

func (f *FuncCall) mathDivide() Block {
	f.assertArgLen(2)
	var fieldOp string
	switch *f.Name {
	case "mod":
		fieldOp = "MODULO"
	case "rem":
		fieldOp = "REMAINDER"
	case "quot":
		fieldOp = "QUOTIENT"
	}
	return Block{
		Type:   "math_divide",
		Fields: []Field{{Name: "OP", Value: fieldOp}},
		Values: MakeValues(f.Args, "DIVIDEND", "DIVISOR"),
	}
}

func (f *FuncCall) modeOf() Block {
	f.assertArgLen(1)
	return Block{
		Type:   "math_mode_of_list",
		Values: MakeValues(f.Args, "LIST"),
	}
}

func (f *FuncCall) mathOnList() Block {
	f.assertArgLen(1)
	var fieldOp string
	switch *f.Name {
	case "avgOf":
		fieldOp = "AVG"
	case "maxOf":
		fieldOp = "MAX"
	case "minOf":
		fieldOp = "MIN"
	case "geoMeanOf":
		fieldOp = "GM"
	case "stdDevOf":
		fieldOp = "SD"
	case "stdErrOf":
		fieldOp = "SE"
	}
	return Block{
		Type:   "math_on_list2",
		Fields: []Field{{Name: "OP", Value: fieldOp}},
		Values: MakeValues(f.Args, "LIST"),
	}
}

func (f *FuncCall) minOrMax() Block {
	argSize := len(f.Args)
	if argSize == 0 {
		f.Where.Error("No arguments provided for %()", *f.Name)
	}
	var fieldOp string
	switch *f.Name {
	case "min":
		fieldOp = "MIN"
	case "max":
		fieldOp = "MAX"
	}
	return Block{
		Type:     "math_on_list",
		Fields:   []Field{{Name: "OP", Value: fieldOp}},
		Mutation: &Mutation{ItemCount: argSize},
		Values:   ToValues("NUM", f.Args),
	}
}

func (f *FuncCall) setRandSeed() Block {
	f.assertArgLen(1)
	return Block{
		Type:   "math_random_set_seed",
		Values: MakeValues(f.Args, "NUM"),
	}
}

func (f *FuncCall) randFloat() Block {
	f.assertArgLen(0)
	return Block{Type: "math_random_float"}
}

func (f *FuncCall) randInt() Block {
	f.assertArgLen(2)
	return Block{
		Type:   "math_random_int",
		Values: MakeValues(f.Args, "FROM", "TO"),
	}
}

func (f *FuncCall) mathRadix() Block {
	f.assertArgLen(1)
	var fieldOp string
	switch *f.Name {
	case "bin":
		fieldOp = "BIN"
	case "octal":
		fieldOp = "OCT"
	case "hexa":
		fieldOp = "HEX"
	}
	textExpr, ok := f.Args[0].(*TextExpr)
	if !ok {
		f.Where.Error("Expected a numeric string argument for %v()", *f.Name)
	}
	return Block{
		Type: "math_number_radix",
		Fields: []Field{
			{Name: "OP", Value: fieldOp},
			{Name: "NUM", Value: *textExpr.Content},
		},
	}
}

func (f *FuncCall) assertArgLen(expectLen int) {
	argsLen := len(f.Args)
	if argsLen != expectLen {
		f.Where.Error("Expected % argument for %() but got %", strconv.Itoa(expectLen), *f.Name, strconv.Itoa(argsLen))
	}
}
