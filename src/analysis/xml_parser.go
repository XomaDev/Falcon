package analysis

import (
	blky "Falcon/ast/blockly"
	"Falcon/ast/common"
	"Falcon/ast/control"
	dtypes "Falcon/ast/fundamentals"
	"Falcon/ast/list"
	"Falcon/ast/method"
	"Falcon/ast/procedures"
	"Falcon/ast/variables"
	l "Falcon/lex"
	"encoding/xml"
	"strconv"
	"strings"
)

type XMLParser struct {
	xmlContent string
}

func NewXMLParser(xmlContent string) *XMLParser {
	return &XMLParser{xmlContent: xmlContent}
}

func (p *XMLParser) ParseBlockly() []blky.Expr {
	return p.parseAllBlocks(p.decodeXML())
}

func (p *XMLParser) decodeXML() []blky.Block {
	decoder := xml.NewDecoder(strings.NewReader(p.xmlContent))
	decoder.Strict = false
	decoder.DefaultSpace = ""

	var root blky.XmlRoot
	if err := decoder.Decode(&root); err != nil {
		panic(err)
	}
	return root.Blocks
}

func (p *XMLParser) parseAllBlocks(allBlocks []blky.Block) []blky.Expr {
	var parsedBlocks []blky.Expr
	for i := range allBlocks {
		parsedBlocks = append(parsedBlocks, p.parseBlock(allBlocks[i]))
	}
	return parsedBlocks
}

