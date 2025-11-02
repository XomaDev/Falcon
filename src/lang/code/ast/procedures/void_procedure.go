package procedures

import (
	ast2 "Falcon/lang/code/ast"
	"Falcon/lang/code/sugar"
	"strings"
)

type VoidProcedure struct {
	Name       string
	Parameters []string
	Body       []ast2.Expr
}

func (v *VoidProcedure) Yail() string {
	yail := "(def ("
	yail += v.Name
	yail += " "
	yail += strings.Join(v.Parameters, "$param_")
	yail += ")"
	yail += ast2.PadBodyYail(v.Body)
	yail += ")"
	return yail
}

func (v *VoidProcedure) String() string {
	return sugar.Format("func %(%) {\n%}", v.Name, strings.Join(v.Parameters, ", "), ast2.PadBody(v.Body))
}

func (v *VoidProcedure) Blockly() ast2.Block {
	var statements []ast2.Statement
	if len(v.Body) > 0 {
		statements = []ast2.Statement{ast2.CreateStatement("STACK", v.Body)}
	}
	return ast2.Block{
		Type:       "procedures_defnoreturn",
		Mutation:   &ast2.Mutation{Args: ast2.ToArgs(v.Parameters)},
		Fields:     append(ast2.ToFields("VAR", v.Parameters), ast2.Field{Name: "NAME", Value: v.Name}),
		Statements: statements,
	}
}

func (v *VoidProcedure) Continuous() bool {
	return false
}

func (v *VoidProcedure) Consumable() bool {
	return false
}
