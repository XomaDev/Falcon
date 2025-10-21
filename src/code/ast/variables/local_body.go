package variables

import (
	"Falcon/code/ast/blockly"
	"strings"
)

type Var struct {
	Names  []string
	Values []blockly.Expr
	Body   []blockly.Expr
}

func (v *Var) String() string {
	var builder strings.Builder
	builder.WriteString("local(\n")

	var varLines []string
	for i, name := range v.Names {
		varLines = append(varLines, blockly.PadDirect(name+" = "+v.Values[i].String()))
	}
	builder.WriteString(strings.Join(varLines, ",\n"))
	builder.WriteString("\n) {\n")
	builder.WriteString(blockly.PadBody(v.Body))
	builder.WriteString("}")
	return builder.String()
}

func (v *Var) Blockly() blockly.Block {
	return blockly.Block{
		Type:       "local_declaration_statement",
		Mutation:   &blockly.Mutation{LocalNames: blockly.MakeLocalNames(v.Names...)},
		Fields:     blockly.ToFields("VAR", v.Names),
		Values:     blockly.ValuesByPrefix("DECL", v.Values),
		Statements: []blockly.Statement{blockly.CreateStatement("STACK", v.Body)},
	}
}

func (v *Var) Continuous() bool {
	return false
}

func (v *Var) Consumable() bool {
	return false
}
