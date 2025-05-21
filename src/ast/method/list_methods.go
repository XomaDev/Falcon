package method

import "Falcon/ast/blockly"

func (c *Call) listMethods(signature *Signature) blockly.Block {
	switch signature.Name {
	case "lists_add_items":
		return c.listAdd()
	case "lists_is_in":
		return c.listContainsItem()
	case "lists_position_in":
		return c.listIndexOf()
	case "lists_insert_item":
		return c.listInsertItem()
	case "lists_remove_item":
		return c.listRemoveAt()
	case "lists_append_list":
		return c.listAppendList()
	case "lists_lookup_in_pairs":
		return c.listLookupInPairs()
	case "lists_join_with_separator":
		return c.listJoin()
	case "lists_slice":
		return c.listSlice()
	default:
		panic("Unknown list method " + signature.Name)
	}
}

func (c *Call) listSlice() blockly.Block {
	return blockly.Block{
		Type:   "lists_slice",
		Values: blockly.MakeValueArgs(c.On, "LIST", c.Args, "INDEX1", "INDEX2"),
	}
}

func (c *Call) listJoin() blockly.Block {
	return blockly.Block{
		Type:   "lists_join_with_separator",
		Values: blockly.MakeValueArgs(c.On, "LIST", c.Args, "SEPARATOR"),
	}
}

func (c *Call) listLookupInPairs() blockly.Block {
	return blockly.Block{
		Type:   "lists_lookup_in_pairs",
		Values: blockly.MakeValueArgs(c.On, "LIST", c.Args, "KEY", "NOTFOUND"),
	}
}

func (c *Call) listAppendList() blockly.Block {
	return blockly.Block{
		Type:   "lists_append_list",
		Values: blockly.MakeValueArgs(c.On, "LIST0", c.Args, "LIST1"),
	}
}

func (c *Call) listRemoveAt() blockly.Block {
	return blockly.Block{
		Type:   "lists_remove_item",
		Values: blockly.MakeValueArgs(c.On, "LIST", c.Args, "INDEX"),
	}
}

func (c *Call) listInsertItem() blockly.Block {
	return blockly.Block{
		Type:   "lists_insert_item",
		Values: blockly.MakeValueArgs(c.On, "LIST", c.Args, "INDEX", "ITEM"),
	}
}

func (c *Call) listIndexOf() blockly.Block {
	return blockly.Block{
		Type:   "lists_position_in",
		Values: blockly.MakeValueArgs(c.On, "LIST", c.Args, "ITEM"),
	}
}

func (c *Call) listContainsItem() blockly.Block {
	return blockly.Block{
		Type:   "lists_is_in",
		Values: blockly.MakeValueArgs(c.On, "LIST", c.Args, "ITEM"),
	}
}

func (c *Call) listAdd() blockly.Block {
	return blockly.Block{
		Type:     "lists_add_items",
		Mutation: &blockly.Mutation{ItemCount: len(c.Args)},
		Values:   blockly.ValueArgsByPrefix(c.On, "LIST", "ITEM", c.Args),
	}
}
