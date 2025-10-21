package procedures

import (
	"Falcon/code/ast/blockly"
	"Falcon/code/sugar"
	"strings"
)

type RetProcedure struct {
	Name       string
	Parameters []string
	Result     blockly.Expr
}

func (v *RetProcedure) String() string {
	return sugar.Format("func %(%) =\n\t%", v.Name, strings.Join(v.Parameters, ", "), blockly.Pad(v.Result))
}

func (v *RetProcedure) Blockly() blockly.Block {
	return blockly.Block{
		Type:     "procedures_defreturn",
		Mutation: &blockly.Mutation{Args: blockly.ToArgs(v.Parameters)},
		Fields:   append(blockly.ToFields("VAR", v.Parameters), blockly.Field{Name: "NAME", Value: v.Name}),
		Values:   []blockly.Value{{Name: "RETURN", Block: v.Result.Blockly()}},
	}
}

func (v *RetProcedure) Continuous() bool {
	return false
}

func (v *RetProcedure) Consumable() bool {
	return false
}
