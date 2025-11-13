package procedures

import (
	"Falcon/code/ast"
	"Falcon/code/sugar"
	"strings"
)

type RetProcedure struct {
	Name       string
	Parameters []string
	Result     ast.Expr
}

func (v *RetProcedure) Yail() string {
	yail := "(def ("
	yail += v.Name
	yail += " "
	yail += strings.Join(v.Parameters, "$param_")
	yail += ") "
	yail += v.Result.Yail()
	yail += ")"
	return yail
}

func (v *RetProcedure) String() string {
	return sugar.Format("func %(%) =\n%", v.Name, strings.Join(v.Parameters, ", "), ast.Pad(v.Result.String()))
}

func (v *RetProcedure) Blockly(flags ...bool) ast.Block {
	return ast.Block{
		Type:     "procedures_defreturn",
		Mutation: &ast.Mutation{Args: ast.ToArgs(v.Parameters)},
		Fields:   append(ast.ToFields("VAR", v.Parameters), ast.Field{Name: "NAME", Value: v.Name}),
		Values:   []ast.Value{{Name: "RETURN", Block: v.Result.Blockly()}},
	}
}

func (v *RetProcedure) Continuous() bool {
	return false
}

func (v *RetProcedure) Consumable() bool {
	return false
}

func (v *RetProcedure) Signature() []ast.Signature {
	return v.Result.Signature()
}
