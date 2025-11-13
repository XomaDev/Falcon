package variables

import (
	"Falcon/code/ast"
	"strings"
)

type VarResult struct {
	Names  []string
	Values []ast.Expr
	Result ast.Expr
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
	yail += ast.PadDirect(v.Result.Yail())
	yail += " )"
	return yail
}

func (v *VarResult) String() string {
	var builder strings.Builder
	builder.WriteString("{\n")
	localLines := make([]string, len(v.Names))
	for k, name := range v.Names {
		localLines[k] = "local " + name + " = " + v.Values[k].String()
	}
	builder.WriteString(ast.PadDirect(strings.Join(localLines, "\n")))
	builder.WriteString("\n")
	builder.WriteString(ast.PadDirect(v.Result.String()))
	builder.WriteString("\n}")
	return builder.String()
}

func (v *VarResult) Blockly() ast.Block {
	return ast.Block{
		Type:     "local_declaration_expression",
		Mutation: &ast.Mutation{LocalNames: ast.MakeLocalNames(v.Names...)},
		Fields:   ast.ToFields("VAR", v.Names),
		Values: append(ast.ValuesByPrefix("DECL", v.Values),
			ast.Value{Name: "RETURN", Block: v.Result.Blockly()}),
	}
}

func (v *VarResult) Continuous() bool {
	return true
}

func (v *VarResult) Consumable() bool {
	return true
}

func (v *VarResult) Signature() []ast.Signature {
	return v.Result.Signature()
}
