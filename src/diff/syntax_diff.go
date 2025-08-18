package diff

import (
	"Falcon/context"
	"Falcon/lex"
	"Falcon/sugar"
	"strings"
)

type SyntaxDiff struct {
	humanSyntax   string
	machineSyntax string

	humanWords   []string
	machineWords []string

	hSize int
	mSize int

	hIndex int
	mIndex int

	mergedWords strings.Builder
}

func MakeSyntaxDiff(humanSyntax string, machineSyntax string) *SyntaxDiff {
	humanWords := makeSimpleWords(humanSyntax, lex.NewLexer(&context.CodeContext{SourceCode: &humanSyntax, FileName: "human.diff"}).Lex())
	machineWords := makeSimpleWords(machineSyntax, lex.NewLexer(&context.CodeContext{SourceCode: &machineSyntax, FileName: "machine.diff"}).Lex())

	return &SyntaxDiff{
		humanSyntax:   humanSyntax,
		machineSyntax: machineSyntax,
		humanWords:    humanWords,
		machineWords:  machineWords,
		hSize:         len(humanWords),
		mSize:         len(machineWords),
		hIndex:        0,
		mIndex:        0,
		mergedWords:   strings.Builder{},
	}
}

func makeSimpleWords(sourceCode string, tokens []*lex.Token) []string {
	var words []string

	var lastTokenColumn = -1
	var lastTokenRow = 1

	for _, token := range tokens {
		var tokenContent string
		if token.Type == lex.Text {
			tokenContent = "\"" + *token.Content + "\""
		} else {
			tokenContent = *token.Content
		}
		words = append(words, tokenContent)
		if lastTokenColumn != -1 {
			if lastTokenColumn != token.Column {
				var separation strings.Builder

				separationLines := token.Column - lastTokenColumn
				separation.WriteString(strings.Repeat("\n", separationLines))

				afterNthLine := sugar.IndexAfterNthOccurrence(sourceCode, token.Column-1, '\n') + 1
				beforeNthSpace := afterNthLine + token.Row - len(tokenContent)
				separation.WriteString(sourceCode[afterNthLine:beforeNthSpace])
				words = append(words, separation.String())
			} else {
				afterNthLine := sugar.IndexAfterNthOccurrence(sourceCode, token.Column-1, '\n') + 1
				afterNthRow := afterNthLine + lastTokenRow
				beforeNthRow := afterNthLine + token.Row - len(tokenContent)

				separation := sourceCode[afterNthRow:beforeNthRow]
				words = append(words, separation)
			}
		} else {
			words = append(words, "")
		}

		lastTokenColumn = token.Column
		lastTokenRow = token.Row
	}
	return words
}

func (s *SyntaxDiff) Merge() string {
	for s.mIndex < s.mSize && s.hIndex < s.hSize {
		mWord := s.machineWords[s.mIndex]
		hWord := s.humanWords[s.hIndex]

		if mWord == hWord {
			s.hIndex++
			// append human spacing
			s.mergedWords.WriteString(s.humanWords[s.hIndex])
			// append matching word
			s.mergedWords.WriteString(hWord)
			s.hIndex++
			// skip machine word and spacing
			s.mIndex += 2
		} else {
			s.mIndex++
			// append machine spacing
			s.mergedWords.WriteString(s.machineWords[s.mIndex])
			// append machine word
			s.mergedWords.WriteString(mWord)
			s.mIndex++
			// skip mismatched word and spacing
			s.hIndex += 2
		}
	}
	// drain all remaining mWords
	for s.mIndex < s.mSize {
		mWord := s.machineWords[s.mIndex]
		s.mIndex++
		// append machine spacing
		s.mergedWords.WriteString(s.machineWords[s.mIndex])
		// append machine word
		s.mergedWords.WriteString(mWord)
		s.mIndex++
	}
	return s.mergedWords.String()
}
