package variables

import (
	"Falcon/code/ast/blockly"
	"strings"
)

type SimpleVar struct {
	Name  string
	Value blockly.Expr
	Body  []blockly.Expr
}

func (v *SimpleVar) String() string {
	var builder strings.Builder
	builder.WriteString("local ")
	builder.WriteString(v.Name)
	builder.WriteString(" = ")
	builder.WriteString(v.Value.String())
	builder.WriteString("\n")
	builder.WriteString(blockly.JoinExprs("\n", v.Body))
	return builder.String()
}

func (v *SimpleVar) Blockly() blockly.Block {
	var statements []blockly.Statement
	if len(v.Body) > 0 {
		statements = []blockly.Statement{blockly.CreateStatement("STACK", v.Body)}
	}
	return blockly.Block{
		Type:       "local_declaration_statement",
		Mutation:   &blockly.Mutation{LocalNames: blockly.MakeLocalNames(v.Name)},
		Fields:     []blockly.Field{{Name: "VAR0", Value: v.Name}},
		Values:     []blockly.Value{{Name: "DECL0", Block: v.Value.Blockly()}},
		Statements: statements,
	}
}

func (v *SimpleVar) Continuous() bool {
	return false
}

func (v *SimpleVar) Consumable() bool {
	return false
}
