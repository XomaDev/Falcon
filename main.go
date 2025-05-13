package main

import (
	"Falcon/ast"
	"encoding/xml"
)

func main() {
	println("Hello from Falcon!")
	sourceCode := `123 + 456`

	tokens := NewLexer(sourceCode).Lex()
	for _, token := range tokens {
		println(token.String())
	}

	println("\n=== AST ===")

	expressions := NewParser(tokens).ParseAll()
	for _, expression := range expressions {
		println(expression.String())
	}

	println("\n=== Blockly XML ===\n")

	blocks := make([]ast.Block, len(expressions))
	for i, expression := range expressions {
		blocks[i] = expression.Blockly()
	}
	xmlBlock := ast.XmlRoot{
		Blocks: blocks,
		XMLNS:  "https://developers.google.com/blockly/xml",
	}
	bytes, _ := xml.MarshalIndent(xmlBlock, "", "  ")
	println(string(bytes))
}
