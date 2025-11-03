package variables

import (
	ast2 "Falcon/code/ast"
	"strings"
)

type VarResult struct {
	Names  []string
	Values []ast2.Expr
	Result ast2.Expr
}

func (v *VarResult) Yail() string {
	yail := "(let ( "
	for k, name := range v.Names {
		yail += "($local_"
		yail += name
		yail += " "
		yail += v.Values[k].Yail()
		yail += ") "
	}
	yail += ") "
	yail += ast2.PadDirect(v.Result.Yail())
	yail += " )"
	return yail
}

func (v *VarResult) String() string {
	var builder strings.Builder
	builder.WriteString("compute(\n")

	var varLines []string
	for i, name := range v.Names {
		varLines = append(varLines, ast2.PadDirect(name+" = "+v.Values[i].String()))
	}
	builder.WriteString(strings.Join(varLines, ",\n"))
	builder.WriteString("\n) -> ")
	builder.WriteString(v.Result.String())
	return builder.String()
}

func (v *VarResult) Blockly() ast2.Block {
	return ast2.Block{
		Type:     "local_declaration_expression",
		Mutation: &ast2.Mutation{LocalNames: ast2.MakeLocalNames(v.Names...)},
		Fields:   ast2.ToFields("VAR", v.Names),
		Values: append(ast2.ValuesByPrefix("DECL", v.Values),
			ast2.Value{Name: "RETURN", Block: v.Result.Blockly()}),
	}
}

func (v *VarResult) Continuous() bool {
	return true
}

func (v *VarResult) Consumable() bool {
	return true
}
