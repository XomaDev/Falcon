package control

import (
	blockly2 "Falcon/code/ast/blockly"
	"strings"
)

type If struct {
	Conditions []blockly2.Expr
	Bodies     [][]blockly2.Expr
	ElseBody   []blockly2.Expr
}

func (i *If) String() string {
	var builder strings.Builder

	numConditions := len(i.Conditions)
	currI := 0

	builder.WriteString("if ")
	for {
		builder.WriteString(i.Conditions[currI].String())
		builder.WriteString(" {\n")
		builder.WriteString(blockly2.PadBody(i.Bodies[currI]))
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
		builder.WriteString(blockly2.PadBody(i.ElseBody))
		builder.WriteString("}")
	}
	return builder.String()
}

func (i *If) Blockly() blockly2.Block {
	conditions := blockly2.ValuesByPrefix("IF", i.Conditions)
	bodies := blockly2.ToStatements("DO", i.Bodies)

	numbElifs := len(i.Conditions) - 1
	var numbElse int

	if i.ElseBody != nil {
		bodies = append(bodies, blockly2.CreateStatement("ELSE", i.ElseBody))
		numbElse = 1
	} else {
		numbElifs = 0
	}
	return blockly2.Block{
		Type:       "controls_if",
		Mutation:   &blockly2.Mutation{ElseIfCount: numbElifs, ElseCount: numbElse},
		Values:     conditions,
		Statements: bodies,
	}
}

func (i *If) Continuous() bool {
	return false
}

func (i *If) Consumable() bool {
	return false
}
