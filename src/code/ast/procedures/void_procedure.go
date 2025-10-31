package procedures

import (
	"Falcon/code/ast"
	"Falcon/code/sugar"
	"strings"
)

type VoidProcedure struct {
	Name       string
	Parameters []string
	Body       []ast.Expr
}

func (v *VoidProcedure) Yail() string {
	yail := "(def ("
	yail += v.Name
	yail += " "
	yail += strings.Join(v.Parameters, "$param_")
	yail += ")"
	yail += ast.PadBody(v.Body)
	yail += ")"
	return yail
}

func (v *VoidProcedure) String() string {
	return sugar.Format("func %(%) {\n%}", v.Name, strings.Join(v.Parameters, ", "), ast.PadBody(v.Body))
}

func (v *VoidProcedure) Blockly() ast.Block {
	var statements []ast.Statement
	if len(v.Body) > 0 {
		statements = []ast.Statement{ast.CreateStatement("STACK", v.Body)}
	}
	return ast.Block{
		Type:       "procedures_defnoreturn",
		Mutation:   &ast.Mutation{Args: ast.ToArgs(v.Parameters)},
		Fields:     append(ast.ToFields("VAR", v.Parameters), ast.Field{Name: "NAME", Value: v.Name}),
		Statements: statements,
	}
}

func (v *VoidProcedure) Continuous() bool {
	return false
}

func (v *VoidProcedure) Consumable() bool {
	return false
}
