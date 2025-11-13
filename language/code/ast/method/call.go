package method

import (
	"Falcon/code/ast"
	"Falcon/code/lex"
	"Falcon/code/sugar"
	"strconv"
	"strings"
)

type Call struct {
	Where *lex.Token
	On    ast.Expr
	Name  string
	Args  []ast.Expr
}

type CallSignature struct {
	Module       string
	BlocklyName  string
	YailName     string
	YailArgTypes string
	ParamCount   int
	Consumable   bool
	Signature    ast.Signature
}

func makeSignature(
	module string,
	blocklyName string,
	yailName string,
	yailArgTypes string,
	paramCount int,
	consumable bool,
	signature ast.Signature,
) *CallSignature {
	return &CallSignature{
		Module:       module,
		BlocklyName:  blocklyName,
		YailName:     yailName,
		YailArgTypes: yailArgTypes,
		ParamCount:   paramCount,
		Consumable:   consumable,
		Signature:    signature,
	}
}

var signatures = map[string]*CallSignature{
	"textLen":                 makeSignature("text", "text_length", "string-length", "text", 0, true, ast.SignNumb),
	"trim":                    makeSignature("text", "text_trim", "string-trim", "text", 0, true, ast.SignText),
	"uppercase":               makeSignature("text", "text_changeCase", "string-to-upper-case", "text", 0, true, ast.SignText),
	"lowercase":               makeSignature("text", "text_changeCase", "string-to-lower-case", "text", 0, true, ast.SignText),
	"startsWith":              makeSignature("text", "text_starts_at", "string-starts-at", "text text", 1, true, ast.SignBool),
	"contains":                makeSignature("text", "text_contains", "string-contains", "text text", 1, true, ast.SignBool),
	"containsAny":             makeSignature("text", "text_contains", "string-contains-any", "text list", 1, true, ast.SignBool),
	"containsAll":             makeSignature("text", "text_contains", "string-contains-all", "text list", 1, true, ast.SignBool),
	"split":                   makeSignature("text", "text_split", "string-split", "text text", 1, true, ast.SignList),
	"splitAtFirst":            makeSignature("text", "text_split", "string-split-at-first", "text text", 1, true, ast.SignList),
	"splitAtAny":              makeSignature("text", "text_split", "string-split-at-any", "text list", 1, true, ast.SignList),
	"splitAtFirstOfAny":       makeSignature("text", "text_split", "string-split-at-first-of-any", "text list", 1, true, ast.SignList),
	"splitAtSpaces":           makeSignature("text", "text_split_at_spaces", "string-split-at-spaces", "text", 0, true, ast.SignList),
	"reverse":                 makeSignature("text", "text_reverse", "reverse", "text", 0, true, ast.SignText),
	"csvRowToList":            makeSignature("text", "lists_from_csv_row", "yail-list-from-csv-row", "text", 0, true, ast.SignList),
	"csvTableToList":          makeSignature("text", "lists_from_csv_table", "yail-list-from-csv-table", "text", 0, true, ast.SignList),
	"segment":                 makeSignature("text", "text_segment", "string-substring", "text number number", 2, true, ast.SignText),
	"replace":                 makeSignature("text", "text_replace_all", "string-replace-all", "text text text", 2, true, ast.SignText),
	"replaceFrom":             makeSignature("text", "text_replace_mappings", "string-replace-mappings-dictionary", "text dictionary", 1, true, ast.SignText),
	"replaceFromLongestFirst": makeSignature("text", "text_replace_mappings", "string-replace-mappings-longest-string", "text dictionary", 1, true, ast.SignText),

	"listLen":       makeSignature("list", "lists_length", "yail-list-length", "list", 0, true, ast.SignNumb),
	"add":           makeSignature("list", "lists_add_items", "yail-list-add-to-list!", "", -1, false, ast.SignVoid),
	"containsItem":  makeSignature("list", "lists_is_in", "yail-list-member?", "any list", 1, true, ast.SignBool),
	"indexOf":       makeSignature("list", "lists_position_in", "yail-list-index", "any list", 1, true, ast.SignNumb),
	"insert":        makeSignature("list", "lists_insert_item", "yail-list-insert-item!", "list number any", 2, true, ast.SignVoid),
	"remove":        makeSignature("list", "lists_remove_item", "yail-list-remove-item!", "list-number", 1, false, ast.SignVoid),
	"appendList":    makeSignature("list", "lists_append_list", "yail-list-append!", "list list", 1, false, ast.SignVoid),
	"lookupInPairs": makeSignature("list", "lists_lookup_in_pairs", "yail-alist-lookup", "any list any", 2, true, ast.SignAny),
	"join":          makeSignature("list", "lists_join_with_separator", "yail-list-join-with-separator", "list text", 1, true, ast.SignText),
	"slice":         makeSignature("list", "lists_slice", "yail-list-slice", "list number number", 2, true, ast.SignList),
	"random":        makeSignature("list", "lists_pick_random_item", "yail-list-pick-random", "list", 0, true, ast.SignAny),
	"reverseList":   makeSignature("list", "lists_reverse", "yail-list-reverse", "list", 0, true, ast.SignList),
	"toCsvRow":      makeSignature("list", "lists_to_csv_row", "yail-list-to-csv-row", "list", 0, true, ast.SignText),
	"toCsvTable":    makeSignature("list", "lists_to_csv_table", "yail-list-to-csv-table", "list", 0, true, ast.SignText),
	"sort":          makeSignature("list", "lists_sort", "yail-list-sort", "list", 0, true, ast.SignList),
	"allButFirst":   makeSignature("list", "lists_but_first", "yail-list-but-first", "list", 0, true, ast.SignAny),
	"allButLast":    makeSignature("list", "lists_but_last", "yail-list-but-last", "list", 0, true, ast.SignAny),
	"pairsToDict":   makeSignature("list", "dictionaries_alist_to_dict", "yail-dictionary-alist-to-dict", "list", 0, true, ast.SignDict),

	"dictLen":     makeSignature("dict", "dictionaries_length", "yail-dictionary-length", "dictionary", 0, true, ast.SignNumb),
	"get":         makeSignature("dict", "dictionaries_lookup", "yail-dictionary-lookup", "key dictionary any", 2, true, ast.SignAny),
	"set":         makeSignature("dict", "dictionaries_set_pair", "yail-dictionary-set-pair", "key dictionary any", 2, false, ast.SignVoid),
	"delete":      makeSignature("dict", "dictionaries_delete_pair", "yail-dictionary-delete-pair", "dictionary key", 1, false, ast.SignVoid),
	"getAtPath":   makeSignature("dict", "dictionaries_recursive_lookup", "yail-dictionary-recursive-lookup", "list dictionary any", 2, true, ast.SignAny),
	"setAtPath":   makeSignature("dict", "dictionaries_recursive_set", "yail-dictionary-recursive-set", "list dictionary any", 2, false, ast.SignVoid),
	"containsKey": makeSignature("dict", "dictionaries_is_key_in", "yail-dictionary-is-key-in", "key dictionary", 1, true, ast.SignBool),
	"mergeInto":   makeSignature("dict", "dictionaries_combine_dicts", "yail-dictionary-combine-dicts", "dictionary dictionary", 1, false, ast.SignDict),
	"walkTree":    makeSignature("dict", "dictionaries_walk_tree", "yail-dictionary-walk", "list any", 1, true, ast.SignAny),
	"keys":        makeSignature("dict", "dictionaries_getters", "yail-dictionary-get-keys", "dictionary", 0, true, ast.SignList),
	"values":      makeSignature("dict", "dictionaries_getters", "yail-dictionary-get-values", "dictionary", 0, true, ast.SignList),
	"toPairs":     makeSignature("dict", "dictionaries_dict_to_alist", "yail-dictionary-dict-to-alist", "dictionary", 0, true, ast.SignList),
}