func (p *XMLParser) parseBlock(block blky.Block) blky.Expr {
	switch block.Type {
	case "controls_if":
		return p.ctrlIf(block)
	case "controls_forRange":
		return p.ctrlForRange(block)
	case "controls_forEach":
		return &control.Each{
			IName:    block.SingleField(),
			Iterable: p.parseBlock(block.SingleValue()),
			Body:     p.optSingleBody(block)}
	case "controls_for_each_dict":
		return p.ctrlForEachDict(block)
	case "controls_while":
		return &control.While{
			Condition: p.parseBlock(block.SingleValue()),
			Body:      p.optSingleBody(block)}
	case "controls_choose":
		return p.ctrlChoose(block)
	case "controls_do_then_return":
		return &control.Do{Body: p.optSingleBody(block), Result: p.parseBlock(block.SingleValue())}
	case "controls_eval_but_ignore":
		return makeFuncCall("println", p.parseBlock(block.SingleValue()))
	case "controls_openAnotherScreen":
		return makeFuncCall("openScreen", p.parseBlock(block.SingleValue()))
	case "controls_openAnotherScreenWithStartValue":
		return makeFuncCall("openScreenWithValue", p.parseBlock(block.SingleValue()))
	case "controls_getStartValue":
		return makeFuncCall("getStartValue")
	case "controls_closeScreen":
		return makeFuncCall("closeScreen")
	case "controls_closeScreenWithValue":
		return makeFuncCall("closeScreenWithValue", p.parseBlock(block.SingleValue()))
	case "controls_closeApplication":
		return makeFuncCall("closeApp")
	case "controls_getPlainStartText":
		return makeFuncCall("getPlainStartText")
	case "controls_closeScreenWithPlainText":
		return makeFuncCall("closeScreenWithPlainText", p.parseBlock(block.SingleValue()))
	case "controls_break":
		return &control.Break{}

	case "logic_boolean", "logic_true", "logic_false":
		return &dtypes.Boolean{Value: block.SingleField() == "TRUE"}
	case "logic_negate":
		return &dtypes.Not{Expr: p.parseBlock(block.SingleValue())}
	case "logic_compare", "logic_operation":
		return p.logicExpr(block)

	case "text":
		return &dtypes.Text{Content: block.SingleField()}
	case "text_join":
		return p.makeBinary("_", p.fromMinVals(block.Values, 1))
	case "text_length":
		return p.makePropCall("textLen", p.parseBlock(block.SingleValue()))
	case "text_isEmpty":
		return p.makeQuestion(l.Text, block.SingleValue(), "emptyText")
	case "text_trim":
		return p.makePropCall("trim", p.parseBlock(block.SingleValue()))
	case "text_reverse":
		return p.makePropCall("reverse", p.parseBlock(block.SingleValue()))
	case "text_split_at_spaces":
		return p.makePropCall("splitAtSpaces", p.parseBlock(block.SingleValue()))
	case "text_compare":
		return p.textCompare(block)
	case "text_changeCase":
		return p.textChangeCase(block)
	case "text_starts_at":
		return p.textStartsWith(block)
	case "text_contains":
		return p.textContains(block)
	case "text_split":
		return p.textSplit(block)
	case "text_segment":
		return p.textSegment(block)
	case "text_replace_all":
		return p.textReplace(block)
	case "obfuscated_text":
		return p.textObfuscate(block)
	case "text_replace_mappings":
		return p.textReplaceMap(block)
	case "text_is_string":
		return p.makeQuestion(l.Text, block.SingleValue(), "text")

	case "math_number":
		return &dtypes.Number{Content: block.SingleField()}
	case "math_compare", "math_bitwise":
		return p.mathExpr(block)
	case "math_add":
		return p.makeBinary("+", p.fromMinVals(block.Values, 2))
	case "math_subtract":
		return p.makeBinary("-", p.fromMinVals(block.Values, 2))
	case "math_multiply":
		return p.makeBinary("*", p.fromMinVals(block.Values, 2))
	case "math_division":
		return p.makeBinary("/", p.fromMinVals(block.Values, 2))
	case "math_power":
		return p.makeBinary("^", p.fromMinVals(block.Values, 2))
	case "math_random_int":
		return p.mathRandom(block)
	case "math_random_float":
		return makeFuncCall("randFloat")
	case "math_random_set_seed":
		return makeFuncCall("setRandSeed", p.parseBlock(block.SingleValue()))
	case "math_number_radix":
		return p.mathRadix(block)
	case "math_on_list": // min() and max()
		return makeFuncCall(strings.ToLower(block.SingleField()), p.fromMinVals(block.Values, 1)...)
	case "math_on_list2":
		return p.mathOnList2(block)
	case "math_mode_of_list":
		return makeFuncCall("modeOf", p.parseBlock(block.SingleValue()))
	case "math_trig", "math_sin", "math_cos", "math_tan":
		return p.mathTrig(block)
	case "math_single":
		return p.mathSingle(block)
	case "math_atan2":
		return makeFuncCall("aTan2", p.fromVals(block.Values)...)
	case "math_format_as_decimal":
		return makeFuncCall("formatDecimal", p.fromMinVals(block.Values, 2)...)
	case "math_divide":
		return p.mathDivide(block)
	case "math_is_a_number":
		return p.mathIsNumber(block)
	case "math_convert_number":
		return p.mathConvertNumber(block)

	case "lists_create_with":
		return &dtypes.List{Elements: p.fromMinVals(block.Values, 0)}
	case "lists_add_items":
		return p.listAddItem(block)
	case "lists_is_in":
		return p.listContainsItem(block)
	case "lists_length":
		return p.makePropCall("listLength", p.parseBlock(block.SingleValue()))
	case "lists_is_empty":
		return p.makeQuestion(l.OpenSquare, block.SingleValue(), "emptyList")
	case "lists_pick_random_item":
		return p.makePropCall("random", p.parseBlock(block.SingleValue()))
	case "lists_position_in":
		return p.listIndexOf(block)
	case "lists_select_item":
		return p.listSelectItem(block)
	case "lists_insert_item":
		return p.listInsertItem(block)
	case "lists_replace_item":
		return p.listReplaceItem(block)
	case "lists_remove_item":
		return p.listRemoveItem(block)
	case "lists_copy":
		return makeFuncCall("copyList", p.parseBlock(block.SingleValue()))
	case "lists_reverse":
		return p.makePropCall("reverseList", p.parseBlock(block.SingleValue()))
	case "lists_to_csv_row":
		return p.makePropCall("toCsvRow", p.parseBlock(block.SingleValue()))
	case "lists_to_csv_table":
		return p.makePropCall("toCsvTable", p.parseBlock(block.SingleValue()))
	case "lists_sort":
		return p.makePropCall("sort", p.parseBlock(block.SingleValue()))
	case "lists_is_list":
		return p.makeQuestion(l.OpenSquare, block.SingleValue(), "list")
	case "lists_from_csv_row":
		return p.makePropCall("csvRowToList", p.parseBlock(block.SingleValue()))
	case "lists_from_csv_table":
		return p.makePropCall("csvTableToList", p.parseBlock(block.SingleValue()))
	case "lists_but_first":
		return p.makePropCall("allButFirst", p.parseBlock(block.SingleValue()))
	case "lists_but_last":
		return p.makePropCall("allButLast", p.parseBlock(block.SingleValue()))
	case "lists_lookup_in_pairs":
		return p.listLookupPairs(block)
	case "lists_join_with_separator":
		return p.listJoin(block)
	case "lists_slice":
		return p.listSlice(block)
	case "lists_map":
		return p.listMap(block)
	case "lists_filter":
		return p.listFilter(block)
	case "lists_reduce":
		return p.listReduce(block)
	case "lists_sort_comparator":
		return p.listSortComparator(block)
	case "lists_sort_key":
		return p.listSortKeyComparator(block)
	case "lists_minimum_value":
		return p.listTransMin(block)
	case "lists_maximum_value":
		return p.listTransMax(block)

	case "pair":
		return p.dictPair(block)
	case "dictionaries_create_with":
		return &dtypes.Dictionary{Elements: p.fromMinVals(block.Values, 0)}
	case "dictionaries_lookup":
		return p.dictLookup(block)
	case "dictionaries_set_pair":
		return p.dictSet(block)
	case "dictionaries_delete_pair":
		return p.dictRemove(block)
	case "dictionaries_recursive_lookup":
		return p.dictLookupPath(block)
	case "dictionaries_recursive_set":
		return p.dictSetPath(block)
	case "dictionaries_getters":
		return p.dictGetters(block)
	case "dictionaries_is_key_in":
		return p.dictHasKey(block)
	case "dictionaries_length":
		return p.makePropCall("dictLen", p.parseBlock(block.SingleValue()))
	case "dictionaries_alist_to_dict":
		return p.makePropCall("pairsToDict", p.parseBlock(block.SingleValue()))
	case "dictionaries_dict_to_alist":
		return p.makePropCall("toPairs", p.parseBlock(block.SingleValue()))
	case "dictionaries_copy":
		return makeFuncCall("copyDict", p.parseBlock(block.SingleValue()))
	case "dictionaries_combine_dicts":
		return p.dictCombine(block)
	case "dictionaries_walk_tree":
		return p.dictWalkTree(block)
	case "dictionaries_walk_all":
		return &dtypes.WalkAll{}
	case "dictionaries_is_dict":
		return p.makeQuestion(l.OpenCurly, block.SingleValue(), "dict")

	case "color_black":
		return p.makeColor("black")
	case "color_white":
		return p.makeColor("white")
	case "color_red":
		return p.makeColor("red")
	case "color_pink":
		return p.makeColor("pink")
	case "color_orange":
		return p.makeColor("orange")
	case "color_yellow":
		return p.makeColor("yellow")
	case "color_green":
		return p.makeColor("green")
	case "color_cyan":
		return p.makeColor("cyan")
	case "color_blue":
		return p.makeColor("blue")
	case "color_magenta":
		return p.makeColor("magenta")
	case "color_light_gray":
		return p.makeColor("light_gray")
	case "color_dark_gray":
		return p.makeColor("dark_gray")
	case "color_make_color":
		return makeFuncCall("makeColor", p.parseBlock(block.SingleValue()))
	case "color_split_color":
		return makeFuncCall("splitColor", p.parseBlock(block.SingleValue()))

	case "global_declaration":
		return &variables.Global{Name: block.SingleField(), Value: p.parseBlock(block.SingleValue())}
	case "lexical_variable_get":
		return p.variableGet(block)
	case "lexical_variable_set":
		return p.variableSet(block)
	case "local_declaration_statement", "local_declaration_expression":
		return p.variableSmts(block)

	case "procedures_defnoreturn":
		return p.voidProcedure(block)
	case "procedures_defreturn":
		return p.returnProcedure(block)
	case "procedures_callnoreturn", "procedures_callreturn":
		return p.procedureCall(block)

	case "helpers_assets":
		return &dtypes.Text{Content: block.SingleField()}
	case "helpers_dropdown":
		return &dtypes.HelperDropdown{Key: block.Mutation.Key, Option: block.SingleField()}

	// TODO: impl component blocks
	default:
		panic("Unsupported block type: " + block.Type)
	}
}

