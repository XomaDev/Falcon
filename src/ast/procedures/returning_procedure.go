package procedures

import (
	blky "Falcon/ast/blockly"
	"Falcon/sugar"
	"strings"
)

type RetProcedure struct {
	Name       string
	Parameters []string
	Result     blky.Expr
}

func (v *RetProcedure) String() string {
	return sugar.Format("func %(%)\n\t=%", v.Name, strings.Join(v.Parameters, ", "), blky.Pad(v.Result))
}

func (v *RetProcedure) Blockly() blky.Block {
	return blky.Block{
		Type:       "procedures_defreturn",
		Mutation:   &blky.Mutation{Args: blky.ToArgs(v.Parameters)},
		Fields:     append(blky.ToFields("VAR", v.Parameters), blky.Field{Name: "NAME", Value: v.Name}),
		Values:     []blky.Value{{Name: "RETURN", Block: v.Result.Blockly()}},
		Consumable: false,
	}
}

func (v *RetProcedure) Continuous() bool {
	return false
}
