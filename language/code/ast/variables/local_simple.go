package variables

import (
	"Falcon/code/ast"
	"strings"
)

type SimpleVar struct {
	Name  string
	Value ast.Expr
	Body  []ast.Expr
}

func (v *SimpleVar) Yail() string {
	yail := "(let ($local_"
	yail += v.Name
	yail += " "
	yail += v.Value.Yail()
	yail += "))"
	yail += ast.PadBodyYail(v.Body)
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
	builder.WriteString(ast.JoinExprs("\n", v.Body))
	return builder.String()
}

func (v *SimpleVar) Blockly(flags ...bool) ast.Block {
	var statements []ast.Statement
	if len(v.Body) > 0 {
		statements = []ast.Statement{ast.CreateStatement("STACK", v.Body)}
	}
	return ast.Block{
		Type:       "local_declaration_statement",
		Mutation:   &ast.Mutation{LocalNames: ast.MakeLocalNames(v.Name)},
		Fields:     []ast.Field{{Name: "VAR0", Value: v.Name}},
		Values:     []ast.Value{{Name: "DECL0", Block: v.Value.Blockly(false)}},
		Statements: statements,
	}
}

func (v *SimpleVar) Continuous() bool {
	return false
}

func (v *SimpleVar) Consumable(flags ...bool) bool {
	return false
}

func (v *SimpleVar) Signature() []ast.Signature {
	return []ast.Signature{ast.SignVoid}
}
