package method

import "Falcon/ast/blockly"

func (m *Call) textMethods(signature *Signature) blockly.Block {
	switch signature.Name {
	case "text_starts_at":
		return m.textStartsAt()
	case "text_contains":
		return m.textContains()
	case "text_split":
		return m.textSplit()
	case "text_segment":
		return m.textSegment()
	case "text_replace_all":
		return m.textReplace()
	case "text_replace_mappings":
		return m.textReplaceFrom()
	default:
		panic("Unknown text method " + signature.Name)
	}
}

func (m *Call) textReplaceFrom() blockly.Block {
	var fieldOp string
	if m.Name == "replaceFrom" {
		fieldOp = "DICTIONARY_ORDER"
	} else {
		fieldOp = "LONGEST_STRING_FIRST"
	}
	return blockly.Block{
		Type:   "text_replace_mappings",
		Fields: []blockly.Field{{Name: "OP", Value: fieldOp}},
		Values: []blockly.Value{
			{Name: "MAPPINGS", Block: m.Args[0].Blockly()},
			{Name: "TEXT", Block: m.On.Blockly()},
		},
	}
}

func (m *Call) textReplace() blockly.Block {
	return blockly.Block{
		Type: "text_replace_all",
		Values: []blockly.Value{
			{Name: "TEXT", Block: m.On.Blockly()},
			{Name: "SEGMENT", Block: m.Args[0].Blockly()},
			{Name: "REPLACEMENT", Block: m.Args[1].Blockly()},
		},
	}
}

func (m *Call) textSegment() blockly.Block {
	return blockly.Block{
		Type: "text_segment",
		Values: []blockly.Value{
			{Name: "TEXT", Block: m.On.Blockly()},
			{Name: "START", Block: m.Args[0].Blockly()},
			{Name: "LENGTH", Block: m.Args[1].Blockly()},
		},
	}
}

func (m *Call) textSplit() blockly.Block {
	var fieldOp string
	switch m.Name {
	case "split":
		fieldOp = "SPLIT"
	case "splitAtFirst":
		fieldOp = "SPLITATFIRST"
	case "splitAtAny":
		fieldOp = "SPLITATANY"
	case "splitAtFirstOfAny":
		fieldOp = "SPLITATFIRSTOFANY"
	}
	return blockly.Block{
		Type:     "text_split",
		Mutation: &blockly.Mutation{Mode: fieldOp},
		Fields:   []blockly.Field{{Name: "OP", Value: fieldOp}},
		Values: []blockly.Value{
			{Name: "TEXT", Block: m.On.Blockly()},
			{Name: "AT", Block: m.Args[0].Blockly()},
		},
	}
}

func (m *Call) textContains() blockly.Block {
	var fieldOp string
	switch m.Name {
	case "contains":
		fieldOp = "CONTAINS"
	case "containsAny":
		fieldOp = "CONTAINS_ANY"
	case "containsAll":
		fieldOp = "CONTAINS_ALL"
	}
	return blockly.Block{
		Type:     "text_contains",
		Mutation: &blockly.Mutation{Mode: fieldOp},
		Fields:   []blockly.Field{{Name: "OP", Value: fieldOp}},
		Values: []blockly.Value{
			{Name: "TEXT", Block: m.On.Blockly()},
			{Name: "PIECE", Block: m.Args[0].Blockly()},
		},
	}
}

func (m *Call) textStartsAt() blockly.Block {
	return blockly.Block{
		Type: "text_starts_at",
		Values: []blockly.Value{
			{Name: "TEXT", Block: m.On.Blockly()},
			{Name: "PIECE", Block: m.Args[0].Blockly()},
		},
	}
}
