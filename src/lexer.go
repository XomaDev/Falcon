package main

import (
	"Falcon/sugar"
	"Falcon/types"
	"strconv"
	"strings"
)

type Lexer struct {
	source    string
	sourceLen int
	currIndex int
	currLine  int
}

func NewLexer(source string) *Lexer {
	return &Lexer{
		source:    source,
		sourceLen: len(source),
		currIndex: 0,
		currLine:  1,
	}
}

var simpleCharTypes = map[string]types.Type{
	"+": types.Operator,
	// '-' parsed manually
	"*": types.Operator, "/": types.Operator, "^": types.Operator,
	"&": types.Operator, "|": types.Operator,
	"~": types.Operator,

	"(": types.OpenCurve, ")": types.CloseCurve,
	"[": types.OpenSquare, "]": types.CloseSquare,
	"{": types.OpenCurly, "}": types.CloseCurly,
	"=": types.Equals,
	".": types.Dot,
	",": types.Comma,
}

var keywordTypes = map[string]types.Type{
	"true":  types.Bool,
	"false": types.Bool,
}

func (l *Lexer) Lex() []types.Token {
	var tokens []types.Token
	for {
		l.trim()
		if l.isEOF() {
			break
		}
		tokens = append(tokens, l.parse())
	}
	return tokens
}

func (l *Lexer) trim() {
	for l.notEOF() {
		p := l.peek()
		switch p {
		case ' ':
			l.currIndex++
			continue
		case '\n':
			l.currIndex++
			l.currLine++
			continue
		}
		break
	}
}

func (l *Lexer) parse() types.Token {
	char := l.next()
	s := string(char)
	resType, ok := simpleCharTypes[s]

	if ok {
		return l.makeToken(resType, s)
	}

	switch char {
	case '"':
		return l.parseText()
	case '<':
		if l.consume('=') {
			return l.makeToken(types.LesserThanEquals, "<=")
		}
		return l.makeToken(types.LesserThan, "<")
	case '>':
		if l.consume('=') {
			return l.makeToken(types.GreaterThanEquals, ">=")
		}
		return l.makeToken(types.GreaterThan, ">")
	case '-':
		if l.consume('>') {
			return l.makeToken(types.RightArrow, "->")
		}
		return l.makeToken(types.Operator, "-")
	default:
		l.currIndex--
		if l.isDigit() {
			return l.parseNumeric()
		} else if l.isAlpha() {
			return l.parseAlpha()
		}
	}
	l.lexError("Unexpected character '%'", s)
	return types.Token{}
}

func (l *Lexer) parseText() types.Token {
	startIndex := l.currIndex
	for l.notEOF() && l.peek() != '"' {
		l.currIndex++
	}
	l.eat('"')
	return l.makeToken(types.Text, l.source[startIndex:l.currIndex-1])
}

func (l *Lexer) parseAlpha() types.Token {
	startIndex := l.currIndex
	l.currIndex++
	for l.notEOF() && l.isAlphaNumeric() {
		l.currIndex++
	}
	content := l.source[startIndex:l.currIndex]
	resType, ok := keywordTypes[content]
	if ok {
		// it's a keyword!
		return l.makeToken(resType, content)
	}
	return l.makeToken(types.Alpha, content)
}

func (l *Lexer) parseNumeric() types.Token {
	var numb strings.Builder
	l.writeNumeric(&numb)

	if l.notEOF() && l.peek() == '.' {
		l.next()
		numb.WriteByte('.')
		l.writeNumeric(&numb)
	}
	return l.makeToken(types.Number, numb.String())
}

func (l *Lexer) writeNumeric(builder *strings.Builder) {
	startIndex := l.currIndex
	for l.notEOF() && l.isDigit() {
		l.currIndex++
	}
	builder.WriteString(l.source[startIndex:l.currIndex])
}

func (l *Lexer) isAlphaNumeric() bool {
	return l.isAlpha() || l.isDigit()
}

func (l *Lexer) isAlpha() bool {
	c := l.peek()
	return c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c == '_'
}

func (l *Lexer) isDigit() bool {
	c := l.peek()
	return c >= '0' && c <= '9'
}

func (l *Lexer) eat(expect uint8) {
	got := l.source[l.currIndex]
	if got != expect {
		l.lexError("Expected '%', but got '%'", string(expect), string(got))
	}
	l.currIndex++
}

func (l *Lexer) lexError(message string, args ...string) {
	panic("[line " + strconv.Itoa(l.currLine) + "] " + sugar.Format(message, args...))
}

func (l *Lexer) consume(expect uint8) bool {
	if l.source[l.currIndex] != expect {
		return false
	}
	l.currIndex++
	return true
}

func (l *Lexer) peek() uint8 {
	return l.source[l.currIndex]
}

func (l *Lexer) next() uint8 {
	ch := l.source[l.currIndex]
	l.currIndex++
	return ch
}

func (l *Lexer) isEOF() bool {
	return l.currIndex >= l.sourceLen
}

func (l *Lexer) notEOF() bool {
	return l.currIndex < l.sourceLen
}

func (l *Lexer) makeToken(t types.Type, content string) types.Token {
	return types.Token{Type: t, Line: l.currLine, Content: &content}
}
