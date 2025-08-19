package fundamentals

import (
	"Falcon/ast/blockly"
)

type Component struct {
	Name string
}

func (c *Component) String() string {
	return "@" + c.Name
}

func (c *Component) Blockly() blockly.Block {
	return blockly.Block{
		// TODO: add component_type to Mutation
		Mutation: &blockly.Mutation{InstanceName: c.Name},
		Fields:   []blockly.Field{{Name: "COMPONENT_SELECTOR", Value: c.Name}},
	}
}

func (c *Component) Continuous() bool {
	return true
}
