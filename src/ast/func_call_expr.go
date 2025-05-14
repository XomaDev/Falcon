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
	case "Bin", "Octal", "Hexa":
		return f.mathRadix()
	case "RandInt":
		return f.randInt()
	case "RandFloat":
		return f.randFloat()
	case "SetRandSeed":
		return f.setRandSeed()
	case "Min", "Max":
		return f.minOrMax()
	default:
		panic("Unimplemented")
	}
}

func (f *FuncCall) minOrMax() Block {
	argSize := len(f.Args)
	if argSize == 0 {
		f.Where.Error("No arguments provided for %()", *f.Name)
	}
	var fieldOp string
	switch *f.Name {
	case "Min":
		fieldOp = "MIN"
	case "Max":
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
	case "Bin":
		fieldOp = "BIN"
	case "Octal":
		fieldOp = "OCT"
	case "Hexa":
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
