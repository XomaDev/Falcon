package fundamentals

import (
	"Falcon/code/ast/blockly"
)

type HelperDropdown struct {
	Key    string
	Option string
}

func (h *HelperDropdown) String() string {
	return h.Key + "@" + h.Option
}

func (h *HelperDropdown) Blockly() blockly.Block {
	return blockly.Block{
		Type:     "helpers_dropdown",
		Mutation: &blockly.Mutation{Key: h.Key},
		Fields:   []blockly.Field{{Name: "OPTION", Value: h.Option}},
	}
}

func (h *HelperDropdown) Continuous() bool {
	return true
}

func (h *HelperDropdown) Consumable() bool {
	return true
}