func (p *XMLParser) ctrlChoose(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return &control.SimpleIf{Condition: pVals["TEST"], Then: pVals["THENRETURN"], Else: pVals["ELSERETURN"]}
}

func (p *XMLParser) ctrlForEachDict(block blky.Block) blky.Expr {
	pFields := p.makeFieldMap(block.Fields)
	return &control.EachPair{
		KeyName:   pFields["KEY"],
		ValueName: pFields["VALUE"],
		Iterable:  p.parseBlock(block.SingleValue()),
		Body:      p.optSingleBody(block),
	}
}

func (p *XMLParser) ctrlForRange(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return &control.For{
		IName: block.SingleField(),
		From:  pVals["START"],
		To:    pVals["END"],
		By:    pVals["STEP"],
		Body:  p.optSingleBody(block),
	}
}

func (p *XMLParser) ctrlIf(block blky.Block) blky.Expr {
	conditions := p.fromVals(block.Values)
	var bodies [][]blky.Expr
	var elseBody []blky.Expr
	for _, smt := range block.Statements {
		if strings.HasPrefix(smt.Name, "DO") {
			bodies = append(bodies, p.recursiveParse(*smt.Block))
		} else {
			elseBody = p.recursiveParse(*smt.Block)
		}
	}
	return &control.If{Conditions: conditions, Bodies: bodies, ElseBody: elseBody}
}

