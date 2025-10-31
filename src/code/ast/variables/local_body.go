package variables

import (
	"Falcon/code/ast"
	"strings"
)

type Var struct {
	Names  []string
	Values []ast.Expr
	Body   []ast.Expr
}

func (v *Var) Yail() string {
	yail := "(let ( "
	for k, name := range v.Names {
		yail += "($local_"
		yail += name
		yail += " "
		yail += v.Values[k].Yail()
		yail += ") "
	}
	yail += ") "
	yail += ast.PadBodyYail(v.Body)
	yail += ")"
	return yail
}

func (v *Var) String() string {
	var builder strings.Builder
	builder.WriteString("local(\n")

	var varLines []string
	for i, name := range v.Names {
		varLines = append(varLines, ast.PadDirect(name+" = "+v.Values[i].String()))
	}
	builder.WriteString(strings.Join(varLines, ",\n"))
	builder.WriteString("\n) {\n")
	builder.WriteString(ast.PadBody(v.Body))
	builder.WriteString("}")
	return builder.String()
}

func (v *Var) Blockly() ast.Block {
	return ast.Block{
		Type:       "local_declaration_statement",
		Mutation:   &ast.Mutation{LocalNames: ast.MakeLocalNames(v.Names...)},
		Fields:     ast.ToFields("VAR", v.Names),
		Values:     ast.ValuesByPrefix("DECL", v.Values),
		Statements: []ast.Statement{ast.CreateStatement("STACK", v.Body)},
	}
}

func (v *Var) Continuous() bool {
	return false
}

func (v *Var) Consumable() bool {
	return false
}
