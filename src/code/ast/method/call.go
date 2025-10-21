package method

import (
	blockly2 "Falcon/code/ast/blockly"
	"Falcon/code/lex"
	"Falcon/code/sugar"
	"strconv"
)

type Call struct {
	Where *lex.Token
	On    blockly2.Expr
	Name  string
	Args  []blockly2.Expr
}

type Signature struct {
	Module     string
	Name       string
	ParamCount int
	Consumable bool
}

func makeSignature(module string, name string, paramCount int, consumable bool) *Signature {
	return &Signature{Module: module, Name: name, ParamCount: paramCount}
}

var signatures = map[string]*Signature{
	"textLen":                 makeSignature("text", "text_length", 0, true),
	"trim":                    makeSignature("text", "text_trim", 0, true),
	"uppercase":               makeSignature("text", "text_changeCase", 0, true),
	"lowercase":               makeSignature("text", "text_changeCase", 0, true),
	"startsWith":              makeSignature("text", "text_starts_at", 1, true),
	"contains":                makeSignature("text", "text_contains", 1, true),
	"containsAny":             makeSignature("text", "text_contains", 1, true),
	"containsAll":             makeSignature("text", "text_contains", 1, true),
	"split":                   makeSignature("text", "text_split", 1, true),
	"splitAtFirst":            makeSignature("text", "text_split", 1, true),
	"splitAtAny":              makeSignature("text", "text_split", 1, true),
	"splitAtFirstOfAny":       makeSignature("text", "text_split", 1, true),
	"splitAtSpaces":           makeSignature("text", "text_split_at_spaces", 0, true),
	"reverse":                 makeSignature("text", "text_reverse", 0, true),
	"csvRowToList":            makeSignature("text", "lists_from_csv_row", 0, true),
	"csvTableToList":          makeSignature("text", "lists_from_csv_table", 0, true),
	"segment":                 makeSignature("text", "text_segment", 2, true),
	"replace":                 makeSignature("text", "text_replace_all", 2, true),
	"replaceFrom":             makeSignature("text", "text_replace_mappings", 1, true),
	"replaceFromLongestFirst": makeSignature("text", "text_replace_mappings", 1, true),

	"listLen":       makeSignature("list", "lists_length", 0, true),
	"add":           makeSignature("list", "lists_add_items", -1, false),
	"containsItem":  makeSignature("list", "lists_is_in", 1, true),
	"indexOf":       makeSignature("list", "lists_position_in", 1, true),
	"insert":        makeSignature("list", "lists_insert_item", 2, true),
	"remove":        makeSignature("list", "lists_remove_item", 1, false),
	"appendList":    makeSignature("list", "lists_append_list", 1, false),
	"lookupInPairs": makeSignature("list", "lists_lookup_in_pairs", 2, true),
	"join":          makeSignature("list", "lists_join_with_separator", 1, true),
	"slice":         makeSignature("list", "lists_slice", 2, true),
	"random":        makeSignature("list", "lists_pick_random_item", 0, true),
	"reverseList":   makeSignature("list", "lists_reverse", 0, true),
	"toCsvRow":      makeSignature("list", "lists_to_csv_row", 0, true),
	"toCsvTable":    makeSignature("list", "lists_to_csv_table", 0, true),
	"sort":          makeSignature("list", "lists_sort", 0, true),
	"allButFirst":   makeSignature("list", "lists_but_first", 0, true),
	"allButLast":    makeSignature("list", "lists_but_last", 0, true),
	"pairsToDict":   makeSignature("list", "dictionaries_alist_to_dict", 0, true),

	"dictLen":     makeSignature("dict", "dictionaries_length", 0, true),
	"get":         makeSignature("dict", "dictionaries_lookup", 2, true),
	"set":         makeSignature("dict", "dictionaries_set_pair", 2, false),
	"delete":      makeSignature("dict", "dictionaries_delete_pair", 1, false),
	"getAtPath":   makeSignature("dict", "dictionaries_recursive_lookup", 2, true),
	"setAtPath":   makeSignature("dict", "dictionaries_recursive_set", 2, false),
	"containsKey": makeSignature("dict", "dictionaries_is_key_in", 1, true),
	"mergeInto":   makeSignature("dict", "dictionaries_combine_dicts", 1, false),
	"walkTree":    makeSignature("dict", "dictionaries_walk_tree", 1, true),
	"keys":        makeSignature("dict", "dictionaries_getters", 0, true),
	"values":      makeSignature("dict", "dictionaries_getters", 0, true),
	"toPairs":     makeSignature("dict", "dictionaries_dict_to_alist", 0, true),
}

func (c *Call) String() string {
	pFormat := "%.%(%)"
	if !c.On.Continuous() {
		pFormat = "(%).%(%)"
	}
	return sugar.Format(pFormat, c.On.String(), c.Name, blockly2.JoinExprs(", ", c.Args))
}

func (c *Call) Blockly() blockly2.Block {
	signature, ok := signatures[c.Name]
	if !ok {
		c.Where.Error("Cannot find method .%", c.Name)
	}
	gotArgLen := len(c.Args)
	if signature.ParamCount >= 0 {
		if signature.ParamCount != gotArgLen {
			c.Where.Error("Expected % args but got % for method .%",
				strconv.Itoa(signature.ParamCount), strconv.Itoa(gotArgLen), c.Name)
		}
	} else {
		minArgs := -signature.ParamCount
		if gotArgLen < minArgs {
			c.Where.Error("Expected at least % args but got only % for method .%",
				strconv.Itoa(minArgs), strconv.Itoa(gotArgLen), c.Name)
		}
	}
	switch signature.Module {
	case "text":
		return c.textMethods(signature)
	case "list":
		return c.listMethods(signature)
	case "dict":
		return c.dictMethods(signature)
	default:
		panic("Unknown module " + signature.Module)
	}
}

func (c *Call) Continuous() bool {
	return true
}

func (c *Call) Consumable() bool {
	signature, ok := signatures[c.Name]
	if !ok {
		c.Where.Error("Cannot find method .%", c.Name)
	}
	return signature.Consumable
}

func (c *Call) simpleOperand(blockType string, valueName string) blockly2.Block {
	return blockly2.Block{Type: blockType, Values: []blockly2.Value{{Name: valueName, Block: c.On.Blockly()}}}
}