func (p *XMLParser) logicExpr(block blky.Block) blky.Expr {
	var pOperation string
	switch block.SingleField() {
	case "EQ":
		pOperation = "=="
	case "NEQ":
		pOperation = "!="
	case "AND":
		pOperation = "&&"
	case "OR":
		pOperation = "||"
	default:
		panic("Unknown Logic Compare operation: " + block.SingleField())
	}
	return p.makeBinary(pOperation, p.fromMinVals(block.Values, 2))
}

func (p *XMLParser) procedureCall(block blky.Block) blky.Expr {
	var mutArgsNames []blky.Arg
	if block.Mutation != nil {
		mutArgsNames = block.Mutation.Args
	}
	paramNames := make([]string, len(mutArgsNames))
	for i := range mutArgsNames {
		paramNames[i] = mutArgsNames[i].Name
	}
	procedureName := block.SingleField()
	args := p.fromVals(block.Values)
	return &procedures.Call{
		Name:       procedureName,
		Parameters: paramNames,
		Arguments:  args,
		Returning:  block.Type == "procedures_callreturn",
	}
}

func (p *XMLParser) returnProcedure(block blky.Block) blky.Expr {
	procedureName := p.makeFieldMap(block.Fields)["NAME"]
	var mutArgs []blky.Arg
	if block.Mutation != nil {
		mutArgs = block.Mutation.Args
	}
	paramNames := make([]string, len(mutArgs))
	for i := range mutArgs {
		paramNames[i] = mutArgs[i].Name
	}
	return &procedures.RetProcedure{
		Name:       procedureName,
		Parameters: paramNames,
		Result:     p.parseBlock(block.SingleValue()),
	}
}

func (p *XMLParser) voidProcedure(block blky.Block) blky.Expr {
	procedureName := p.makeFieldMap(block.Fields)["NAME"]
	var mutArgs []blky.Arg
	if block.Mutation != nil {
		mutArgs = block.Mutation.Args
	}
	paramNames := make([]string, len(mutArgs))
	for i := range mutArgs {
		paramNames[i] = mutArgs[i].Name
	}
	return &procedures.VoidProcedure{
		Name:       procedureName,
		Parameters: paramNames,
		Body:       p.optSingleBody(block),
	}
}

func (p *XMLParser) variableSmts(block blky.Block) blky.Expr {
	numOfVars := len(block.Mutation.LocalNames)
	fieldMap := p.makeFieldMap(block.Fields)
	valueMap := p.makeValueMap(block.Values)

	varNames := make([]string, numOfVars)
	varValues := make([]blky.Expr, numOfVars)

	for i := 0; i < numOfVars; i++ {
		varNames[i] = fieldMap["VAR"+strconv.Itoa(i)]
		varValues[i] = valueMap["DECL"+strconv.Itoa(i)]
	}
	if block.GetType() == "local_declaration_statement" {
		return &variables.Var{
			Names:  varNames,
			Values: varValues,
			Body:   p.optSingleBody(block),
		}
	}
	return &variables.VarResult{Names: varNames, Values: varValues, Result: valueMap["RETURN"]}
}

