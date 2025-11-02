package procedures

import (
	ast2 "Falcon/lang/code/ast"
	"Falcon/lang/code/sugar"
	"strings"
)

type RetProcedure struct {
	Name       string
	Parameters []string
	Result     ast2.Expr
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
	return sugar.Format("func %(%) =\n\t%", v.Name, strings.Join(v.Parameters, ", "), ast2.Pad(v.Result.String()))
}

func (v *RetProcedure) Blockly() ast2.Block {
	return ast2.Block{
		Type:     "procedures_defreturn",
		Mutation: &ast2.Mutation{Args: ast2.ToArgs(v.Parameters)},
		Fields:   append(ast2.ToFields("VAR", v.Parameters), ast2.Field{Name: "NAME", Value: v.Name}),
		Values:   []ast2.Value{{Name: "RETURN", Block: v.Result.Blockly()}},
	}
}

func (v *RetProcedure) Continuous() bool {
	return false
}

func (v *RetProcedure) Consumable() bool {
	return false
}
