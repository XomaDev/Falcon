package common

import (
	ast2 "Falcon/lang/code/ast"
	"Falcon/lang/code/ast/fundamentals"
	"Falcon/lang/code/lex"
	"Falcon/lang/code/sugar"
)

type Transform struct {
	Where *lex.Token
	On    ast2.Expr
	Name  string
}

func (t *Transform) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (t *Transform) String() string {
	return sugar.Format("%::%", t.On.String(), t.Name)
}

func (t *Transform) Blockly() ast2.Block {
	switch t.Name {
	case "obfuscate":
		textExpr, ok := t.On.(*fundamentals.Text)
		if ok {
			return ast2.Block{
				Type:     "obfuscated_text",
				Mutation: &ast2.Mutation{Cofounder: "Falcon"},
				Fields:   []ast2.Field{{Name: "TEXT", Value: textExpr.Content}},
			}
		}
		t.Where.Error("Cannot obfuscate a non string object!")
	default:
		t.Where.Error("Unknown constant transform call ::%", t.Name)
	}
	panic("")
}

func (t *Transform) Continuous() bool {
	return true
}

func (t *Transform) Consumable() bool {
	return false
}