func (p *XMLParser) variableSet(block blky.Block) blky.Expr {
	varName := block.SingleField()
	isGlobal := strings.HasPrefix(varName, "global ")
	if isGlobal {
		varName = varName[len("global "):]
	}
	return p.makeBinary("=",
		[]blky.Expr{
			&variables.Get{
				Where:  makeFakeToken(l.Global),
				Global: isGlobal,
				Name:   varName,
			},
			p.parseBlock(block.SingleValue()),
		},
	)
}

func (p *XMLParser) variableGet(block blky.Block) blky.Expr {
	varName := block.Fields[0].Name
	if varName == "VAR" {
		varName = block.SingleField()
	}
	isGlobal := strings.HasPrefix(varName, "global ")
	if isGlobal {
		varName = varName[len("global "):]
	}
	return &variables.Get{Where: makeFakeToken(l.Global), Global: isGlobal, Name: varName}
}

func (p *XMLParser) dictWalkTree(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("walkTree", pVals["DICT"], pVals["PATH"])
}

func (p *XMLParser) dictCombine(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("mergeInto", pVals["DICT2"], pVals["DICT1"])
}

func (p *XMLParser) dictHasKey(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("containsKey", pVals["DICT"], pVals["KEY"])
}

func (p *XMLParser) dictGetters(block blky.Block) blky.Expr {
	var pOperation string
	switch block.SingleField() {
	case "KEYS":
		pOperation = "keys"
	case "VALUES":
		pOperation = "values"
	default:
		panic("Unknown DictGetters operation: " + block.SingleField())
	}
	return p.makePropCall(pOperation, p.parseBlock(block.SingleValue()))
}

func (p *XMLParser) dictSetPath(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("setAtPath", pVals["DICT"], pVals["KEYS"], pVals["VALUE"])
}

func (p *XMLParser) dictLookupPath(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("getAtPath", pVals["DICT"], pVals["KEYS"], pVals["NOTFOUND"])
}

func (p *XMLParser) dictRemove(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("remove", pVals["DICT"], pVals["KEY"])
}

func (p *XMLParser) dictSet(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("set", pVals["KEY"], pVals["VALUE"])
}

func (p *XMLParser) dictLookup(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("get", pVals["DICT"], pVals["KEY"], pVals["NOTFOUND"])
}

func (p *XMLParser) dictPair(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makeBinary(":", []blky.Expr{pVals["KEY"], pVals["VALUE"]})
}

func (p *XMLParser) listTransMax(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	pFields := p.makeFieldMap(block.Fields)
	return &list.Transformer{
		Where:       makeFakeToken(l.OpenSquare),
		List:        pVals["LIST"],
		Name:        "max",
		Args:        []blky.Expr{},
		Names:       []string{pFields["VAR1"], pFields["VAR2"]},
		Transformer: pVals["COMPARE"],
	}
}

func (p *XMLParser) listTransMin(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	pFields := p.makeFieldMap(block.Fields)
	return &list.Transformer{
		Where:       makeFakeToken(l.OpenSquare),
		List:        pVals["LIST"],
		Name:        "min",
		Args:        []blky.Expr{},
		Names:       []string{pFields["VAR1"], pFields["VAR2"]},
		Transformer: pVals["COMPARE"],
	}
}

func (p *XMLParser) listSortKeyComparator(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return &list.Transformer{
		Where:       makeFakeToken(l.OpenSquare),
		List:        pVals["LIST"],
		Name:        "sortByKey",
		Args:        []blky.Expr{},
		Names:       []string{block.SingleField()},
		Transformer: pVals["KEY"],
	}
}

func (p *XMLParser) listSortComparator(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	pFields := p.makeFieldMap(block.Fields)
	return &list.Transformer{
		Where:       makeFakeToken(l.OpenSquare),
		List:        pVals["LIST"],
		Name:        "sort",
		Args:        []blky.Expr{},
		Names:       []string{pFields["VAR1"], pFields["VAR2"]},
		Transformer: pVals["COMPARE"],
	}
}

func (p *XMLParser) listReduce(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	pFields := p.makeFieldMap(block.Fields)
	return &list.Transformer{
		Where:       makeFakeToken(l.OpenSquare),
		List:        pVals["LIST"],
		Name:        "reduce",
		Args:        []blky.Expr{pVals["INITANSWER"]},
		Names:       []string{pFields["VAR1"], pFields["VAR2"]},
		Transformer: pVals["COMBINE"],
	}
}

