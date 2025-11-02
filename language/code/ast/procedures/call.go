package procedures

import (
	ast2 "Falcon/code/ast"
	"Falcon/code/sugar"
)

type Call struct {
	Name       string
	Parameters []string
	Arguments  []ast2.Expr
	Returning  bool
}

func (v *Call) Yail() string {
	yail := "((get-var "
	yail += v.Name
	yail += ") "
	yail += ast2.JoinYailExprs(" ", v.Arguments)
	yail += ")"
	return yail
}

func (v *Call) String() string {
	return sugar.Format("%(%)", v.Name, ast2.JoinExprs(", ", v.Arguments))
}

func (v *Call) Blockly() ast2.Block {
	var blockType string
	if v.Returning {
		blockType = "procedures_callreturn"
	} else {
		blockType = "procedures_callnoreturn"
	}
	return ast2.Block{
		Type:     blockType,
		Mutation: &ast2.Mutation{Name: v.Name, Args: ast2.ToArgs(v.Parameters)},
		Fields:   []ast2.Field{{Name: "PROCNAME", Value: v.Name}},
		Values:   ast2.ValuesByPrefix("ARG", v.Arguments),
	}
}

func (v *Call) Continuous() bool {
	return true
}

func (v *Call) Consumable() bool {
	return v.Returning
}
