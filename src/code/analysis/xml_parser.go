package analysis

import (
	blky "Falcon/code/ast"
	common2 "Falcon/code/ast/common"
	components2 "Falcon/code/ast/components"
	control2 "Falcon/code/ast/control"
	"Falcon/code/ast/fundamentals"
	list2 "Falcon/code/ast/list"
	"Falcon/code/ast/method"
	procedures2 "Falcon/code/ast/procedures"
	variables2 "Falcon/code/ast/variables"
	"Falcon/code/lex"
	"encoding/xml"
	"strconv"
	"strings"
)

type ValueMap struct {
	valueMap map[string]blky.Expr
}

func (v *ValueMap) getUnsafe(name string) blky.Expr {
	return v.valueMap[name]
}

func (v *ValueMap) get(name string) blky.Expr {
	value := v.valueMap[name]
	if value == nil {
		return &common2.EmptySocket{}
	}
	return value
}

type StatementMap struct {
	statementMap map[string][]blky.Expr
}

func (s *StatementMap) getUnsafe(name string) []blky.Expr {
	return s.statementMap[name]
}

func (s *StatementMap) get(name string) []blky.Expr {
	value := s.statementMap[name]
	if value == nil {
		return []blky.Expr{}
	}
	return value
}

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

func (p *XMLParser) singleExpr(block blky.Block) blky.Expr {
	if len(block.Values) == 0 {
		return &common2.EmptySocket{}
	}
	return p.parseBlock(block.Values[0].Block)
}

