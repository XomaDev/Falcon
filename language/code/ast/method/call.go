package method

import (
	ast2 "Falcon/code/ast"
	"Falcon/code/lex"
	"Falcon/code/sugar"
	"strconv"
	"strings"
)

type Call struct {
	Where *lex.Token
	On    ast2.Expr
	Name  string
	Args  []ast2.Expr
}

type Signature struct {
	Module       string
	BlocklyName  string
	YailName     string
	YailArgTypes string
	ParamCount   int
	Consumable   bool
}

func makeSignature(
	module string,
	blocklyName string,
	yailName string,
	yailArgTypes string,
	paramCount int,
	consumable bool,
) *Signature {
	return &Signature{
		Module:       module,
		BlocklyName:  blocklyName,
		YailName:     yailName,
		YailArgTypes: yailArgTypes,
		ParamCount:   paramCount,
		Consumable:   consumable,
	}
}

var signatures = map[string]*Signature{
	"textLen":                 makeSignature("text", "text_length", "string-length", "text", 0, true),
	"trim":                    makeSignature("text", "text_trim", "string-trim", "text", 0, true),
	"uppercase":               makeSignature("text", "text_changeCase", "string-to-upper-case", "text", 0, true),
	"lowercase":               makeSignature("text", "text_changeCase", "string-to-lower-case", "text", 0, true),
	"startsWith":              makeSignature("text", "text_starts_at", "string-starts-at", "text text", 1, true),
	"contains":                makeSignature("text", "text_contains", "string-contains", "text text", 1, true),
	"containsAny":             makeSignature("text", "text_contains", "string-contains-any", "text list", 1, true),
	"containsAll":             makeSignature("text", "text_contains", "string-contains-all", "text list", 1, true),
	"split":                   makeSignature("text", "text_split", "string-split", "text text", 1, true),
	"splitAtFirst":            makeSignature("text", "text_split", "string-split-at-first", "text text", 1, true),
	"splitAtAny":              makeSignature("text", "text_split", "string-split-at-any", "text list", 1, true),
	"splitAtFirstOfAny":       makeSignature("text", "text_split", "string-split-at-first-of-any", "text list", 1, true),
	"splitAtSpaces":           makeSignature("text", "text_split_at_spaces", "string-split-at-spaces", "text", 0, true),
	"reverse":                 makeSignature("text", "text_reverse", "reverse", "text", 0, true),
	"csvRowToList":            makeSignature("text", "lists_from_csv_row", "yail-list-from-csv-row", "text", 0, true),
	"csvTableToList":          makeSignature("text", "lists_from_csv_table", "yail-list-from-csv-table", "text", 0, true),
	"segment":                 makeSignature("text", "text_segment", "string-substring", "text number number", 2, true),
	"replace":                 makeSignature("text", "text_replace_all", "string-replace-all", "text text text", 2, true),
	"replaceFrom":             makeSignature("text", "text_replace_mappings", "string-replace-mappings-dictionary", "text dictionary", 1, true),
	"replaceFromLongestFirst": makeSignature("text", "text_replace_mappings", "string-replace-mappings-longest-string", "text dictionary", 1, true),

	"listLen":       makeSignature("list", "lists_length", "yail-list-length", "list", 0, true),
	"add":           makeSignature("list", "lists_add_items", "yail-list-add-to-list!", "", -1, false),
	"containsItem":  makeSignature("list", "lists_is_in", "yail-list-member?", "any list", 1, true),
	"indexOf":       makeSignature("list", "lists_position_in", "yail-list-index", "any list", 1, true),
	"insert":        makeSignature("list", "lists_insert_item", "yail-list-insert-item!", "list number any", 2, true),
	"remove":        makeSignature("list", "lists_remove_item", "yail-list-remove-item!", "list-number", 1, false),
	"appendList":    makeSignature("list", "lists_append_list", "yail-list-append!", "list list", 1, false),
	"lookupInPairs": makeSignature("list", "lists_lookup_in_pairs", "yail-alist-lookup", "any list any", 2, true),
	"join":          makeSignature("list", "lists_join_with_separator", "yail-list-join-with-separator", "list text", 1, true),
	"slice":         makeSignature("list", "lists_slice", "yail-list-slice", "list number number", 2, true),
	"random":        makeSignature("list", "lists_pick_random_item", "yail-list-pick-random", "list", 0, true),
	"reverseList":   makeSignature("list", "lists_reverse", "yail-list-reverse", "list", 0, true),
	"toCsvRow":      makeSignature("list", "lists_to_csv_row", "yail-list-to-csv-row", "list", 0, true),
	"toCsvTable":    makeSignature("list", "lists_to_csv_table", "yail-list-to-csv-table", "list", 0, true),
	"sort":          makeSignature("list", "lists_sort", "yail-list-sort", "list", 0, true),
	"allButFirst":   makeSignature("list", "lists_but_first", "yail-list-but-first", "list", 0, true),
	"allButLast":    makeSignature("list", "lists_but_last", "yail-list-but-last", "list", 0, true),
	"pairsToDict":   makeSignature("list", "dictionaries_alist_to_dict", "yail-dictionary-alist-to-dict", "list", 0, true),

	"dictLen":     makeSignature("dict", "dictionaries_length", "yail-dictionary-length", "dictionary", 0, true),
	"get":         makeSignature("dict", "dictionaries_lookup", "yail-dictionary-lookup", "key dictionary any", 2, true),
	"set":         makeSignature("dict", "dictionaries_set_pair", "yail-dictionary-set-pair", "key dictionary any", 2, false),
	"delete":      makeSignature("dict", "dictionaries_delete_pair", "yail-dictionary-delete-pair", "dictionary key", 1, false),
	"getAtPath":   makeSignature("dict", "dictionaries_recursive_lookup", "yail-dictionary-recursive-lookup", "list dictionary any", 2, true),
	"setAtPath":   makeSignature("dict", "dictionaries_recursive_set", "yail-dictionary-recursive-set", "list dictionary any", 2, false),
	"containsKey": makeSignature("dict", "dictionaries_is_key_in", "yail-dictionary-is-key-in", "key dictionary", 1, true),
	"mergeInto":   makeSignature("dict", "dictionaries_combine_dicts", "yail-dictionary-combine-dicts", "dictionary dictionary", 1, false),
	"walkTree":    makeSignature("dict", "dictionaries_walk_tree", "yail-dictionary-walk", "list any", 1, true),
	"keys":        makeSignature("dict", "dictionaries_getters", "yail-dictionary-get-keys", "dictionary", 0, true),
	"values":      makeSignature("dict", "dictionaries_getters", "yail-dictionary-get-values", "dictionary", 0, true),
	"toPairs":     makeSignature("dict", "dictionaries_dict_to_alist", "yail-dictionary-dict-to-alist", "dictionary", 0, true),
}

func (c *Call) String() string {
	pFormat := "%.%(%)"
	if !c.On.Continuous() {
		pFormat = "(%).%(%)"
	}
	return sugar.Format(pFormat, c.On.String(), c.Name, ast2.JoinExprs(", ", c.Args))
}

func (c *Call) Yail() string {
	signature := c.getSignature()

	// handle special cases for arg types
	var argTypes string
	switch signature.BlocklyName {
	case "lists_add_items": // list.add(...) has user defined number of arguments
		argTypes = strings.Repeat("any ", len(c.Args))
	default:
		argTypes = signature.YailArgTypes
	}
	return ast2.PrimitiveCall(signature.YailName, c.Name, c.Args, argTypes)
}

func (c *Call) Blockly() ast2.Block {
	signature := c.getSignature()
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

func (c *Call) getSignature() *Signature {
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
	return signature
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

func (c *Call) simpleOperand(blockType string, valueName string) ast2.Block {
	return ast2.Block{Type: blockType, Values: []ast2.Value{{Name: valueName, Block: c.On.Blockly()}}}
}
