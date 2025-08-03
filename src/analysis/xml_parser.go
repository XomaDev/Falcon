package analysis

import (
	blky "Falcon/ast/blockly"
	"Falcon/ast/common"
	dtypes "Falcon/ast/datatypes"
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
	case "lists_create_with":
		return &dtypes.List{Elements: p.fromMinVals(block.Values, 1)}

	case "math_number":
		return &dtypes.Number{Content: block.SingleField()}
	case "math_compare", "math_bitwise":
		return p.mathExpr(block)
	case "math_add":
		return makeBinary("+", p.fromMinVals(block.Values, 2))
	case "math_subtract":
		return makeBinary("-", p.fromMinVals(block.Values, 2))
	case "math_multiply":
		return makeBinary("*", p.fromMinVals(block.Values, 2))
	case "math_division":
		return makeBinary("/", p.fromMinVals(block.Values, 2))
	case "math_power":
		return makeBinary("^", p.fromMinVals(block.Values, 2))
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
	default:
		panic("Unsupported block type: " + block.Type)
	}
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
	return &common.Question{Where: makeFakeToken(l.Number), Question: question, On: p.parseBlock(block.SingleValue())}
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
	return makeBinary(mathOp, p.fromMinVals(block.Values, 2))
}

func makeBinary(operator string, operands []blky.Expr) blky.Expr {
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
