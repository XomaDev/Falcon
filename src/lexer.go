package main

import (
	"Falcon/label"
	"Falcon/sugar"
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

var symbols = map[string][]label.Quality{
	"+": {label.Operator},
	// '-' parsed manually
	"*": {label.Operator}, "/": {label.Operator}, "^": {label.Operator},
	"&": {label.Operator}, "|": {label.Operator},
	"~": {label.Operator},

	"(": {label.OpenCurve}, ")": {label.CloseCurve},
	"[": {label.OpenSquare}, "]": {label.CloseSquare},
	"{": {label.OpenCurly}, "}": {label.CloseCurly},
	".": {label.Dot},
	",": {label.Comma},
	"?": {label.Question},
	"!": {label.Not},

	"==": {label.Equality},
	"!=": {label.Equality},
}

var keywords = map[string]label.Quality{
	"true":  label.Bool,
	"false": label.Bool,
	"if":    label.If,
	"elif":  label.Elif,
	"else":  label.Else,
}

func (l *Lexer) Lex() []label.Token {
	var tokens []label.Token
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

func (l *Lexer) parse() label.Token {
	char := l.next()
	s := string(char)
	qualities, ok := symbols[s]

	if ok {
		return l.makeToken(qualities, s)
	}

	switch char {
	case '"':
		return l.parseText()
	case '<':
		if l.consume('=') {
			return l.simpleToken(label.LesserThanEquals, "<=")
		}
		return l.simpleToken(label.LesserThan, "<")
	case '>':
		if l.consume('=') {
			return l.simpleToken(label.GreaterThanEquals, ">=")
		}
		return l.simpleToken(label.GreaterThan, ">")
	case '-':
		if l.consume('>') {
			return l.simpleToken(label.RightArrow, "->")
		}
		return l.simpleToken(label.Operator, "-")
	default:
		l.currIndex--
		if l.isDigit() {
			return l.parseNumeric()
		} else if l.isAlpha() {
			return l.parseAlpha()
		}
	}
	l.lexError("Unexpected character '%'", s)
	return label.Token{}
}

func (l *Lexer) parseText() label.Token {
	startIndex := l.currIndex
	for l.notEOF() && l.peek() != '"' {
		l.currIndex++
	}
	l.eat('"')
	return l.simpleToken(label.Text, l.source[startIndex:l.currIndex-1])
}

func (l *Lexer) parseAlpha() label.Token {
	startIndex := l.currIndex
	l.currIndex++
	for l.notEOF() && l.isAlphaNumeric() {
		l.currIndex++
	}
	content := l.source[startIndex:l.currIndex]
	resType, ok := keywords[content]
	if ok {
		// it's a keyword!
		return l.simpleToken(resType, content)
	}
	return l.simpleToken(label.Alpha, content)
}

func (l *Lexer) parseNumeric() label.Token {
	var numb strings.Builder
	l.writeNumeric(&numb)

	if l.notEOF() && l.peek() == '.' {
		l.next()
		numb.WriteByte('.')
		l.writeNumeric(&numb)
	}
	return l.simpleToken(label.Number, numb.String())
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

func (l *Lexer) simpleToken(quality label.Quality, content string) label.Token {
	return label.Token{Quality: quality, Content: &content}
}

func (l *Lexer) makeToken(q []label.Quality, content string) label.Token {
	return label.Token{Quality: q[0], AllQualities: q[1:], Line: l.currLine, Content: &content}
}
