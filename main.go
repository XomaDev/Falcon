package main

func main() {
	println("Hello from Falcon!")
	sourceCode := `123 + 456`
	tokens := NewLexer(sourceCode).Lex()
	for _, token := range tokens {
		println(token.String())
	}
}
