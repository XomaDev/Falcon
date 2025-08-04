package analysis

import (
	blky "Falcon/ast/blockly"
	"Falcon/ast/common"
	dtypes "Falcon/ast/datatypes"
	"Falcon/ast/method"
	l "Falcon/lex"
	"encoding/xml"
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
		return &dtypes.List{Elements: p.fromMinVals(block.Values, 1)}

	case "dictionaries_create_with":
		return &dtypes.Dictionary{Elements: p.fromMinVals(block.Values, 1)}
	default:
		panic("Unsupported block type: " + block.Type)
	}
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
