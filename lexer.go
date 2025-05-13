package main

import (
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

func (l *Lexer) Lex() []Token {
	var tokens []Token
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
	for l.notEOF() && l.peek() == ' ' {
		l.currIndex++
	}
}

func (l *Lexer) parse() Token {
	char := l.next()
	s := string(char)
	switch char {
	case '\n':
		l.currLine++
		break
	case '+', '-', '*', '/':
		return l.makeToken(Operator, s)
	case '(':
		return l.simpleToken(OpenCurve)
	case ')':
		return l.simpleToken(CloseCurve)
	case '[':
		return l.simpleToken(OpenSquare)
	case ']':
		return l.simpleToken(CloseSquare)
	case '{':
		return l.simpleToken(OpenCurly)
	case '}':
		return l.simpleToken(CloseCurly)
	default:
		l.currIndex--
		if l.isDigit() {
			return l.parseNumeric()
		}
	}
	l.lexError("Unexpected character '%'", s)
	return Token{}
}

func (l *Lexer) parseNumeric() Token {
	var numb strings.Builder
	l.writeNumeric(&numb)

	if l.notEOF() && l.peek() == '.' {
		numb.WriteByte('.')
		l.writeNumeric(&numb)
	}
	return l.makeToken(Number, numb.String())
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
}

func (l *Lexer) lexError(message string, args ...string) {
	panic("[line " + strconv.Itoa(l.currLine) + "] " + Format(message, args...))
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

func (l *Lexer) simpleToken(t Type) Token {
	return Token{Type: t, Line: l.currLine}
}

func (l *Lexer) makeToken(t Type, content string) Token {
	return Token{Type: t, Line: l.currLine, Content: &content}
}
