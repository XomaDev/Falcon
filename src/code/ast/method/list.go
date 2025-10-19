package method

import (
	blockly2 "Falcon/code/ast/blockly"
)

func (c *Call) listMethods(signature *Signature) blockly2.Block {
	switch signature.Name {
	case "lists_length", "lists_pick_random_item", "lists_reverse", "lists_to_csv_row",
		"lists_to_csv_table", "lists_sort", "lists_but_first", "lists_but_last":
		return c.simpleOperand(signature.Name, "LIST")
	case "dictionaries_alist_to_dict":
		return c.simpleOperand(signature.Name, "PAIRS")
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

func (c *Call) listSlice() blockly2.Block {
	return blockly2.Block{
		Type:   "lists_slice",
		Values: blockly2.MakeValueArgs(c.On, "LIST", c.Args, "INDEX1", "INDEX2"),
	}
}

func (c *Call) listJoin() blockly2.Block {
	return blockly2.Block{
		Type:   "lists_join_with_separator",
		Values: blockly2.MakeValueArgs(c.On, "LIST", c.Args, "SEPARATOR"),
	}
}

func (c *Call) listLookupInPairs() blockly2.Block {
	return blockly2.Block{
		Type:   "lists_lookup_in_pairs",
		Values: blockly2.MakeValueArgs(c.On, "LIST", c.Args, "KEY", "NOTFOUND"),
	}
}

func (c *Call) listAppendList() blockly2.Block {
	return blockly2.Block{
		Type:   "lists_append_list",
		Values: blockly2.MakeValueArgs(c.On, "LIST0", c.Args, "LIST1"),
	}
}

func (c *Call) listRemoveAt() blockly2.Block {
	return blockly2.Block{
		Type:   "lists_remove_item",
		Values: blockly2.MakeValueArgs(c.On, "LIST", c.Args, "INDEX"),
	}
}

func (c *Call) listInsertItem() blockly2.Block {
	return blockly2.Block{
		Type:   "lists_insert_item",
		Values: blockly2.MakeValueArgs(c.On, "LIST", c.Args, "INDEX", "ITEM"),
	}
}

func (c *Call) listIndexOf() blockly2.Block {
	return blockly2.Block{
		Type:   "lists_position_in",
		Values: blockly2.MakeValueArgs(c.On, "LIST", c.Args, "ITEM"),
	}
}

func (c *Call) listContainsItem() blockly2.Block {
	return blockly2.Block{
		Type:   "lists_is_in",
		Values: blockly2.MakeValueArgs(c.On, "LIST", c.Args, "ITEM"),
	}
}

func (c *Call) listAdd() blockly2.Block {
	return blockly2.Block{
		Type:     "lists_add_items",
		Mutation: &blockly2.Mutation{ItemCount: len(c.Args)},
		Values:   blockly2.ValueArgsByPrefix(c.On, "LIST", "ITEM", c.Args),
	}
}
