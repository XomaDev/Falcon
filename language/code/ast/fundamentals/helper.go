package fundamentals

import (
	"Falcon/code/ast"
)

type HelperDropdown struct {
	Key    string
	Option string
}

func (h *HelperDropdown) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (h *HelperDropdown) String() string {
	return h.Key + "@" + h.Option
}

func (h *HelperDropdown) Blockly() ast.Block {
	return ast.Block{
		Type:     "helpers_dropdown",
		Mutation: &ast.Mutation{Key: h.Key},
		Fields:   []ast.Field{{Name: "OPTION", Value: h.Option}},
	}
}

func (h *HelperDropdown) Continuous() bool {
	return true
}

func (h *HelperDropdown) Consumable() bool {
	return true
}

func (h *HelperDropdown) Signature() ast.Signature {
	return ast.SignHelper
}