func (p *XMLParser) parseBlock(block blky.Block) blky.Expr {
	switch block.Type {
	case "controls_if":
		return p.ctrlIf(block)
	case "controls_forRange":
		return p.ctrlForRange(block)
	case "controls_forEach":
		return &control2.Each{
			IName:    block.SingleField(),
			Iterable: p.singleExpr(block),
			Body:     p.optSingleBody(block)}
	case "controls_for_each_dict":
		return p.ctrlForEachDict(block)
	case "controls_while":
		return &control2.While{
			Condition: p.singleExpr(block),
			Body:      p.optSingleBody(block)}
	case "controls_choose":
		return p.ctrlChoose(block)
	case "controls_do_then_return":
		return &control2.Do{Body: p.optSingleBody(block), Result: p.singleExpr(block)}
	case "controls_eval_but_ignore":
		return makeFuncCall("println", p.singleExpr(block))
	case "controls_openAnotherScreen":
		return makeFuncCall("openScreen", p.singleExpr(block))
	case "controls_openAnotherScreenWithStartValue":
		return makeFuncCall("openScreenWithValue", p.singleExpr(block))
	case "controls_getStartValue":
		return makeFuncCall("getStartValue")
	case "controls_closeScreen":
		return makeFuncCall("closeScreen")
	case "controls_closeScreenWithValue":
		return makeFuncCall("closeScreenWithValue", p.singleExpr(block))
	case "controls_closeApplication":
		return makeFuncCall("closeApp")
	case "controls_getPlainStartText":
		return makeFuncCall("getPlainStartText")
	case "controls_closeScreenWithPlainText":
		return makeFuncCall("closeScreenWithPlainText", p.singleExpr(block))
	case "controls_break":
		return &control2.Break{}

	case "logic_boolean", "logic_true", "logic_false":
		return &fundamentals.Boolean{Value: block.SingleField() == "TRUE"}
	case "logic_negate":
		return &fundamentals.Not{Expr: p.singleExpr(block)}
	case "logic_compare", "logic_operation":
		return p.logicExpr(block)

	case "text":
		return &fundamentals.Text{Content: block.SingleField()}
	case "text_join":
		return p.makeBinary("_", p.fromMinVals(block.Values, 1))
	case "text_length":
		return p.makePropCall("textLen", p.singleExpr(block))
	case "text_isEmpty":
		return p.makeQuestion(lex.Text, block, "emptyText")
	case "text_trim":
		return p.makePropCall("trim", p.singleExpr(block))
	case "text_reverse":
		return p.makePropCall("reverse", p.singleExpr(block))
	case "text_split_at_spaces":
		return p.makePropCall("splitAtSpaces", p.singleExpr(block))
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
		return p.makeQuestion(lex.Text, block, "text")

	case "math_number":
		return &fundamentals.Number{Content: block.SingleField()}
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
		return makeFuncCall("setRandSeed", p.singleExpr(block))
	case "math_number_radix":
		return p.mathRadix(block)
	case "math_on_list": // min() and max()
		return makeFuncCall(strings.ToLower(block.SingleField()), p.fromMinVals(block.Values, 1)...)
	case "math_on_list2":
		return p.mathOnList2(block)
	case "math_mode_of_list":
		return makeFuncCall("modeOf", p.singleExpr(block))
	case "math_trig", "math_sin", "math_cos", "math_tan":
		return makeFuncCall(strings.ToLower(block.SingleField()), p.singleExpr(block))
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
		return &fundamentals.List{Elements: p.fromMinVals(block.Values, 0)}
	case "lists_add_items":
		return p.listAddItem(block)
	case "lists_is_in":
		return p.listContainsItem(block)
	case "lists_length":
		return p.makePropCall("listLength", p.singleExpr(block))
	case "lists_is_empty":
		return p.makeQuestion(lex.OpenSquare, block, "emptyList")
	case "lists_pick_random_item":
		return p.makePropCall("random", p.singleExpr(block))
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
		return makeFuncCall("copyList", p.singleExpr(block))
	case "lists_reverse":
		return p.makePropCall("reverseList", p.singleExpr(block))
	case "lists_to_csv_row":
		return p.makePropCall("toCsvRow", p.singleExpr(block))
	case "lists_to_csv_table":
		return p.makePropCall("toCsvTable", p.singleExpr(block))
	case "lists_sort":
		return p.makePropCall("sort", p.singleExpr(block))
	case "lists_is_list":
		return p.makeQuestion(lex.OpenSquare, block, "list")
	case "lists_from_csv_row":
		return p.makePropCall("csvRowToList", p.singleExpr(block))
	case "lists_from_csv_table":
		return p.makePropCall("csvTableToList", p.singleExpr(block))
	case "lists_but_first":
		return p.makePropCall("allButFirst", p.singleExpr(block))
	case "lists_but_last":
		return p.makePropCall("allButLast", p.singleExpr(block))
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
		return &fundamentals.Dictionary{Elements: p.fromMinVals(block.Values, 0)}
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
		return p.makePropCall("dictLen", p.singleExpr(block))
	case "dictionaries_alist_to_dict":
		return p.makePropCall("pairsToDict", p.singleExpr(block))
	case "dictionaries_dict_to_alist":
		return p.makePropCall("toPairs", p.singleExpr(block))
	case "dictionaries_copy":
		return makeFuncCall("copyDict", p.singleExpr(block))
	case "dictionaries_combine_dicts":
		return p.dictCombine(block)
	case "dictionaries_walk_tree":
		return p.dictWalkTree(block)
	case "dictionaries_walk_all":
		return &fundamentals.WalkAll{}
	case "dictionaries_is_dict":
		return p.makeQuestion(lex.OpenCurly, block, "dict")

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
		return makeFuncCall("makeColor", p.singleExpr(block))
	case "color_split_color":
		return makeFuncCall("splitColor", p.singleExpr(block))

	case "global_declaration":
		return &variables2.Global{Name: block.SingleField(), Value: p.singleExpr(block)}
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
		return &fundamentals.Text{Content: block.SingleField()}
	case "helpers_dropdown":
		return &fundamentals.HelperDropdown{Key: block.Mutation.Key, Option: block.SingleField()}

	case "component_component_block":
		return &fundamentals.Component{Name: block.SingleField(), Type: block.Mutation.ComponentType}
	case "component_set_get":
		return p.componentProp(block)
	case "component_event":
		return p.componentEvent(block)
	case "component_method":
		return p.componentMethod(block)
	case "component_all_component_block":
		return &components2.EveryComponent{Type: block.Mutation.ComponentType}
	default:
		panic("Unsupported block type: " + block.Type)
	}
}

