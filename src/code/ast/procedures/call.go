package procedures

import (
	"Falcon/code/ast"
	"Falcon/code/sugar"
)

type Call struct {
	Name       string
	Parameters []string
	Arguments  []ast.Expr
	Returning  bool
}

func (v *Call) Yail() string {
	yail := "((get-var "
	yail += v.Name
	yail += ") "
	yail += ast.JoinYailExprs(" ", v.Arguments)
	yail += ")"
	return yail
}

func (v *Call) String() string {
	return sugar.Format("%(%)", v.Name, ast.JoinExprs(", ", v.Arguments))
}

func (v *Call) Blockly() ast.Block {
	var blockType string
	if v.Returning {
		blockType = "procedures_callreturn"
	} else {
		blockType = "procedures_callnoreturn"
	}
	return ast.Block{
		Type:     blockType,
		Mutation: &ast.Mutation{Name: v.Name, Args: ast.ToArgs(v.Parameters)},
		Fields:   []ast.Field{{Name: "PROCNAME", Value: v.Name}},
		Values:   ast.ValuesByPrefix("ARG", v.Arguments),
	}
}

func (v *Call) Continuous() bool {
	return true
}

func (v *Call) Consumable() bool {
	return v.Returning
}
