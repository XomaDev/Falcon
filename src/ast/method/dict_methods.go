package method

import blky "Falcon/ast/blockly"

func (c *Call) dictMethods(signature *Signature) blky.Block {
	switch signature.Name {
	case "dictionaries_length", "dictionaries_dict_to_alist":
		return c.simpleOperand(signature.Name, "DICT")
	case "dictionaries_lookup":
		return c.dictGet()
	case "dictionaries_set_pair":
		return c.dictSet()
	case "dictionaries_delete_pair":
		return c.dictDelete()
	case "dictionaries_recursive_lookup":
		return c.dictGetAtPath()
	case "dictionaries_recursive_set":
		return c.dictSetAtPath()
	case "dictionaries_is_key_in":
		return c.dictContainsKey()
	case "dictionaries_combine_dicts":
		return c.dictMergeInto()
	case "dictionaries_walk_tree":
		return c.dictWalkTree()
	default:
		panic("Unknown text method " + signature.Name)
	}
}

func (c *Call) dictGetters() blky.Block {
	var fieldOp string
	if c.Name == "values" {
		fieldOp = "VALUES"
	} else {
		fieldOp = "KEYS"
	}
	return blky.Block{
		Type:   "dictionaries_getters",
		Fields: []blky.Field{{Name: "OP", Value: fieldOp}},
		Values: []blky.Value{{Name: "DICT", Block: c.On.Blockly()}},
	}
}

func (c *Call) dictWalkTree() blky.Block {
	return blky.Block{
		Type:   "dictionaries_walk_tree",
		Values: blky.MakeValueArgs(c.On, "DICT", c.Args, "PATH"),
	}
}

func (c *Call) dictMergeInto() blky.Block {
	return blky.Block{
		Type: "dictionaries_combine_dicts",
		Values: []blky.Value{
			{Name: "DICT1", Block: c.Args[0].Blockly()},
			{Name: "DICT2", Block: c.On.Blockly()},
		},
	}
}

func (c *Call) dictContainsKey() blky.Block {
	return blky.Block{
		Type:   "dictionaries_is_key_in",
		Values: blky.MakeValueArgs(c.On, "DICT", c.Args, "KEY"),
	}
}

func (c *Call) dictSetAtPath() blky.Block {
	return blky.Block{
		Type:   "dictionaries_recursive_set",
		Values: blky.MakeValueArgs(c.On, "DICT", c.Args, "KEYS", "VALUE"),
	}
}

func (c *Call) dictGetAtPath() blky.Block {
	return blky.Block{
		Type:   "dictionaries_recursive_lookup",
		Values: blky.MakeValueArgs(c.On, "DICT", c.Args, "KEYS", "NOTFOUND"),
	}
}

func (c *Call) dictDelete() blky.Block {
	return blky.Block{
		Type:   "dictionaries_delete_pair",
		Values: blky.MakeValueArgs(c.On, "DICT", c.Args, "KEY"),
	}
}

func (c *Call) dictSet() blky.Block {
	return blky.Block{
		Type:   "dictionaries_set_pair",
		Values: blky.MakeValueArgs(c.On, "DICT", c.Args, "KEY", "VALUE"),
	}
}

func (c *Call) dictGet() blky.Block {
	return blky.Block{
		Type:   "dictionaries_lookup",
		Values: blky.MakeValueArgs(c.On, "DICT", c.Args, "KEY", "NOTFOUND"),
	}
}
