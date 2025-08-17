//go:build js && wasm
// +build js,wasm

// GOOS=js GOARCH=wasm go build -o web/falcon.wasm

package main

import (
	"Falcon/analysis"
	"Falcon/ast/blockly"
	"Falcon/context"
	"Falcon/lex"
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
		if len(p) < 1 {
			return js.ValueOf("No Mist content provided")
		}
		sourceCode := p[0].String()
		codeContext := &context.CodeContext{SourceCode: &sourceCode, FileName: "appinventor.live"}

		tokens := lex.NewLexer(codeContext).Lex()
		expressions := analysis.NewLangParser(tokens).ParseAll()

		var xmlCode strings.Builder

		for _, expression := range expressions {
			xmlBlock := blockly.XmlRoot{
				Blocks: []blockly.Block{expression.Blockly()},
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
		exprs := analysis.NewXMLParser(xmlContent).ParseBlockly()
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

func main() {
	println("Hello from falcon.go!")

	c := make(chan struct{}, 0)
	js.Global().Set("mistToXml", js.FuncOf(mistToXml))
	js.Global().Set("xmlToMist", js.FuncOf(xmlToMist))
	<-c
}
