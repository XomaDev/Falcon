package procedures

import (
	"Falcon/code/ast/blockly"
	"Falcon/code/sugar"
)

type Call struct {
	Name       string
	Parameters []string
	Arguments  []blockly.Expr
	Returning  bool
}

func (v *Call) String() string {
	return sugar.Format("%(%)", v.Name, blockly.JoinExprs(", ", v.Arguments))
}

func (v *Call) Blockly() blockly.Block {
	var blockType string
	if v.Returning {
		blockType = "procedures_callreturn"
	} else {
		blockType = "procedures_callnoreturn"
	}
	return blockly.Block{
		Type:     blockType,
		Mutation: &blockly.Mutation{Name: v.Name, Args: blockly.ToArgs(v.Parameters)},
		Fields:   []blockly.Field{{Name: "PROCNAME", Value: v.Name}},
		Values:   blockly.ValuesByPrefix("ARG", v.Arguments),
	}
}

func (v *Call) Continuous() bool {
	return true
}

func (v *Call) Consumable() bool {
	return v.Returning
}
