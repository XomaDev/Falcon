//go:build !js && !wasm

package main

import (
	codeAnalysis "Falcon/code/analysis"
	"Falcon/code/ast"
	"Falcon/code/context"
	"Falcon/code/diff"
	"Falcon/code/lex"
	designAnalysis "Falcon/design"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	println("Hello from Falcon!\n")

	//diffTest()
	analyzeSyntax()
	//xmlTest()
	//designTest()
	//testSyntheticData()
}

func testSyntheticData() {
	example_codes := "/home/ekina/GolandProjects/Falcon/examples"

	// Determine the output file path in the parent directory of example_codes (the /Falcon dir)
	falconDir := filepath.Dir("/home/ekina/GolandProjects/Falcon/")
	outputFile := filepath.Join(falconDir, "falcon_dsl.txt")

	// Use a strings.Builder for efficient string concatenation
	var combinedContent strings.Builder

	// 1. Read all the files in the directory
	files, err := ioutil.ReadDir(example_codes)
	if err != nil {
		// Note: This path is user-specific and likely won't exist in a general environment.
		// Error handling for this path is important.
		fmt.Printf("Error reading directory %s: %v\n", example_codes, err)
		fmt.Println("Please ensure the directory exists and is accessible.")
		return
	}

	fileCount := 0
	// 2. Iterate through the files and combine content
	for _, file := range files {
		// Skip directories and non-regular files
		if file.IsDir() || !file.Mode().IsRegular() {
			continue
		}

		filename := file.Name()

		fmt.Println("Concatenating file " + filename)
		filepath := filepath.Join(example_codes, filename)

		// 3. Read the content of the file
		content, err := os.ReadFile(filepath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filename, err)
			continue
		}

		// Append file content
		combinedContent.Write(content)

		// Add the separator "----" only between files
		// We add a newline for readability if the last content didn't end with one.
		combinedContent.WriteString("\n----")

		fileCount++
	}

	// Check if any content was actually loaded
	if combinedContent.Len() == 0 {
		fmt.Println("No files with content were processed.")
		return
	}

	// Get the final concatenated string. We remove the last appended "----"
	// to avoid an extra empty element in the final split array.
	finalContent := combinedContent.String()
	// Safely remove the trailing separator
	if strings.HasSuffix(finalContent, "\n----") {
		finalContent = finalContent[:len(finalContent)-len("\n----")]
	}

	fmt.Printf("\nSuccessfully combined content from %d files. Starting unified data split...\n", fileCount)

	falconDocs := "/home/ekina/GolandProjects/Falcon/README.md"
	codeBytes, err := os.ReadFile(falconDocs)
	if err != nil {
		panic(err)
	}
	combinedResult := string(codeBytes) + "\n----" + "\n" + finalContent

	// NEW STEP 6: Write the finalContent to the file in the parent directory
	// 0644 is standard read/write permissions for the owner.
	err = os.WriteFile(outputFile, []byte(combinedResult), 0644)
	if err != nil {
		fmt.Printf("Error writing combined content to file %s: %v\n", outputFile, err)
	} else {
		fmt.Printf("Successfully wrote combined content to: %s\n", outputFile)
	}

	// 4. Split the entire combined content using "----" as delimiter
	parts := strings.Split(finalContent, "----")

	// 5. Trim all the strings and verify
	for i, part := range parts {
		trimmed := strings.TrimSpace(part)
		// Since the file origin is lost, we pass a generic identifier for verification
		verifySyntFunc(fmt.Sprintf("combined_part_%d", i+1), trimmed)
	}

	fmt.Println("\nSynthetic data loading and verification complete.")
}

func verifySyntFunc(fileName string, sourceCode string) {
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

	blocks := make([]ast.Block, len(expressions))
	for i, expression := range expressions {
		blocks[i] = expression.Blockly(true)
	}
	xmlBlock := ast.XmlRoot{
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
}

func designTest() {
	xmlFile := "Screen1.aiml"
	xmlPath := "/home/ekina/GolandProjects/Falcon/testing/" + xmlFile
	codeBytes, err := os.ReadFile(xmlPath)
	if err != nil {
		panic(err)
	}
	xmlString := string(codeBytes)
	schemaString, err := designAnalysis.NewXmlParser(xmlString).ConvertXmlToSchema()
	if err != nil {
		panic(err)
	}
	println(schemaString)
	xmlString, err = designAnalysis.NewSchemaParser(schemaString).ConvertSchemaToXml()
	if err != nil {
		panic(err)
	}
	println("Produced XML: ")
	println(xmlString)
}

func xmlTest() {
	xmlFile := "xml.txt"
	xmlPath := "/home/ekina/GolandProjects/Falcon/testing/" + xmlFile
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

	diff0Path := "/home/ekina/GolandProjects/Falcon/testing/" + diff0
	codeBytes, err := os.ReadFile(diff0Path)
	if err != nil {
		panic(err)
	}
	diff0Code := string(codeBytes)

	diff1Path := "/home/ekina/GolandProjects/Falcon/testing/" + diff1
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
	filePath := "/home/ekina/GolandProjects/Falcon/testing/" + fileName
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

	blocks := make([]ast.Block, len(expressions))
	for i, expression := range expressions {
		blocks[i] = expression.Blockly(true)
	}
	xmlBlock := ast.XmlRoot{
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

	//// Generate a merged syntax
	//println("\n=== DIFF ===\n")
	//syntaxDiff := diff.MakeSyntaxDiff(sourceCode, machineSourceCode.String())
	//println(syntaxDiff.Merge())
}
