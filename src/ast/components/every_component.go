package components

import (
	blky "Falcon/ast/blockly"
	"Falcon/sugar"
)

type EveryComponent struct {
	Type string
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
