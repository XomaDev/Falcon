package ast

type NumExpr struct {
	Expr
	Content *string
}

func (n *NumExpr) String() string {
	return *n.Content
}

func (n *NumExpr) Blockly() Block {
	return Block{
		Type:   "math_number",
		Fields: ToFields(map[string]string{"NUM": *n.Content}),
	}
}
