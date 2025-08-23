package variables

import (
	blky "Falcon/ast/blockly"
	"strings"
)

type Var struct {
	Names  []string
	Values []blky.Expr
	Body   []blky.Expr
}

func (v *Var) String() string {
	var builder strings.Builder
	builder.WriteString("local(\n")

	var varLines []string
	for i, name := range v.Names {
		varLines = append(varLines, blky.PadDirect(name+" = "+v.Values[i].String()))
	}
	builder.WriteString(strings.Join(varLines, ",\n"))
	builder.WriteString("\n) {\n")
	builder.WriteString(blky.PadBody(v.Body))
	builder.WriteString("}")
	return builder.String()
}

func (v *Var) Blockly() blky.Block {
	return blky.Block{
		Type:       "local_declaration_statement",
		Mutation:   &blky.Mutation{LocalNames: blky.MakeLocalNames(v.Names...)},
		Fields:     blky.ToFields("VAR", v.Names),
		Values:     blky.ValuesByPrefix("DECL", v.Values),
		Statements: []blky.Statement{blky.CreateStatement("STACK", v.Body)},
	}
}

func (v *Var) Continuous() bool {
	return false
}

func (v *Var) Consumable() bool {
	return false
}
