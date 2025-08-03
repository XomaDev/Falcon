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
	case "math_number":
		return &dtypes.Number{Content: block.SingleField()}
	case "math_compare":
		return p.mathCompare(block)
	default:
		panic("Unsupported block type: " + block.Type)
	}
}

func (p *XMLParser) mathCompare(block blky.Block) blky.Expr {
	var pOperation string
	switch block.SingleField() {
	case "EQ":
		pOperation = "=="
	case "NEQ":
		pOperation = "!="
	case "LT":
		pOperation = "<"
	case "LTE":
		pOperation = "<="
	case "GT":
		pOperation = ">"
	case "GTE":
		pOperation = ">="
	default:
		panic("Unsupported MathCompare operation: " + block.SingleField())
	}
	token := makeToken(pOperation)
	return &common.BinaryExpr{
		Where:    token,
		Operator: token.Type,
		Operands: p.fromMinVals(block.Values, 2),
	}
}

func makeToken(symbol string) *l.Token {
	sToken := l.Symbols[symbol]
	return sToken.Normal(-1, -1, nil, symbol)
}

func (p *XMLParser) fromMinVals(allValues []blky.Value, minCount int) []blky.Expr {
	arrExprs := make([]blky.Expr, max(minCount, len(allValues)))
	for i := range allValues {
		arrExprs[i] = p.parseBlock(allValues[i].Block)
	}
	return arrExprs
}
