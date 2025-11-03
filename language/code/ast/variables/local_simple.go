package variables

import (
	ast2 "Falcon/code/ast"
	"strings"
)

type SimpleVar struct {
	Name  string
	Value ast2.Expr
	Body  []ast2.Expr
}

func (v *SimpleVar) Yail() string {
	yail := "(let ($local_"
	yail += v.Name
	yail += " "
	yail += v.Value.Yail()
	yail += "))"
	yail += ast2.PadBodyYail(v.Body)
	yail += ")"
	return yail
}

func (v *SimpleVar) String() string {
	var builder strings.Builder
	builder.WriteString("local ")
	builder.WriteString(v.Name)
	builder.WriteString(" = ")
	builder.WriteString(v.Value.String())
	builder.WriteString("\n")
	builder.WriteString(ast2.JoinExprs("\n", v.Body))
	return builder.String()
}

func (v *SimpleVar) Blockly() ast2.Block {
	var statements []ast2.Statement
	if len(v.Body) > 0 {
		statements = []ast2.Statement{ast2.CreateStatement("STACK", v.Body)}
	}
	return ast2.Block{
		Type:       "local_declaration_statement",
		Mutation:   &ast2.Mutation{LocalNames: ast2.MakeLocalNames(v.Name)},
		Fields:     []ast2.Field{{Name: "VAR0", Value: v.Name}},
		Values:     []ast2.Value{{Name: "DECL0", Block: v.Value.Blockly()}},
		Statements: statements,
	}
}

func (v *SimpleVar) Continuous() bool {
	return false
}

func (v *SimpleVar) Consumable() bool {
	return false
}