func (c *Call) String() string {
	pFormat := "%.%(%)"
	if !c.On.Continuous() {
		pFormat = "(%).%(%)"
	}
	return sugar.Format(pFormat, c.On.String(), c.Name, ast.JoinExprs(", ", c.Args))
}

func (c *Call) Yail() string {
	signature := c.getCallSignature()

	// handle special cases for arg types
	var argTypes string
	switch signature.BlocklyName {
	case "lists_add_items": // list.add(...) has user defined number of arguments
		argTypes = strings.Repeat("any ", len(c.Args))
	default:
		argTypes = signature.YailArgTypes
	}
	return ast.PrimitiveCall(signature.YailName, c.Name, c.Args, argTypes)
}

func (c *Call) Blockly(flags ...bool) ast.Block {
	signature := c.getCallSignature()
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

func (c *Call) getCallSignature() *CallSignature {
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

func (c *Call) Consumable(flags ...bool) bool {
	signature, ok := signatures[c.Name]
	if !ok {
		c.Where.Error("Cannot find method .%", c.Name)
	}
	return signature.Consumable
}

func (c *Call) Signature() []ast.Signature {
	return []ast.Signature{c.getCallSignature().Signature}
}

func (c *Call) simpleOperand(blockType string, valueName string) ast.Block {
	return ast.Block{Type: blockType, Values: []ast.Value{{Name: valueName, Block: c.On.Blockly()}}}
}