func (p *XMLParser) componentMethod(block blky.Block) blky.Expr {
	if block.Mutation.IsGeneric {
		pVals := p.makeValueMap(block.Values)
		var callArgs []blky.Expr

		for i := 0; ; i++ {
			aArg := pVals.getUnsafe("ARG" + strconv.Itoa(i))
			if aArg == nil {
				break
			}
			callArgs = append(callArgs, aArg)
		}
		return &components2.GenericMethodCall{
			Component:     pVals.get("COMPONENT"),
			ComponentType: block.Mutation.ComponentType,
			Method:        block.Mutation.MethodName,
			Args:          callArgs,
		}
	}
	return &components2.MethodCall{
		ComponentName: block.Mutation.InstanceName,
		ComponentType: block.Mutation.ComponentType,
		Method:        block.Mutation.MethodName,
		Args:          p.fromVals(block.Values),
	}
}

func (p *XMLParser) componentEvent(block blky.Block) blky.Expr {
	// TODO: supply parameters to events later
	if block.Mutation.IsGeneric {
		return &components2.GenericEvent{
			ComponentType: block.Mutation.ComponentType,
			Event:         block.Mutation.EventName,
			Parameters:    make([]string, 0),
			Body:          p.optSingleBody(block),
		}
	}
	return &components2.Event{
		ComponentName: block.Mutation.InstanceName,
		ComponentType: block.Mutation.ComponentType,
		Event:         block.Mutation.EventName,
		Parameters:    make([]string, 0),
		Body:          p.optSingleBody(block),
	}
}

func (p *XMLParser) componentProp(block blky.Block) blky.Expr {
	pFields := p.makeFieldMap(block.Fields)
	property := pFields["PROP"]
	isSet := block.Mutation.SetOrGet == "set"

	if block.Mutation.IsGeneric {
		pVals := p.makeValueMap(block.Values)
		if isSet {
			return &components2.GenericPropertySet{
				Component:     pVals.get("COMPONENT"),
				ComponentType: block.Mutation.ComponentType,
				Property:      property,
				Value:         pVals.get("VALUE"),
			}
		}
		return &components2.GenericPropertyGet{
			Component:     pVals.get("COMPONENT"),
			ComponentType: block.Mutation.ComponentType,
			Property:      property,
		}
	}

	if isSet {
		return &components2.PropertySet{
			ComponentName: pFields["COMPONENT_SELECTOR"],
			ComponentType: block.Mutation.ComponentType,
			Property:      property,
			Value:         p.singleExpr(block),
		}
	}
	return &components2.PropertyGet{
		ComponentName: pFields["COMPONENT_SELECTOR"],
		ComponentType: block.Mutation.ComponentType,
		Property:      property,
	}
}

func (p *XMLParser) ctrlChoose(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return &control2.SimpleIf{Condition: pVals.get("TEST"), Then: pVals.get("THENRETURN"), Else: pVals.get("ELSERETURN")}
}

func (p *XMLParser) ctrlForEachDict(block blky.Block) blky.Expr {
	pFields := p.makeFieldMap(block.Fields)
	return &control2.EachPair{
		KeyName:   pFields["KEY"],
		ValueName: pFields["VALUE"],
		Iterable:  p.singleExpr(block),
		Body:      p.optSingleBody(block),
	}
}

func (p *XMLParser) ctrlForRange(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return &control2.For{
		IName: block.SingleField(),
		From:  pVals.get("START"),
		To:    pVals.get("END"),
		By:    pVals.get("STEP"),
		Body:  p.optSingleBody(block),
	}
}

