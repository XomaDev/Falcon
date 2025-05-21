package method

import (
	"Falcon/ast/blockly"
	"Falcon/lex"
	"Falcon/sugar"
	"strconv"
)

type Call struct {
	Where lex.Token
	On    blockly.Expr
	Name  string
	Args  []blockly.Expr
}

type Signature struct {
	Module     string
	Name       string
	ParamCount int
}

func makeSignature(module string, name string, paramCount int) *Signature {
	return &Signature{Module: module, Name: name, ParamCount: paramCount}
}

var signatures = map[string]*Signature{
	"textLen":                 makeSignature("text", "text_length", 0),
	"trim":                    makeSignature("text", "text_trim", 0),
	"uppercase":               makeSignature("text", "text_changeCase", 0),
	"lowercase":               makeSignature("text", "text_changeCase", 0),
	"startsWith":              makeSignature("text", "text_starts_at", 1),
	"contains":                makeSignature("text", "text_contains", 1),
	"containsAny":             makeSignature("text", "text_contains", 1),
	"containsAll":             makeSignature("text", "text_contains", 1),
	"split":                   makeSignature("text", "text_split", 1),
	"splitAtFirst":            makeSignature("text", "text_split", 1),
	"splitAtAny":              makeSignature("text(", "text_split", 1),
	"splitAtFirstOfAny":       makeSignature("text", "text_split", 1),
	"splitAtSpaces":           makeSignature("text", "text_split_at_spaces", 0),
	"reverse":                 makeSignature("text", "text_reverse", 0),
	"csvRowToList":            makeSignature("text", "lists_from_csv_row", 0),
	"csvTableToList":          makeSignature("text", "lists_from_csv_table", 0),
	"segment":                 makeSignature("text", "text_segment", 2),
	"replace":                 makeSignature("text", "text_replace_all", 2),
	"replaceFrom":             makeSignature("text", "text_replace_mappings", 1),
	"replaceFromLongestFirst": makeSignature("text", "text_replace_mappings", 1),

	"listLen":       makeSignature("list", "lists_length", 0),
	"add":           makeSignature("list", "lists_add_items", -1),
	"containsItem":  makeSignature("list", "lists_is_in", 1),
	"indexOf":       makeSignature("list", "lists_position_in", 1),
	"insert":        makeSignature("list", "lists_insert_item", 2),
	"remove":        makeSignature("list", "lists_remove_item", 1),
	"appendList":    makeSignature("list", "lists_append_list", 1),
	"lookupInPairs": makeSignature("list", "lists_lookup_in_pairs", 2),
	"join":          makeSignature("list", "lists_join_with_separator", 1),
	"slice":         makeSignature("list", "lists_slice", 2),
	"random":        makeSignature("list", "lists_pick_random_item", 0),
	"reverseList":   makeSignature("list", "lists_reverse", 0),
	"toCsvRow":      makeSignature("list", "lists_to_csv_row", 0),
	"toCsvTable":    makeSignature("list", "lists_to_csv_table", 0),
	"sort":          makeSignature("list", "lists_sort", 0),
	"allButFirst":   makeSignature("list", "lists_but_first", 0),
	"allButLast":    makeSignature("list", "lists_but_last", 0),
	"pairsToDict":   makeSignature("list", "dictionaries_alist_to_dict", 0),

	"get":         makeSignature("dict", "dictionaries_lookup", 2),
	"set":         makeSignature("dict", "dictionaries_set_pair", 2),
	"delete":      makeSignature("dict", "dictionaries_delete_pair", 1),
	"getAtPath":   makeSignature("dict", "dictionaries_recursive_lookup", 2),
	"setAtPath":   makeSignature("dict", "dictionaries_recursive_set", 2),
	"containsKey": makeSignature("dict", "dictionaries_is_key_in", 1),
	"mergeInto":   makeSignature("dict", "dictionaries_combine_dicts", 1),
	"walkTree":    makeSignature("dict", "dictionaries_walk_tree", 1),
}

func (c *Call) String() string {
	return sugar.Format("%.%(%)", c.On.String(), c.Name, blockly.JoinExprs(", ", c.Args))
}

func (c *Call) Blockly() blockly.Block {
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

func (c *Call) simpleOperand(blockType string, valueName string) blockly.Block {
	return blockly.Block{Type: blockType, Values: []blockly.Value{{Name: valueName, Block: c.On.Blockly()}}}
}
