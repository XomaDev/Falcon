//go:build !js && !wasm

package main

import (
	codeAnalysis "Falcon/code/analysis"
	"Falcon/code/ast/blockly"
	"Falcon/code/context"
	"Falcon/code/diff"
	"Falcon/code/lex"

	designAnalysis "Falcon/design"
	"encoding/xml"
	"os"
	"strings"
)

func main() {
	println("Hello from Falcon!\n")

	//diffTest()
	//analyzeSyntax()
	//xmlTest()
	designTest()
}

func designTest() {
	xmlFile := "Screen1.aiml"
	xmlPath := "/home/melon/GolandProjects/Falcon/" + xmlFile
	codeBytes, err := os.ReadFile(xmlPath)
	if err != nil {
		panic(err)
	}
	xmlString := string(codeBytes)
	result, err := designAnalysis.NewXmlParser(xmlString).ConvertXmlToSchema()
	if err != nil {
		panic(err)
	}
	println(result)
}

func xmlTest() {
	xmlFile := "xml.txt"
	xmlPath := "/home/ekina/GolandProjects/Falcon/" + xmlFile
	codeBytes, err := os.ReadFile(xmlPath)
	if err != nil {
		panic(err)
	}
	xmlString := string(codeBytes)
	exprs := codeAnalysis.NewXMLParser(xmlString).ParseBlockly()
	var machineSourceCode strings.Builder
	for _, expr := range exprs {
		machineSourceCode.WriteString(expr.String())
		machineSourceCode.WriteRune('\n')
	}
	println(machineSourceCode.String())
}

func diffTest() {
	diff0 := "diff0.mist"
	diff1 := "diff1.mist"

	diff0Path := "/home/ekina/GolandProjects/Falcon/" + diff0
	codeBytes, err := os.ReadFile(diff0Path)
	if err != nil {
		panic(err)
	}
	diff0Code := string(codeBytes)

	diff1Path := "/home/ekina/GolandProjects/Falcon/" + diff1
	codeBytes, err = os.ReadFile(diff1Path)
	if err != nil {
		panic(err)
	}
	diff1Code := string(codeBytes)

	syntaxDiff := diff.MakeSyntaxDiff(diff0Code, diff1Code)
	println(syntaxDiff.Merge())
}

func analyzeSyntax() {
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
	langParser := codeAnalysis.NewLangParser(tokens)
	expressions := langParser.ParseAll()
	println(langParser.GetComponentDefinitionsCode())
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
	exprs := codeAnalysis.NewXMLParser(xmlContent).ParseBlockly()
	var machineSourceCode strings.Builder
	for _, expr := range exprs {
		machineSourceCode.WriteString(expr.String())
		machineSourceCode.WriteRune('\n')
	}
	println(machineSourceCode.String())

	// Generate a merged syntax
	println("\n=== DIFF ===\n")
	syntaxDiff := diff.MakeSyntaxDiff(sourceCode, machineSourceCode.String())
	println(syntaxDiff.Merge())
}