func (p *XMLParser) listFilter(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return &list.Transformer{
		Where:       makeFakeToken(l.OpenSquare),
		List:        pVals["LIST"],
		Name:        "filter",
		Args:        []blky.Expr{},
		Names:       []string{block.SingleField()},
		Transformer: pVals["TEST"],
	}
}

func (p *XMLParser) listMap(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return &list.Transformer{
		Where:       makeFakeToken(l.OpenSquare),
		List:        pVals["LIST"],
		Name:        "map",
		Args:        []blky.Expr{},
		Names:       []string{block.SingleField()},
		Transformer: pVals["TO"],
	}
}

func (p *XMLParser) listSlice(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("slice", pVals["LIST"], pVals["INDEX1"], pVals["INDEX2"])
}

func (p *XMLParser) listJoin(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("join", pVals["LIST"], pVals["SEPARATOR"])
}

func (p *XMLParser) listLookupPairs(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("lookupInPairs", pVals["LIST"], pVals["KEY"], pVals["NOTFOUND"])
}

func (p *XMLParser) listRemoveItem(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("remove", pVals["LIST"], pVals["INDEX"])
}

func (p *XMLParser) listReplaceItem(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return &list.Set{List: pVals["LIST"], Index: pVals["NUM"], Value: pVals["ITEM"]}
}

func (p *XMLParser) listInsertItem(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("insert", pVals["LIST"], pVals["INDEX"], pVals["ITEM"])
}

func (p *XMLParser) listSelectItem(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return &list.Get{List: pVals["LIST"], Index: pVals["NUM"]}
}

func (p *XMLParser) listIndexOf(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("indexOf", pVals["LIST"], pVals["ITEM"])
}

func (p *XMLParser) listContainsItem(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("containsItem", pVals["LIST"], pVals["ITEM"])
}

func (p *XMLParser) listAddItem(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	numElements := block.Mutation.ItemCount
	arrElements := make([]blky.Expr, numElements)
	for i := 0; i < numElements; i++ {
		arrElements[i] = pVals["ITEM"+strconv.Itoa(i)]
	}
	return p.makePropCall("add", pVals["LIST"], arrElements...)
}

func (p *XMLParser) textReplaceMap(block blky.Block) blky.Expr {
	var pOperation string
	switch block.SingleField() {
	case "LONGEST_STRING_FIRST":
		pOperation = "replaceFromLongestFirst"
	case "DICTIONARY_ORDER":
		pOperation = "replaceFrom"
	default:
		panic("Unknown Text Replace Map operation: " + block.SingleField())
	}
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall(pOperation, pVals["TEXT"], pVals["MAPPINGS"])
}

func (p *XMLParser) textObfuscate(block blky.Block) blky.Expr {
	return &common.Transform{
		Where: makeFakeToken(l.Text),
		On:    &dtypes.Text{Content: block.SingleField()},
		Name:  "obfuscate"}
}

func (p *XMLParser) textSegment(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("segment", pVals["TEXT"], pVals["START"], pVals["LENGTH"])
}

func (p *XMLParser) textReplace(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("replace", pVals["TEXT"], pVals["SEGMENT"], pVals["REPLACEMENT"])
}

func (p *XMLParser) textSplit(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	var pOperation string
	switch block.SingleField() {
	case "SPLIT":
		pOperation = "split"
	case "SPLITATFIRST":
		pOperation = "splitAtFirst"
	case "SPLITATANY":
		pOperation = "splitAtAny"
	case "SPLITATFIRSTOFANY":
		pOperation = "splitAtFirstOfAny"
	default:
		panic("Unsupported Text Split operation: " + block.SingleField())
	}
	return p.makePropCall(pOperation, pVals["TEXT"], pVals["AT"])
}

func (p *XMLParser) textContains(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	var pOperation string
	switch block.SingleField() {
	case "CONTAINS":
		pOperation = "contains"
	case "CONTAINS_ANY":
		pOperation = "containsAny"
	case "CONTAINS_ALL":
		pOperation = "containsAll"
	default:
		panic("Unsupported Text Contains operation: " + block.SingleField())
	}
	return p.makePropCall(pOperation, pVals["TEXT"], pVals["PIECE"])
}

