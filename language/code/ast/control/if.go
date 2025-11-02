package control

import (
	ast2 "Falcon/code/ast"
	"strings"
)

type If struct {
	Conditions []ast2.Expr
	Bodies     [][]ast2.Expr
	ElseBody   []ast2.Expr
}

func (i *If) Yail() string {
	yail := ""
	for k, cond := range i.Conditions {
		if k != 0 {
			yail += " (begin "
		}
		yail += "(if "
		yail += cond.Yail()
		yail += "(begin "
		yail += ast2.PadBodyYail(i.Bodies[k])
		yail += ")"
	}
	if i.ElseBody != nil {
		yail += " (begin "
		yail += ast2.PadBodyYail(i.ElseBody)
		yail += ")"
	}
	for k := 0; k < len(i.Conditions)-1; k++ {
		yail += "))"
	}
	yail += ")"
	return yail
}

func (i *If) String() string {
	var builder strings.Builder

	numConditions := len(i.Conditions)
	currI := 0

	builder.WriteString("if ")
	for {
		builder.WriteString(i.Conditions[currI].String())
		builder.WriteString(" {\n")
		builder.WriteString(ast2.PadBody(i.Bodies[currI]))
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
		builder.WriteString(ast2.PadBody(i.ElseBody))
		builder.WriteString("}")
	}
	return builder.String()
}

func (i *If) Blockly() ast2.Block {
	conditions := ast2.ValuesByPrefix("IF", i.Conditions)
	bodies := ast2.ToStatements("DO", i.Bodies)

	numbElifs := len(i.Conditions) - 1
	var numbElse int

	if i.ElseBody != nil {
		bodies = append(bodies, ast2.CreateStatement("ELSE", i.ElseBody))
		numbElse = 1
	} else {
		numbElifs = 0
	}
	return ast2.Block{
		Type:       "controls_if",
		Mutation:   &ast2.Mutation{ElseIfCount: numbElifs, ElseCount: numbElse},
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
