package procedures

import (
	blky "Falcon/ast/blockly"
	"Falcon/sugar"
	"strings"
)

type VoidProcedure struct {
	Name       string
	Parameters []string
	Body       []blky.Expr
}

func (v *VoidProcedure) String() string {
	return sugar.Format("func %(%) {\n%}", v.Name, strings.Join(v.Parameters, ", "), blky.PadBody(v.Body))
}

func (v *VoidProcedure) Blockly() blky.Block {
	return blky.Block{
		Type:       "procedures_defnoreturn",
		Mutation:   &blky.Mutation{Args: blky.ToArgs(v.Parameters)},
		Fields:     append(blky.ToFields("VAR", v.Parameters), blky.Field{Name: "NAME", Value: v.Name}),
		Statements: []blky.Statement{blky.CreateStatement("STACK", v.Body)},
	}
}
