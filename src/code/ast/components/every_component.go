package components

import (
	blky "Falcon/code/ast"
	"Falcon/code/sugar"
)

type EveryComponent struct {
	Type string
}

func (e *EveryComponent) Yail() string {
	return "(get-all-components " + e.Type + ")"
}

func (e *EveryComponent) String() string {
	return sugar.Format("every(%)", e.Type)
}

func (e *EveryComponent) Blockly() blky.Block {
	return blky.Block{
		Type:     "component_all_component_block",
		Mutation: &blky.Mutation{ComponentType: e.Type},
		Fields:   []blky.Field{{Name: "COMPONENT_SELECTOR", Value: e.Type}},
	}
}

func (e *EveryComponent) Continuous() bool {
	return true
}

func (e *EveryComponent) Consumable() bool {
	return true
}
