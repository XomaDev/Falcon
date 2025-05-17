package method

import (
	"Falcon/ast/blockly"
	"Falcon/lex"
	"Falcon/sugar"
	"strconv"
)

type Call struct {
	Where lex.Token
	On    blockly.Expr
	Name  string
	Args  []blockly.Expr
}

type Signature struct {
	Module     string
	Name       string
	ParamCount int
}

func makeSignature(module string, name string, paramCount int) *Signature {
	return &Signature{Module: module, Name: name, ParamCount: paramCount}
}

var Signatures = map[string]*Signature{
	"startsWith":              makeSignature("text", "text_starts_at", 1),
	"contains":                makeSignature("text", "text_contains", 1),
	"containsAny":             makeSignature("text", "text_contains", 1),
	"containsAll":             makeSignature("text", "text_contains", 1),
	"split":                   makeSignature("text", "text_split", 1),
	"splitAtFirst":            makeSignature("text", "text_split", 1),
	"splitAtAny":              makeSignature("text", "text_split", 1),
	"splitAtFirstOfAny":       makeSignature("text", "text_split", 1),
	"segment":                 makeSignature("text", "text_segment", 2),
	"replace":                 makeSignature("text", "text_replace_all", 2),
	"replaceFrom":             makeSignature("text", "text_replace_mappings", 1),
	"replaceFromLongestFirst": makeSignature("text", "text_replace_mappings", 1),
}

func (m *Call) String() string {
	return sugar.Format("%.%(%)", m.On.String(), m.Name, blockly.JoinExprs(", ", m.Args))
}

func (m *Call) Blockly() blockly.Block {
	signature, ok := Signatures[m.Name]
	if !ok {
		m.Where.Error("Cannot find method .%", m.Name)
	}
	gotArgLen := len(m.Args)
	if signature.ParamCount != gotArgLen {
		m.Where.Error("Expected % args but got % for method .%",
			strconv.Itoa(signature.ParamCount), strconv.Itoa(gotArgLen), m.Name)
	}
	switch signature.Module {
	case "text":
		return m.textMethods(signature)
	default:
		panic("Unknown module " + signature.Module)
	}
}
