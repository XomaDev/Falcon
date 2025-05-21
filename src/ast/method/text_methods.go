package method

import "Falcon/ast/blockly"

func (c *Call) textMethods(signature *Signature) blockly.Block {
	switch signature.Name {
	case "text_trim":
		return c.simpleOperand("text_trim", "TEXT")
	case "text_split_at_spaces":
		return c.simpleOperand("text_split_at_spaces", "TEXT")
	case "text_reverse":
		return c.simpleOperand("text_reverse", "VALUE")
	case "lists_from_csv_row":
		return c.simpleOperand("lists_from_csv_row", "TEXT")
	case "lists_from_csv_table":
		return c.simpleOperand("lists_from_csv_table", "TEXT")
	case "text_changeCase":
		return c.textChangeCase()
	case "text_starts_at":
		return c.textStartsAt()
	case "text_contains":
		return c.textContains()
	case "text_split":
		return c.textSplit()
	case "text_segment":
		return c.textSegment()
	case "text_replace_all":
		return c.textReplace()
	case "text_replace_mappings":
		return c.textReplaceFrom()
	default:
		panic("Unknown text method " + signature.Name)
	}
}

func (c *Call) textChangeCase() blockly.Block {
	var fieldOp string
	if c.Name == "uppercase" {
		fieldOp = "UPCASE"
	} else {
		fieldOp = "DOWNCASE"
	}
	return blockly.Block{
		Type:   "text_changeCase",
		Fields: []blockly.Field{{Name: "OP", Value: fieldOp}},
		Values: []blockly.Value{{Name: "TEXT", Block: c.On.Blockly()}},
	}
}

func (c *Call) textReplaceFrom() blockly.Block {
	var fieldOp string
	if c.Name == "replaceFrom" {
		fieldOp = "DICTIONARY_ORDER"
	} else {
		fieldOp = "LONGEST_STRING_FIRST"
	}
	return blockly.Block{
		Type:   "text_replace_mappings",
		Fields: []blockly.Field{{Name: "OP", Value: fieldOp}},
		Values: []blockly.Value{
			{Name: "MAPPINGS", Block: c.Args[0].Blockly()},
			{Name: "TEXT", Block: c.On.Blockly()},
		},
	}
}

func (c *Call) textReplace() blockly.Block {
	return blockly.Block{
		Type: "text_replace_all",
		Values: []blockly.Value{
			{Name: "TEXT", Block: c.On.Blockly()},
			{Name: "SEGMENT", Block: c.Args[0].Blockly()},
			{Name: "REPLACEMENT", Block: c.Args[1].Blockly()},
		},
	}
}

func (c *Call) textSegment() blockly.Block {
	return blockly.Block{
		Type: "text_segment",
		Values: []blockly.Value{
			{Name: "TEXT", Block: c.On.Blockly()},
			{Name: "START", Block: c.Args[0].Blockly()},
			{Name: "LENGTH", Block: c.Args[1].Blockly()},
		},
	}
}

func (c *Call) textSplit() blockly.Block {
	var fieldOp string
	switch c.Name {
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
			{Name: "TEXT", Block: c.On.Blockly()},
			{Name: "AT", Block: c.Args[0].Blockly()},
		},
	}
}

func (c *Call) textContains() blockly.Block {
	var fieldOp string
	switch c.Name {
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
			{Name: "TEXT", Block: c.On.Blockly()},
			{Name: "PIECE", Block: c.Args[0].Blockly()},
		},
	}
}

func (c *Call) textStartsAt() blockly.Block {
	return blockly.Block{
		Type: "text_starts_at",
		Values: []blockly.Value{
			{Name: "TEXT", Block: c.On.Blockly()},
			{Name: "PIECE", Block: c.Args[0].Blockly()},
		},
	}
}
