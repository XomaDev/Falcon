package method

import "Falcon/ast/blockly"

func (m *Call) listMethods(signature *Signature) blockly.Block {
	switch signature.Name {
	case "lists_add_items":
		return m.listAdd()
	case "lists_is_in":
		return m.containsItem()
	default:
		panic("Unknown list method " + signature.Name)
	}
}

func (m *Call) containsItem() blockly.Block {
	return blockly.Block{
		Type:   "lists_is_in",
		Values: blockly.MakeValueArgs(m.On, "LIST", m.Args, "ITEM"),
	}
}

func (m *Call) listAdd() blockly.Block {
	return blockly.Block{
		Type:     "lists_add_items",
		Mutation: &blockly.Mutation{ItemCount: len(m.Args)},
		Values:   blockly.ValueArgsByPrefix(m.On, "LIST", "ITEM", m.Args),
	}
}
