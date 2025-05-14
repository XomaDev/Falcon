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
	default:
		panic("Unimplemented")
	}
}

func (f *FuncCall) mathRadix() Block {
	var fieldOp string
	switch *f.Name {
	case "Bin":
		fieldOp = "BIN"
	case "Octal":
		fieldOp = "OCT"
	case "Hexa":
		fieldOp = "HEX"
	}
	argsLen := len(f.Args)
	if argsLen != 1 {
		f.Where.Error("Expected 1 argument for %() but got %", *f.Name, strconv.Itoa(argsLen))
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