func (p *XMLParser) ctrlIf(block blky.Block) blky.Expr {
	conditions := p.fromVals(block.Values)
	statementMap := p.makeStatementMap(block.Statements)

	var bodies [][]blky.Expr
	elseBody := statementMap.getUnsafe("ELSE")

	for i := range conditions {
		bodies = append(bodies, statementMap.get("DO"+strconv.Itoa(i)))
	}
	return &control2.If{Conditions: conditions, Bodies: bodies, ElseBody: elseBody}
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
	return &procedures2.Call{
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
	return &procedures2.RetProcedure{
		Name:       procedureName,
		Parameters: paramNames,
		Result:     p.singleExpr(block),
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
	return &procedures2.VoidProcedure{
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
		varValues[i] = valueMap.get("DECL" + strconv.Itoa(i))
	}
	if block.GetType() == "local_declaration_statement" {
		return &variables2.Var{
			Names:  varNames,
			Values: varValues,
			Body:   p.optSingleBody(block),
		}
	}
	return &variables2.VarResult{Names: varNames, Values: varValues, Result: valueMap.get("RETURN")}
}

func (p *XMLParser) variableSet(block blky.Block) blky.Expr {
	varName := block.SingleField()
	isGlobal := strings.HasPrefix(varName, "global ")
	if isGlobal {
		varName = varName[len("global "):]
	}
	return p.makeBinary("=",
		[]blky.Expr{
			&variables2.Get{
				Where:  makeFakeToken(lex.Global),
				Global: isGlobal,
				Name:   varName,
			},
			p.singleExpr(block),
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
	return &variables2.Get{Where: makeFakeToken(lex.Global), Global: isGlobal, Name: varName}
}

func (p *XMLParser) dictWalkTree(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("walkTree", pVals.get("DICT"), pVals.get("PATH"))
}

func (p *XMLParser) dictCombine(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("mergeInto", pVals.get("DICT2"), pVals.get("DICT1"))
}

func (p *XMLParser) dictHasKey(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("containsKey", pVals.get("DICT"), pVals.get("KEY"))
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
	return p.makePropCall(pOperation, p.singleExpr(block))
}

func (p *XMLParser) dictSetPath(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("setAtPath", pVals.get("DICT"), pVals.get("KEYS"), pVals.get("VALUE"))
}

func (p *XMLParser) dictLookupPath(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("getAtPath", pVals.get("DICT"), pVals.get("KEYS"), pVals.get("NOTFOUND"))
}

func (p *XMLParser) dictRemove(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("remove", pVals.get("DICT"), pVals.get("KEY"))
}

func (p *XMLParser) dictSet(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("set", pVals.get("KEY"), pVals.get("VALUE"))
}

func (p *XMLParser) dictLookup(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("get", pVals.get("DICT"), pVals.get("KEY"), pVals.get("NOTFOUND"))
}

func (p *XMLParser) dictPair(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makeBinary(":", []blky.Expr{pVals.get("KEY"), pVals.get("VALUE")})
}

func (p *XMLParser) listTransMax(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	pFields := p.makeFieldMap(block.Fields)
	return &list2.Transformer{
		Where:       makeFakeToken(lex.OpenSquare),
		List:        pVals.get("LIST"),
		Name:        "max",
		Args:        []blky.Expr{},
		Names:       []string{pFields["VAR1"], pFields["VAR2"]},
		Transformer: pVals.get("COMPARE"),
	}
}

func (p *XMLParser) listTransMin(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	pFields := p.makeFieldMap(block.Fields)
	return &list2.Transformer{
		Where:       makeFakeToken(lex.OpenSquare),
		List:        pVals.get("LIST"),
		Name:        "min",
		Args:        []blky.Expr{},
		Names:       []string{pFields["VAR1"], pFields["VAR2"]},
		Transformer: pVals.get("COMPARE"),
	}
}

func (p *XMLParser) listSortKeyComparator(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return &list2.Transformer{
		Where:       makeFakeToken(lex.OpenSquare),
		List:        pVals.get("LIST"),
		Name:        "sortByKey",
		Args:        []blky.Expr{},
		Names:       []string{block.SingleField()},
		Transformer: pVals.get("KEY"),
	}
}

func (p *XMLParser) listSortComparator(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	pFields := p.makeFieldMap(block.Fields)
	return &list2.Transformer{
		Where:       makeFakeToken(lex.OpenSquare),
		List:        pVals.get("LIST"),
		Name:        "sort",
		Args:        []blky.Expr{},
		Names:       []string{pFields["VAR1"], pFields["VAR2"]},
		Transformer: pVals.get("COMPARE"),
	}
}

func (p *XMLParser) listReduce(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	pFields := p.makeFieldMap(block.Fields)
	return &list2.Transformer{
		Where:       makeFakeToken(lex.OpenSquare),
		List:        pVals.get("LIST"),
		Name:        "reduce",
		Args:        []blky.Expr{pVals.get("INITANSWER")},
		Names:       []string{pFields["VAR1"], pFields["VAR2"]},
		Transformer: pVals.get("COMBINE"),
	}
}

func (p *XMLParser) listFilter(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return &list2.Transformer{
		Where:       makeFakeToken(lex.OpenSquare),
		List:        pVals.get("LIST"),
		Name:        "filter",
		Args:        []blky.Expr{},
		Names:       []string{block.SingleField()},
		Transformer: pVals.get("TEST"),
	}
}

func (p *XMLParser) listMap(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return &list2.Transformer{
		Where:       makeFakeToken(lex.OpenSquare),
		List:        pVals.get("LIST"),
		Name:        "map",
		Args:        []blky.Expr{},
		Names:       []string{block.SingleField()},
		Transformer: pVals.get("TO"),
	}
}

func (p *XMLParser) listSlice(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("slice", pVals.get("LIST"), pVals.get("INDEX1"), pVals.get("INDEX2"))
}

func (p *XMLParser) listJoin(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("join", pVals.get("LIST"), pVals.get("SEPARATOR"))
}

func (p *XMLParser) listLookupPairs(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("lookupInPairs", pVals.get("LIST"), pVals.get("KEY"), pVals.get("NOTFOUND"))
}

func (p *XMLParser) listRemoveItem(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("remove", pVals.get("LIST"), pVals.get("INDEX"))
}

func (p *XMLParser) listReplaceItem(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return &list2.Set{List: pVals.get("LIST"), Index: pVals.get("NUM"), Value: pVals.get("ITEM")}
}

func (p *XMLParser) listInsertItem(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("insert", pVals.get("LIST"), pVals.get("INDEX"), pVals.get("ITEM"))
}

func (p *XMLParser) listSelectItem(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return &list2.Get{List: pVals.get("LIST"), Index: pVals.get("NUM")}
}

func (p *XMLParser) listIndexOf(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("indexOf", pVals.get("LIST"), pVals.get("ITEM"))
}

func (p *XMLParser) listContainsItem(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("containsItem", pVals.get("LIST"), pVals.get("ITEM"))
}

func (p *XMLParser) listAddItem(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	numElements := block.Mutation.ItemCount
	arrElements := make([]blky.Expr, numElements)
	for i := 0; i < numElements; i++ {
		arrElements[i] = pVals.get("ITEM" + strconv.Itoa(i))
	}
	return p.makePropCall("add", pVals.get("LIST"), arrElements...)
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
	return p.makePropCall(pOperation, pVals.get("TEXT"), pVals.get("MAPPINGS"))
}

func (p *XMLParser) textObfuscate(block blky.Block) blky.Expr {
	return &common2.Transform{
		Where: makeFakeToken(lex.Text),
		On:    &fundamentals.Text{Content: block.SingleField()},
		Name:  "obfuscate"}
}

func (p *XMLParser) textSegment(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("segment", pVals.get("TEXT"), pVals.get("START"), pVals.get("LENGTH"))
}

func (p *XMLParser) textReplace(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("replace", pVals.get("TEXT"), pVals.get("SEGMENT"), pVals.get("REPLACEMENT"))
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
	return p.makePropCall(pOperation, pVals.get("TEXT"), pVals.get("AT"))
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
	return p.makePropCall(pOperation, pVals.get("TEXT"), pVals.get("PIECE"))
}

func (p *XMLParser) textStartsWith(block blky.Block) blky.Expr {
	pVals := p.makeValueMap(block.Values)
	return p.makePropCall("startsWith", pVals.get("TEXT"), pVals.get("PIECE"))
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
	return p.makePropCall(pOperation, p.singleExpr(block))
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
	var funcName string
	switch block.SingleField() {
	case "DEC_TO_HEX":
		funcName = "decToHex"
	case "DEC_TO_BIN":
		funcName = "decToBin"
	case "HEX_TO_DEC":
		funcName = "hexToDec"
	case "BIN_TO_DEC":
		funcName = "binToDec"
	default:
		panic("Unknown MathConvertNumber type: " + block.SingleField())
	}
	return makeFuncCall(funcName, p.singleExpr(block))
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
	return p.makeQuestion(lex.Number, block, question)
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
	funcName := strings.ToLower(block.SingleField())
	switch funcName {
	case "ln":
		funcName = "log"
	case "ceiling":
		funcName = "ceil"
	}
	return makeFuncCall(funcName, p.singleExpr(block))
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
	return makeFuncCall(funcName, p.singleExpr(block))
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
	return makeFuncCall(funcName, &fundamentals.Text{Content: pFields["NUM"]})
}

func (p *XMLParser) mathRandom(block blky.Block) blky.Expr {
	valMap := p.makeValueMap(block.Values)
	return makeFuncCall("randInt", valMap.get("FROM"), valMap.get("TO"))
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
	return &fundamentals.Color{Where: makeFakeToken(lex.Color), Name: name}
}

func (p *XMLParser) makeQuestion(t lex.Type, on blky.Block, name string) blky.Expr {
	return &common2.Question{Where: makeFakeToken(t), On: p.singleExpr(on), Question: name}
}

func (p *XMLParser) makePropCall(name string, on blky.Expr, args ...blky.Expr) blky.Expr {
	return &method.Call{
		Where: makeFakeToken(lex.Text),
		Name:  name,
		On:    on,
		Args:  args,
	}
}

func (p *XMLParser) makeBinary(operator string, operands []blky.Expr) blky.Expr {
	token := makeToken(operator)
	return &common2.BinaryExpr{
		Where:    token,
		Operator: token.Type,
		Operands: operands,
	}
}

func makeFuncCall(name string, args ...blky.Expr) blky.Expr {
	return &common2.FuncCall{
		Where: makeFakeToken(lex.Func),
		Name:  name,
		Args:  args,
	}
}

// TODO: (future) it'll point to something meaningful
func makeFakeToken(t lex.Type) *lex.Token {
	return &lex.Token{
		Column:  -1,
		Row:     -1,
		Context: nil,
		Type:    t,
		Flags:   make([]lex.Flag, 0),
		Content: nil,
	}
}

func makeToken(symbol string) *lex.Token {
	sToken := lex.Symbols[symbol]
	return sToken.Normal(-1, -1, nil, symbol)
}

func (p *XMLParser) optSingleBody(block blky.Block) []blky.Expr {
	if len(block.Statements) > 0 {
		return p.recursiveParse(*block.SingleStatement().Block)
	}
	return []blky.Expr{}
}

func (p *XMLParser) makeStatementMap(allStatements []blky.Statement) StatementMap {
	statementMap := make(map[string][]blky.Expr, len(allStatements))
	for _, stmt := range allStatements {
		statementMap[stmt.Name] = p.recursiveParse(*stmt.Block)
	}
	return StatementMap{statementMap: statementMap}
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

func (p *XMLParser) makeValueMap(allValues []blky.Value) ValueMap {
	valueMap := make(map[string]blky.Expr, len(allValues))
	for _, val := range allValues {
		valueMap[val.Name] = p.parseBlock(val.Block)
	}
	return ValueMap{valueMap: valueMap}
}

func (p *XMLParser) fromVals(allValues []blky.Value) []blky.Expr {
	arrBlocks := make([]blky.Expr, len(allValues))
	for i := range allValues {
		arrBlocks[i] = p.parseBlock(allValues[i].Block)
	}
	return arrBlocks
}

func (p *XMLParser) fromMinVals(allValues []blky.Value, minCount int) []blky.Expr {
	size := max(minCount, len(allValues))
	arrExprs := make([]blky.Expr, size)
	for i := range allValues {
		arrExprs[i] = p.parseBlock(allValues[i].Block)
	}
	for i := len(allValues); i < size; i++ {
		arrExprs[i] = &common2.EmptySocket{}
	}
	return arrExprs
}
