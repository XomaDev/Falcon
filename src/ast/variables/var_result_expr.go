package variables

import (
	blky "Falcon/ast/blockly"
	"strings"
)

type VarResult struct {
	Names  []string
	Values []blky.Expr
	Result blky.Expr
}

func (v *VarResult) String() string {
	var builder strings.Builder
	builder.WriteString("var(\n")

	var varLines []string
	for i, name := range v.Names {
		varLines = append(varLines, blky.PadDirect(name+" = "+v.Values[i].String()))
	}
	builder.WriteString(strings.Join(varLines, "\n"))
	builder.WriteString("\n) -> ")
	builder.WriteString(v.Result.String())
	return builder.String()
}

func (v *VarResult) Blockly() blky.Block {
	return blky.Block{
		Type:     "local_declaration_expression",
		Mutation: &blky.Mutation{LocalNames: blky.MakeLocalNames(v.Names...)},
		Fields:   blky.ToFields("VAR", v.Names),
		Values: append(blky.ValuesByPrefix("DECL", v.Values),
			blky.Value{Name: "RETURN", Block: v.Result.Blockly()}),
		Consumable: true,
	}
}
