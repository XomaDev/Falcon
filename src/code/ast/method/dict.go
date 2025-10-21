package method

import (
	"Falcon/code/ast/blockly"
)

func (c *Call) dictMethods(signature *Signature) blockly.Block {
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
	case "dictionaries_getters":
		return c.dictGetters()
	default:
		panic("Unknown text method " + signature.Name)
	}
}

func (c *Call) dictGetters() blockly.Block {
	var fieldOp string
	if c.Name == "values" {
		fieldOp = "VALUES"
	} else {
		fieldOp = "KEYS"
	}
	return blockly.Block{
		Type:   "dictionaries_getters",
		Fields: []blockly.Field{{Name: "OP", Value: fieldOp}},
		Values: []blockly.Value{{Name: "DICT", Block: c.On.Blockly()}},
	}
}

func (c *Call) dictWalkTree() blockly.Block {
	return blockly.Block{
		Type:   "dictionaries_walk_tree",
		Values: blockly.MakeValueArgs(c.On, "DICT", c.Args, "PATH"),
	}
}

func (c *Call) dictMergeInto() blockly.Block {
	return blockly.Block{
		Type: "dictionaries_combine_dicts",
		Values: []blockly.Value{
			{Name: "DICT1", Block: c.Args[0].Blockly()},
			{Name: "DICT2", Block: c.On.Blockly()},
		},
	}
}

func (c *Call) dictContainsKey() blockly.Block {
	return blockly.Block{
		Type:   "dictionaries_is_key_in",
		Values: blockly.MakeValueArgs(c.On, "DICT", c.Args, "KEY"),
	}
}

func (c *Call) dictSetAtPath() blockly.Block {
	return blockly.Block{
		Type:   "dictionaries_recursive_set",
		Values: blockly.MakeValueArgs(c.On, "DICT", c.Args, "KEYS", "VALUE"),
	}
}

func (c *Call) dictGetAtPath() blockly.Block {
	return blockly.Block{
		Type:   "dictionaries_recursive_lookup",
		Values: blockly.MakeValueArgs(c.On, "DICT", c.Args, "KEYS", "NOTFOUND"),
	}
}

func (c *Call) dictDelete() blockly.Block {
	return blockly.Block{
		Type:   "dictionaries_delete_pair",
		Values: blockly.MakeValueArgs(c.On, "DICT", c.Args, "KEY"),
	}
}

func (c *Call) dictSet() blockly.Block {
	return blockly.Block{
		Type:   "dictionaries_set_pair",
		Values: blockly.MakeValueArgs(c.On, "DICT", c.Args, "KEY", "VALUE"),
	}
}

func (c *Call) dictGet() blockly.Block {
	return blockly.Block{
		Type:   "dictionaries_lookup",
		Values: blockly.MakeValueArgs(c.On, "DICT", c.Args, "KEY", "NOTFOUND"),
	}
}
