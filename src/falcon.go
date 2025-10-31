//go:build js && wasm
// +build js,wasm

// GOOS=js GOARCH=wasm go build -o web/falcon.wasm

package main

import (
	analysis2 "Falcon/code/analysis"
	"Falcon/code/ast"
	"Falcon/code/context"
	"Falcon/code/diff"
	"Falcon/code/lex"
	designAnalysis "Falcon/design"
	"encoding/xml"
	"runtime/debug"
	"strings"
	"syscall/js"
)

func safeExec(fn func() js.Value) js.Value {
	defer func() {
		if r := recover(); r != nil {
			switch err := r.(type) {
			case error:
				println("Recovered from panic:", err.Error())
			default:
				println("Recovered from panic (non-error):", r)
			}
			println("Stack trace:")
			println(string(debug.Stack())) // Print full stack trace
		}
	}()
	return fn()
}

// Code -> Blocks
func mistToXml(this js.Value, p []js.Value) any {
	return safeExec(func() js.Value {
		if len(p) < 2 {
			return js.ValueOf("mistToXML(sourceCode string, componentDefinitions map[string][]string) not provided!")
		}
		sourceCode := p[0].String()

		// Parse the Component Definition Context
		componentContextMap := make(map[string][]string) // Button -> [Button1, Button2]
		reverseComponentMap := make(map[string]string)   // Button1 -> Button, Button2 -> Button
		obj := p[1]
		keys := js.Global().Get("Object").Call("keys", obj)
		length := keys.Length()
		for i := 0; i < length; i++ {
			compType := keys.Index(i).String()
			jsArr := obj.Get(compType)
			var compNames []string
			for j := 0; j < jsArr.Length(); j++ {
				instanceName := jsArr.Index(j).String()
				compNames = append(compNames, instanceName)
				reverseComponentMap[instanceName] = compType
			}
			componentContextMap[compType] = compNames
		}

		// Parse Mist To XML Blockly
		codeContext := &context.CodeContext{SourceCode: &sourceCode, FileName: "appinventor.live"}

		tokens := lex.NewLexer(codeContext).Lex()
		langParser := analysis2.NewLangParser(tokens)
		langParser.SetComponentDefinitions(componentContextMap, reverseComponentMap)
		expressions := langParser.ParseAll()

		var xmlCode strings.Builder

		for _, expression := range expressions {
			xmlBlock := ast.XmlRoot{
				Blocks: []ast.Block{expression.Blockly()},
				XMLNS:  "https://developers.google.com/blockly/xml",
			}
			bytes, _ := xml.MarshalIndent(xmlBlock, "", "  ")

			xmlCode.WriteString(string(bytes))
			xmlCode.WriteByte(0)
		}

		return js.ValueOf(xmlCode.String())
	})
}

// Blocks -> Code
func xmlToMist(this js.Value, p []js.Value) any {
	return safeExec(func() js.Value {
		if len(p) < 1 {
			return js.ValueOf("No XML content provided")
		}
		xmlContent := p[0].String()
		exprs := analysis2.NewXMLParser(xmlContent).ParseBlockly()
		var builder strings.Builder

		for _, expr := range exprs {
			builder.WriteString(expr.String())
			builder.WriteString("\n")

			block := expr.Blockly()
			if block.Order() > 0 {
				builder.WriteString("\n")
			}
		}
		return js.ValueOf(builder.String())
	})
}

func mergeSyntaxDiff(this js.Value, p []js.Value) any {
	return safeExec(func() js.Value {
		if len(p) < 2 {
			return js.ValueOf("Requires two string arguments, [HumanSyntax, MachineSyntax]")
		}
		humanSyntax := p[0].String()
		machineSyntax := p[1].String()
		mergedSyntax := diff.MakeSyntaxDiff(humanSyntax, machineSyntax).Merge()
		return js.ValueOf(mergedSyntax)
	})
}

func convertSchemaToXml(this js.Value, p []js.Value) any {
	return safeExec(func() js.Value {
		if len(p) < 1 {
			return js.ValueOf("No schema provided")
		}
		schemaString, err := designAnalysis.NewSchemaParser(p[0].String()).ConvertSchemaToXml()
		if err != nil {
			panic(err)
		}
		return js.ValueOf(schemaString)
	})
}

func convertXmlToSchema(this js.Value, p []js.Value) any {
	return safeExec(func() js.Value {
		if len(p) < 1 {
			return js.ValueOf("No schema provided")
		}
		schemaString, err := designAnalysis.NewXmlParser(p[0].String()).ConvertXmlToSchema()
		if err != nil {
			panic(err)
		}
		return js.ValueOf(schemaString)
	})
}

func main() {
	println("Hello from falcon.go!")

	c := make(chan struct{}, 0)
	js.Global().Set("mistToXml", js.FuncOf(mistToXml))
	js.Global().Set("xmlToMist", js.FuncOf(xmlToMist))
	js.Global().Set("mergeSyntaxDiff", js.FuncOf(mergeSyntaxDiff))
	js.Global().Set("schemaToXml", js.FuncOf(convertSchemaToXml))
	js.Global().Set("xmlToSchema", js.FuncOf(convertXmlToSchema))
	<-c
}
