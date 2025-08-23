package procedures

import (
	blky "Falcon/ast/blockly"
	"Falcon/sugar"
)

type Call struct {
	Name       string
	Parameters []string
	Arguments  []blky.Expr
	Returning  bool
}

func (v *Call) String() string {
	return sugar.Format("%(%)", v.Name, blky.JoinExprs(", ", v.Arguments))
}

func (v *Call) Blockly() blky.Block {
	var blockType string
	if v.Returning {
		blockType = "procedures_callreturn"
	} else {
		blockType = "procedures_callnoreturn"
	}
	return blky.Block{
		Type:     blockType,
		Mutation: &blky.Mutation{Name: v.Name, Args: blky.ToArgs(v.Parameters)},
		Fields:   []blky.Field{{Name: "PROCNAME", Value: v.Name}},
		Values:   blky.ValuesByPrefix("ARG", v.Arguments),
	}
}

func (v *Call) Continuous() bool {
	return true
}

func (v *Call) Consumable() bool {
	return v.Returning
}
