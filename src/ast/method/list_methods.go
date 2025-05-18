package method

import "Falcon/ast/blockly"

func (m *Call) listMethods(signature *Signature) blockly.Block {
	switch signature.Name {
	case "lists_add_items":
		return m.listAdd()
	case "lists_is_in":
		return m.listContainsItem()
	case "lists_position_in":
		return m.listIndexOf()
	case "lists_insert_item":
		return m.listInsertItem()
	case "lists_remove_item":
		return m.listRemoveAt()
	case "lists_append_list":
		return m.listAppendList()
	case "lists_lookup_in_pairs":
		return m.listLookupInPairs()
	case "lists_join_with_separator":
		return m.listJoin()
	default:
		panic("Unknown list method " + signature.Name)
	}
}

func (m *Call) listJoin() blockly.Block {
	return blockly.Block{
		Type:   "lists_join_with_separator",
		Values: blockly.MakeValueArgs(m.On, "LIST", m.Args, "SEPARATOR"),
	}
}

func (m *Call) listLookupInPairs() blockly.Block {
	return blockly.Block{
		Type:   "lists_lookup_in_pairs",
		Values: blockly.MakeValueArgs(m.On, "LIST", m.Args, "KEY", "NOTFOUND"),
	}
}

func (m *Call) listAppendList() blockly.Block {
	return blockly.Block{
		Type:   "lists_append_list",
		Values: blockly.MakeValueArgs(m.On, "LIST0", m.Args, "LIST1"),
	}
}

func (m *Call) listRemoveAt() blockly.Block {
	return blockly.Block{
		Type:   "lists_remove_item",
		Values: blockly.MakeValueArgs(m.On, "LIST", m.Args, "INDEX"),
	}
}

func (m *Call) listInsertItem() blockly.Block {
	return blockly.Block{
		Type:   "lists_insert_item",
		Values: blockly.MakeValueArgs(m.On, "LIST", m.Args, "INDEX", "ITEM"),
	}
}

func (m *Call) listIndexOf() blockly.Block {
	return blockly.Block{
		Type:   "lists_position_in",
		Values: blockly.MakeValueArgs(m.On, "LIST", m.Args, "ITEM"),
	}
}

func (m *Call) listContainsItem() blockly.Block {
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
