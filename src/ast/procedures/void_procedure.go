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
	var statements []blky.Statement
	if len(v.Body) > 0 {
		statements = []blky.Statement{blky.CreateStatement("STACK", v.Body)}
	}
	return blky.Block{
		Type:       "procedures_defnoreturn",
		Mutation:   &blky.Mutation{Args: blky.ToArgs(v.Parameters)},
		Fields:     append(blky.ToFields("VAR", v.Parameters), blky.Field{Name: "NAME", Value: v.Name}),
		Statements: statements,
	}
}

func (v *VoidProcedure) Continuous() bool {
	return false
}

func (v *VoidProcedure) Consumable() bool {
	return false
}
