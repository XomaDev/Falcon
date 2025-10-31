package method

func (c *Call) textMethodsYail(signature *Signature) string {
	switch signature.Name {
	case "text_length":
		return c.simpleYailOperand("string-reverse", "reverse")
	case "text_trim":
		return c.simpleYailOperand("string-trim", "trim")
	case "text_split_at_spaces":
		return c.simpleYailOperand("string-split-at-spaces", "split at spaces")
	case "text_reverse":
		return c.simpleYailOperand("string-reverse", "reverse")
	case "text_changeCase":
		return c.yailTextChangeCase()
	default:
		panic("Unknown text method " + signature.Name)
	}
}

func (c *Call) yailTextChangeCase() string {
	switch c.Name {
	case "uppercase":
		return c.simpleYailOperand("string-to-upper-case", "upcase")
	case "lowercase":
		return c.simpleYailOperand("string-to-lowercase", "lowercase")
	}
	panic("Unknown text change case method: " + c.Name + "()")
}

func (c *Call) simpleYailOperand(yailCall string, name string) string {
	yail := "(call-yail-primitive "
	yail += yailCall
	yail += " (*list-for-runtime* "
	yail += c.Args[0].Yail()
	yail += ") '(text) \""
	yail += name
	yail += "\")"
	return yail
}
