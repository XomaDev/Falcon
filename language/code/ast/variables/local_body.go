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
	builder.WriteString("{\n")
	for k, name := range v.Names {
		builder.WriteString("local ")
		builder.WriteString(name)
		builder.WriteString(" = ")
		builder.WriteString(v.Values[k].String())
		builder.WriteString("\n")
	}
	builder.WriteString(ast.PadBody(v.Body))
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

func (v *Var) Signature() []ast.Signature {
	return []ast.Signature{ast.SignVoid}
}
