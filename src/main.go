//go:build !js && !wasm

package main

import (
	"Falcon/analysis"
	"Falcon/ast/blockly"
	"Falcon/context"
	"Falcon/lex"
	"encoding/xml"
	"os"
)

func main() {
	println("Hello from Falcon!")

	fileName := "hi.mist"
	filePath := "/home/ekina/GolandProjects/Falcon/" + fileName
	codeBytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	sourceCode := string(codeBytes)
	codeContext := &context.CodeContext{SourceCode: &sourceCode, FileName: fileName}

	// lexical analysis
	tokens := lex.NewLexer(codeContext).Lex()
	for _, token := range tokens {
		println(token.Debug())
	}

	println("\n=== AST ===\n")

	// conversion of Falcon -> Blockly XML
	expressions := analysis.NewLangParser(tokens).ParseAll()
	for _, expression := range expressions {
		println(expression.String())
	}

	println("\n=== Blockly XML ===\n")

	blocks := make([]blockly.Block, len(expressions))
	for i, expression := range expressions {
		blocks[i] = expression.Blockly()
	}
	xmlBlock := blockly.XmlRoot{
		Blocks: blocks,
		XMLNS:  "https://developers.google.com/blockly/xml",
	}
	bytes, _ := xml.MarshalIndent(xmlBlock, "", "  ")
	xmlContent := string(bytes)

	println(xmlContent)
	println()

	// reconversion of Blockly XML -> Falcon
	exprs := analysis.NewXMLParser(xmlContent).ParseBlockly()
	for _, expr := range exprs {
		println(expr.String())
	}
}
