package variables

import (
	ast2 "Falcon/code/ast"
	"strings"
)

type Var struct {
	Names  []string
	Values []ast2.Expr
	Body   []ast2.Expr
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
	yail += ast2.PadBodyYail(v.Body)
	yail += ")"
	return yail
}

func (v *Var) String() string {
	var builder strings.Builder
	builder.WriteString("local(\n")

	var varLines []string
	for i, name := range v.Names {
		varLines = append(varLines, ast2.PadDirect(name+" = "+v.Values[i].String()))
	}
	builder.WriteString(strings.Join(varLines, ",\n"))
	builder.WriteString("\n) {\n")
	builder.WriteString(ast2.PadBody(v.Body))
	builder.WriteString("}")
	return builder.String()
}

func (v *Var) Blockly() ast2.Block {
	return ast2.Block{
		Type:       "local_declaration_statement",
		Mutation:   &ast2.Mutation{LocalNames: ast2.MakeLocalNames(v.Names...)},
		Fields:     ast2.ToFields("VAR", v.Names),
		Values:     ast2.ValuesByPrefix("DECL", v.Values),
		Statements: []ast2.Statement{ast2.CreateStatement("STACK", v.Body)},
	}
}

func (v *Var) Continuous() bool {
	return false
}

func (v *Var) Consumable() bool {
	return false
}