func (p *XMLParser) textStartsWith(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("startsWith", pVals["TEXT"], pVals["PIECE"])
}

func (p *XMLParser) textChangeCase(block blky.Block) blky.Expr {
	var pOperation string
	switch block.SingleField() {
	case "UPCASE":
		pOperation = "uppercase"
	case "DOWNCASE":
		pOperation = "lowercase"
	default:
		panic("Unsupported Text Change Case operation type: " + block.SingleField())
	}
	return p.makePropCall(pOperation, p.parseBlock(block.SingleValue()))
}

func (p *XMLParser) textCompare(block blky.Block) blky.Expr {
	var pOperation string
	switch block.SingleField() {
	case "EQ":
		pOperation = "==="
	case "NEQ":
		pOperation = "!=="
	case "LT":
		pOperation = "<<"
	case "GT":
		pOperation = ">>"
	default:
		panic("Unknown Text Compare operation: " + block.SingleField())
	}
	return p.makeBinary(pOperation, p.fromMinVals(block.Values, 2))
}

func (p *XMLParser) mathConvertNumber(block blky.Block) blky.Expr {
	var opConvert string
	switch block.SingleField() {
	case "DEC_TO_HEX":
		opConvert = "hex"
	case "DEC_TO_BIN":
		opConvert = "bin"
	case "HEX_TO_DEC":
		opConvert = "fromHex"
	case "BIN_TO_DEC":
		opConvert = "fromBin"
	default:
		panic("Unknown MathConvertNumber type: " + block.SingleField())
	}
	return &common.Convert{Where: makeFakeToken(l.Number), Name: opConvert, On: p.parseBlock(block.SingleValue())}
}

func (p *XMLParser) mathTrig(block blky.Block) blky.Expr {
	return &common.Convert{
		Where: makeFakeToken(l.Number),
		Name:  strings.ToLower(block.SingleField()),
		On:    p.parseBlock(block.SingleValue())}
}

func (p *XMLParser) mathIsNumber(block blky.Block) blky.Expr {
	var question string
	switch block.SingleField() {
	case "NUMBER":
		question = "number"
	case "BINARY":
		question = "bin"
	case "HEXADECIMAL":
		question = "hexa"
	case "BASE10":
		question = "base10"
	default:
		panic("Unknown MathIsNumber type: " + block.SingleField())
	}
	return p.makeQuestion(l.Number, block.SingleValue(), question)
}

func (p *XMLParser) mathDivide(block blky.Block) blky.Expr {
	var funcName string
	switch block.SingleField() {
	case "MODULO":
		funcName = "mod"
	case "REMAINDER":
		funcName = "rem"
	case "QUOTIENT":
		funcName = "quot"
	default:
		panic("Unsupported math divide type: " + block.SingleField())
	}
	return makeFuncCall(funcName, p.fromMinVals(block.Values, 2)...)
}

func (p *XMLParser) mathSingle(block blky.Block) blky.Expr {
	mathOp := strings.ToLower(block.SingleField())
	switch mathOp {
	case "ln":
		mathOp = "log"
	case "ceiling":
		mathOp = "ceil"
	}
	return &common.Convert{Where: makeFakeToken(l.Number), Name: mathOp, On: p.parseBlock(block.SingleValue())}
}

func (p *XMLParser) mathOnList2(block blky.Block) blky.Expr {
	var funcName string
	switch block.SingleField() {
	case "AVG":
		funcName = "avgOf"
	case "MIN":
		funcName = "minOf"
	case "MAX":
		funcName = "maxOf"
	case "GM":
		funcName = "geoMeanOf"
	case "SD":
		funcName = "stdDevOf"
	case "SE":
		funcName = "stdErrOf"
	default:
		panic("Unsupported math on list operation: " + block.SingleField())
	}
	return makeFuncCall(funcName, p.parseBlock(block.SingleValue()))
}

func (p *XMLParser) mathRadix(block blky.Block) blky.Expr {
	pFields := p.makeFieldMap(block.Fields)
	var funcName string
	switch pFields["OP"] {
	case "DEC":
		funcName = "dec"
	case "BIN":
		funcName = "bin"
	case "HEX":
		funcName = "hexa"
	case "OCT":
		funcName = "octal"
	default:
		panic("Unknown Math Radix Type: " + pFields["OP"])
	}
	return makeFuncCall(funcName, &dtypes.Text{Content: pFields["NUM"]})
}

