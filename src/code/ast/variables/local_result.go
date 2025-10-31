package variables

import (
	"Falcon/code/ast/blockly"
	"strings"
)

type VarResult struct {
	Names  []string
	Values []blockly.Expr
	Result blockly.Expr
}

func (v *VarResult) String() string {
	var builder strings.Builder
	builder.WriteString("compute(\n")

	var varLines []string
	for i, name := range v.Names {
		varLines = append(varLines, blockly.PadDirect(name+" = "+v.Values[i].String()))
	}
	builder.WriteString(strings.Join(varLines, ",\n"))
	builder.WriteString("\n) -> ")
	builder.WriteString(v.Result.String())
	return builder.String()
}

func (v *VarResult) Blockly() blockly.Block {
	return blockly.Block{
		Type:     "local_declaration_expression",
		Mutation: &blockly.Mutation{LocalNames: blockly.MakeLocalNames(v.Names...)},
		Fields:   blockly.ToFields("VAR", v.Names),
		Values: append(blockly.ValuesByPrefix("DECL", v.Values),
			blockly.Value{Name: "RETURN", Block: v.Result.Blockly()}),
	}
}

func (v *VarResult) Continuous() bool {
	return true
}

func (v *VarResult) Consumable() bool {
	return true
}
