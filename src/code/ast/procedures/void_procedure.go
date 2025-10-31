package procedures

import (
	"Falcon/code/ast/blockly"
	"Falcon/code/sugar"
	"strings"
)

type VoidProcedure struct {
	Name       string
	Parameters []string
	Body       []blockly.Expr
}

func (v *VoidProcedure) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (v *VoidProcedure) String() string {
	return sugar.Format("func %(%) {\n%}", v.Name, strings.Join(v.Parameters, ", "), blockly.PadBody(v.Body))
}

func (v *VoidProcedure) Blockly() blockly.Block {
	var statements []blockly.Statement
	if len(v.Body) > 0 {
		statements = []blockly.Statement{blockly.CreateStatement("STACK", v.Body)}
	}
	return blockly.Block{
		Type:       "procedures_defnoreturn",
		Mutation:   &blockly.Mutation{Args: blockly.ToArgs(v.Parameters)},
		Fields:     append(blockly.ToFields("VAR", v.Parameters), blockly.Field{Name: "NAME", Value: v.Name}),
		Statements: statements,
	}
}

func (v *VoidProcedure) Continuous() bool {
	return false
}

func (v *VoidProcedure) Consumable() bool {
	return false
}
