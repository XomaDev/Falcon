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
	filePath := "/home/kumaraswamy/GolandProjects/Falcon/" + fileName
	codeBytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	sourceCode := string(codeBytes)
	codeContext := context.CodeContext{SourceCode: &sourceCode, FileName: fileName}

	tokens := lex.NewLexer(codeContext).Lex()
	for _, token := range tokens {
		println(token.String())
	}

	println("\n=== AST ===\n")

	expressions := analysis.NewParser(tokens).ParseAll()
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
	println(string(bytes))
}