func (p *XMLParser) mathRandom(block blky.Block) blky.Expr {
	valMap := p.makeValueMap(block.Values)
	return makeFuncCall("randInt", valMap["FROM"], valMap["TO"])
}

func (p *XMLParser) mathExpr(block blky.Block) blky.Expr {
	var mathOp string
	switch block.SingleField() {
	case "EQ":
		mathOp = "=="
	case "NEQ":
		mathOp = "!="
	case "LT":
		mathOp = "<"
	case "LTE":
		mathOp = "<="
	case "GT":
		mathOp = ">"
	case "GTE":
		mathOp = ">="
	case "BITAND":
		mathOp = "&"
	case "BITOR":
		mathOp = "|"
	case "BITXOR":
		mathOp = "~"
	default:
		panic("Unsupported math expression operation: " + block.SingleField())
	}
	return p.makeBinary(mathOp, p.fromMinVals(block.Values, 2))
}

func (p *XMLParser) makeColor(name string) blky.Expr {
	return &dtypes.Color{Where: makeFakeToken(l.Color), Name: name}
}

func (p *XMLParser) makeQuestion(t l.Type, on blky.Block, name string) blky.Expr {
	return &common.Question{Where: makeFakeToken(t), On: p.parseBlock(on), Question: name}
}

func (p *XMLParser) makePropCall(name string, on blky.Expr, args ...blky.Expr) blky.Expr {
	return &method.Call{
		Where: makeFakeToken(l.Text),
		Name:  name,
		On:    on,
		Args:  args,
	}
}

func (p *XMLParser) makeBinary(operator string, operands []blky.Expr) blky.Expr {
	token := makeToken(operator)
	return &common.BinaryExpr{
		Where:    token,
		Operator: token.Type,
		Operands: operands,
	}
}

func makeFuncCall(name string, args ...blky.Expr) blky.Expr {
	return &common.FuncCall{
		Where: makeFakeToken(l.Func),
		Name:  name,
		Args:  args,
	}
}

// TODO: (future) it'll point to something meaningful
func makeFakeToken(t l.Type) *l.Token {
	return &l.Token{
		Column:  -1,
		Row:     -1,
		Context: nil,
		Type:    t,
		Flags:   make([]l.Flag, 0),
		Content: nil,
	}
}

func makeToken(symbol string) *l.Token {
	sToken := l.Symbols[symbol]
	return sToken.Normal(-1, -1, nil, symbol)
}

func (p *XMLParser) optSingleBody(block blky.Block) []blky.Expr {
	if len(block.Statements) > 0 {
		return p.recursiveParse(*block.SingleStatement().Block)
	}
	return []blky.Expr{}
}

func (p *XMLParser) recursiveParse(currBlock blky.Block) []blky.Expr {
	var pParsed []blky.Expr
	for {
		pParsed = append(pParsed, p.parseBlock(currBlock))
		if currBlock.Next == nil {
			break
		}
		currBlock = *currBlock.Next.Block
	}
	return pParsed
}

func (p *XMLParser) makeFieldMap(allFields []blky.Field) map[string]string {
	fieldMap := make(map[string]string, len(allFields))
	for _, fil := range allFields {
		fieldMap[fil.Name] = fil.Value
	}
	return fieldMap
}

func (p *XMLParser) makeValueMap(allValues []blky.Value) map[string]blky.Expr {
	valueMap := make(map[string]blky.Expr, len(allValues))
	for _, val := range allValues {
		valueMap[val.Name] = p.parseBlock(val.Block)
	}
	return valueMap
}

func (p *XMLParser) fromVals(allValues []blky.Value) []blky.Expr {
	arrBlocks := make([]blky.Expr, len(allValues))
	for i := range allValues {
		arrBlocks[i] = p.parseBlock(allValues[i].Block)
	}
	return arrBlocks
}

func (p *XMLParser) fromMinVals(allValues []blky.Value, minCount int) []blky.Expr {
	arrExprs := make([]blky.Expr, max(minCount, len(allValues)))
	for i := range allValues {
		arrExprs[i] = p.parseBlock(allValues[i].Block)
	}
	return arrExprs
}
