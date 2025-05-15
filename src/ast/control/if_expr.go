package control

import (
	"Falcon/ast/blockly"
	"strings"
)

type If struct {
	Conditions []blockly.Expr
	Bodies     [][]blockly.Expr
	ElseBody   []blockly.Expr
}

func (i *If) String() string {
	var builder strings.Builder

	numConditions := len(i.Conditions)
	currI := 0

	builder.WriteString("if ")
	for {
		builder.WriteString(i.Conditions[currI].String())
		builder.WriteString(" {\n")
		builder.WriteString(blockly.PadBody(i.Bodies[currI]))
		builder.WriteString("}")
		currI++
		if currI < numConditions {
			builder.WriteString(" elif ")
		} else {
			break
		}
	}
	if i.ElseBody != nil {
		builder.WriteString(" else {\n")
		builder.WriteString(blockly.PadBody(i.ElseBody))
		builder.WriteString("}")
	}
	return builder.String()
}

func (i *If) Blockly() blockly.Block {
	conditions := blockly.ToValues("IF", i.Conditions)
	bodies := blockly.ToStatements("DO", i.Bodies)

	numbElifs := len(i.Conditions) - 1
	var numbElse int

	if i.ElseBody != nil {
		bodies = append(bodies, blockly.CreateStatement("ELSE", i.ElseBody))
		numbElse = 1
	} else {
		numbElifs = 0
	}
	return blockly.Block{
		Type:       "controls_if",
		Mutation:   &blockly.Mutation{ElseIfCount: numbElifs, ElseCount: numbElse},
		Values:     conditions,
		Statements: bodies,
	}
}
