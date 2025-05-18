package method

import blky "Falcon/ast/blockly"

func (m *Call) dictMethods(signature *Signature) blky.Block {
	switch signature.Name {
	case "dictionaries_lookup":
		return m.dictGet()
	case "dictionaries_set_pair":
		return m.dictSet()
	case "dictionaries_delete_pair":
		return m.dictDelete()
	case "dictionaries_recursive_lookup":
		return m.dictGetAtPath()
	case "dictionaries_recursive_set":
		return m.dictSetAtPath()
	case "dictionaries_is_key_in":
		return m.dictContainsKey()
	case "dictionaries_combine_dicts":
		return m.dictMergeInto()
	case "dictionaries_walk_tree":
		return m.dictWalkTree()
	default:
		panic("Unknown text method " + signature.Name)
	}
}

func (m *Call) dictWalkTree() blky.Block {
	return blky.Block{
		Type:   "dictionaries_walk_tree",
		Values: blky.MakeValueArgs(m.On, "DICT", m.Args, "PATH"),
	}
}

func (m *Call) dictMergeInto() blky.Block {
	return blky.Block{
		Type: "dictionaries_combine_dicts",
		Values: []blky.Value{
			{Name: "DICT1", Block: m.Args[0].Blockly()},
			{Name: "DICT2", Block: m.On.Blockly()},
		},
	}
}

func (m *Call) dictContainsKey() blky.Block {
	return blky.Block{
		Type:   "dictionaries_is_key_in",
		Values: blky.MakeValueArgs(m.On, "DICT", m.Args, "KEY"),
	}
}

func (m *Call) dictSetAtPath() blky.Block {
	return blky.Block{
		Type:   "dictionaries_recursive_set",
		Values: blky.MakeValueArgs(m.On, "DICT", m.Args, "KEYS", "VALUE"),
	}
}

func (m *Call) dictGetAtPath() blky.Block {
	return blky.Block{
		Type:   "dictionaries_recursive_lookup",
		Values: blky.MakeValueArgs(m.On, "DICT", m.Args, "KEYS", "NOTFOUND"),
	}
}

func (m *Call) dictDelete() blky.Block {
	return blky.Block{
		Type:   "dictionaries_delete_pair",
		Values: blky.MakeValueArgs(m.On, "DICT", m.Args, "KEY"),
	}
}

func (m *Call) dictSet() blky.Block {
	return blky.Block{
		Type:   "dictionaries_set_pair",
		Values: blky.MakeValueArgs(m.On, "DICT", m.Args, "KEY", "VALUE"),
	}
}

func (m *Call) dictGet() blky.Block {
	return blky.Block{
		Type:   "dictionaries_lookup",
		Values: blky.MakeValueArgs(m.On, "DICT", m.Args, "KEY", "NOTFOUND"),
	}
}
