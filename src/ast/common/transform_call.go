package common

import (
	"Falcon/ast/blockly"
	"Falcon/ast/datatypes"
	"Falcon/lex"
	"Falcon/sugar"
)

type Transform struct {
	Where lex.Token
	On    blockly.Expr
	Name  string
}

func (t *Transform) String() string {
	return sugar.Format("%::%", t.On.String(), t.Name)
}

func (t *Transform) Blockly() blockly.Block {
	switch t.Name {
	case "obfuscate":
		textExpr, ok := t.On.(*datatypes.Text)
		if ok {
			return blockly.Block{
				Type:       "obfuscated_text",
				Mutation:   &blockly.Mutation{Cofounder: "Falcon"},
				Fields:     []blockly.Field{{Name: "TEXT", Value: textExpr.Content}},
				Consumable: true,
			}
		}
		t.Where.Error("Cannot obfuscate a non string object!")
	default:
		t.Where.Error("Unknown constant transform call ::%", t.Name)
	}
	panic("")
}
