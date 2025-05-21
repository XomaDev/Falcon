//go:build js && wasm
// +build js,wasm

package main

import (
	"Falcon/analysis"
	"Falcon/ast/blockly"
	"Falcon/lex"
	"encoding/xml"
	"strings"
	"syscall/js"
)

func safeExec(fn func() js.Value) js.Value {
	defer func() {
		if r := recover(); r != nil {
			println("Recovering from panic:", r)
		}
	}()
	return fn()
}

func mistToXml(this js.Value, p []js.Value) any {
	return safeExec(func() js.Value {
		if len(p) < 1 {
			return js.ValueOf("No Mist content provided")
		}
		sourceCode := p[0].String()

		tokens := lex.NewLexer(sourceCode).Lex()
		expressions := analysis.NewParser(tokens).ParseAll()

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

func main() {
	println("Hello from falcon.go!")

	c := make(chan struct{}, 0)
	js.Global().Set("mistToXml", js.FuncOf(mistToXml))
	<-c
}
