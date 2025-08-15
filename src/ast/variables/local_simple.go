package variables

import (
	blky "Falcon/ast/blockly"
	"strings"
)

type SimpleVar struct {
	Name  string
	Value blky.Expr
	Body  []blky.Expr
}

func (v *SimpleVar) String() string {
	var builder strings.Builder
	builder.WriteString("local ")
	builder.WriteString(v.Name)
	builder.WriteString(" = ")
	builder.WriteString(v.Value.String())
	builder.WriteString("\n")
	builder.WriteString(blky.JoinExprs("\n", v.Body))
	return builder.String()
}

func (v *SimpleVar) Blockly() blky.Block {
	var statements []blky.Statement
	if len(v.Body) > 0 {
		statements = []blky.Statement{blky.CreateStatement("STACK", v.Body)}
	}
	return blky.Block{
		Type:       "local_declaration_statement",
		Mutation:   &blky.Mutation{LocalNames: blky.MakeLocalNames(v.Name)},
		Fields:     []blky.Field{{Name: "VAR0", Value: v.Name}},
		Values:     []blky.Value{{Name: "DECL0", Block: v.Value.Blockly()}},
		Statements: statements,
		Consumable: false,
	}
}

func (v *SimpleVar) Continuous() bool {
	return false
}
