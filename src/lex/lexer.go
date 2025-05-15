package lex

import (
	"Falcon/sugar"
	"strconv"
	"strings"
)

type Lexer struct {
	source    string
	sourceLen int
	currIndex int
	currLine  int
	Tokens    []Token
}

func NewLexer(source string) *Lexer {
	return &Lexer{
		source:    source,
		sourceLen: len(source),
		currIndex: 0,
		currLine:  1,
		Tokens:    []Token{},
	}
}

func (l *Lexer) Lex() []Token {
	for l.notEOF() {
		l.parse()
	}
	return l.Tokens
}

func (l *Lexer) parse() {
	c := l.next()

	switch c {
	case '\n':
		l.currLine++
		return
	case ' ', '\t':
		return
	}
	switch c {
	case '+':
		l.createOp("+")
	case '-':
		if l.consume('>') {
			l.createOp("->")
		} else {
			l.createOp("-")
		}
	case '*':
		l.createOp("*")
	case '/':
		l.createOp("/")
	case '|':
		if l.consume('|') {
			l.createOp("||")
		} else {
			l.createOp("|")
		}
	case '&':
		if l.consume('&') {
			l.createOp("&&")
		} else {
			l.createOp("&")
		}
	case '~':
		l.createOp("~")
	case '<':
		if l.consume('=') {
			l.createOp("<=")
		} else {
			l.createOp("<")
		}
	case '>':
		if l.consume('=') {
			l.createOp(">=")
		} else {
			l.createOp(">")
		}
	case '(':
		l.createOp("(")
	case ')':
		l.createOp(")")
	case '[':
		l.createOp("[")
	case ']':
		l.createOp("]")
	case '{':
		l.createOp("{")
	case '}':
		l.createOp("}")
	case '=':
		if l.consume('=') {
			l.createOp("==")
		} else {
			l.createOp("=")
		}
	case '.':
		l.createOp(".")
	case ',':
		l.createOp(",")
	case '?':
		l.createOp("?")
	case '!':
		if l.consume('=') {
			l.createOp("!=")
		} else {
			l.createOp("!")
		}
	case ':':
		if l.consume(':') {
			l.createOp("::")
		} else {
			l.createOp(":")
		}
	case '"':
		l.text()
	default:
		l.currIndex--
		if l.isAlpha() {
			l.alpha()
		} else if l.isDigit() {
			l.numeric()
		} else {
			l.error("Unexpected character '%'", string(c))
		}
	}
}

func (l *Lexer) createOp(op string) {
	sToken, ok := Symbols[op]
	if !ok {
		l.error("Bad createOp('%')", op)
	} else {
		l.appendToken(sToken.normal(l.currLine, op))
	}
}

func (l *Lexer) text() {
	startIndex := l.currIndex
	for l.notEOF() && l.peek() != '"' {
		l.currIndex++
	}
	l.eat('"')
	content := l.source[startIndex : l.currIndex-1]
	l.appendToken(Token{Type: Text, Content: &content, Flags: []Flag{Value, ConstantValue}})
}

func (l *Lexer) alpha() {
	startIndex := l.currIndex
	l.currIndex++
	for l.notEOF() && l.isAlphaNumeric() {
		l.currIndex++
	}
	content := l.source[startIndex:l.currIndex]
	sToken, ok := Keywords[content]
	if ok {
		l.appendToken(sToken.normal(l.currLine))
	} else {
		l.appendToken(Token{Type: Name, Content: &content, Flags: []Flag{Value}})
	}
}

func (l *Lexer) numeric() {
	var numb strings.Builder
	l.writeNumeric(&numb)
	if l.notEOF() && l.peek() == '.' {
		l.currIndex++
		numb.WriteByte('.')
		l.writeNumeric(&numb)
	}
	content := numb.String()
	l.appendToken(Token{Type: Number, Content: &content, Flags: []Flag{Value, ConstantValue}})
}

func (l *Lexer) appendToken(token Token) {
	l.Tokens = append(l.Tokens, token)
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
	got := l.next()
	if got != expect {
		l.error("Expected '%', but got '%'", string(expect), string(got))
	}
}

func (l *Lexer) error(message string, args ...string) {
	panic("[line " + strconv.Itoa(l.currLine) + "] " + sugar.Format(message, args...))
}

func (l *Lexer) consume(expect uint8) bool {
	if l.peek() == expect {
		l.currIndex++
		return true
	}
	return false
}

func (l *Lexer) peek() uint8 {
	return l.source[l.currIndex]
}

func (l *Lexer) next() uint8 {
	c := l.source[l.currIndex]
	l.currIndex++
	return c
}

func (l *Lexer) isEOF() bool {
	return l.currIndex >= l.sourceLen
}

func (l *Lexer) notEOF() bool {
	return l.currIndex < l.sourceLen
}
