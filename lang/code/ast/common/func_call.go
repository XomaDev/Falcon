package common

import (
	"Falcon/code/ast"
	"Falcon/code/ast/fundamentals"
	"Falcon/code/ast/variables"
	"Falcon/code/lex"
	"Falcon/code/sugar"
	"strconv"
)

type FuncCall struct {
	Where *lex.Token
	Name  string
	Args  []ast.Expr
}

func CreateFuncCall(where *lex.Token, name string, args []ast.Expr) *FuncCall {
	call := &FuncCall{Where: where, Name: name, Args: args}
	call.Signature() // Ensures a valid function name
	return call
}

func (f *FuncCall) String() string {
	if f.Name == "rem" {
		return f.Args[0].String() + " % " + f.Args[1].String()
	}
	if f.Name == "neg" {
		if f.Args[0].Continuous() {
			return "-" + f.Args[0].String()
		}
		return "-(" + f.Args[0].String() + ")"
	}
	return sugar.Format("%(%)", f.Name, ast.JoinExprs(", ", f.Args))
}

func (f *FuncCall) Blockly(flags ...bool) ast.Block {
	if len(flags) > 0 && !flags[0] && !f.Consumable() {
		f.Where.Error("Expected a consumable but got a statement")
	}
	// TODO:
	//  We have to assert correct number of arguments
	//  both here and in Signature()
	switch f.Name {
	case "sqrt", "abs", "neg", "log", "exp", "round", "ceil", "floor",
		"sin", "cos", "tan", "asin", "acos", "atan", "degrees", "radians",
		"decToHex", "decToBin", "hexToDec", "binToDec":
		return f.mathConversions()

	case "dec", "bin", "octal", "hexa":
		return f.mathRadix()
	case "randInt":
		return f.randInt()
	case "randFloat":
		return f.randFloat()
	case "setRandSeed":
		return f.setRandSeed()
	case "min", "max":
		return f.minOrMax()
	case "avgOf", "maxOf", "minOf", "geoMeanOf", "stdDevOf", "stdErrOf":
		return f.mathOnList()
	case "modeOf":
		return f.modeOf()
	case "mod", "rem", "quot":
		return f.mathDivide()
	case "aTan2":
		return f.atan2()
	case "formatDecimal":
		return f.formatDecimal()

	case "println":
		return f.println()
	case "openScreen":
		return f.openScreen()
	case "openScreenWithValue":
		return f.openScreenWithValue()
	case "closeScreenWithValue":
		return f.closeScreenWithValue()
	case "getStartValue":
		return f.ctrlSimpleBlock("controls_getStartValue")
	case "closeScreen":
		return f.ctrlSimpleBlock("controls_closeScreen")
	case "closeApp":
		return f.ctrlSimpleBlock("controls_closeApplication")
	case "getPlainStartText":
		return f.ctrlSimpleBlock("controls_getPlainStartText")
	case "closeScreenWithPlainText":
		return f.closeScreenWithPlainText()
	case "copyList":
		return f.copyList()
	case "copyDict":
		return f.copyDict()

	case "makeColor":
		return f.makeColor()
	case "splitColor":
		return f.splitColor()

	case "set":
		return f.genericSet()
	case "get":
		return f.genericGet()
	case "call":
		return f.genericCall()
	case "every":
		return f.everyComponent()
	default:
		f.Where.Error("Unknown function %()", f.Name)
		panic("never reached")
	}
}

func (f *FuncCall) Continuous() bool {
	return true
}

func (f *FuncCall) Consumable(flags ...bool) bool {
	if f.Name == "setRandSeed" || f.Name == "println" ||
		f.Name == "openScreen" || f.Name == "openScreenWithValue" ||
		f.Name == "closeScreen" || f.Name == "closeScreenWithValue" ||
		f.Name == "closeApp" || f.Name == "closeScreenWithPlainText" || f.Name == "set" {
		return false
	}
	return true
}

func (f *FuncCall) Signature() []ast.Signature {
	switch f.Name {
	case "root", "abs", "neg", "log", "exp", "round", "ceil", "floor",
		"sin", "cos", "tan", "asin", "acos", "atan", "degrees", "radians",
		"decToHex", "decToBin", "hexToDec", "binToDec":
		return []ast.Signature{ast.SignNumb}

	case "dec", "bin", "octal", "hexa":
		return []ast.Signature{ast.SignNumb}
	case "randInt":
		return []ast.Signature{ast.SignNumb}
	case "randFloat":
		return []ast.Signature{ast.SignNumb}
	case "setRandSeed":
		return []ast.Signature{ast.SignNumb}
	case "min", "max":
		return []ast.Signature{ast.SignNumb}
	case "avgOf", "maxOf", "minOf", "geoMeanOf", "stdDevOf", "stdErrOf":
		return []ast.Signature{ast.SignNumb}
	case "modeOf":
		return []ast.Signature{ast.SignNumb}
	case "mod", "rem", "quot":
		return []ast.Signature{ast.SignNumb}
	case "aTan2":
		return []ast.Signature{ast.SignNumb}
	case "formatDecimal":
		return []ast.Signature{ast.SignNumb}

	case "println":
		return []ast.Signature{ast.SignVoid}
	case "openScreen":
		return []ast.Signature{ast.SignVoid}
	case "openScreenWithValue":
		return []ast.Signature{ast.SignVoid}
	case "closeScreenWithValue":
		return []ast.Signature{ast.SignVoid}
	case "getStartValue":
		return []ast.Signature{ast.SignVoid}
	case "closeScreen":
		return []ast.Signature{ast.SignVoid}
	case "closeApp":
		return []ast.Signature{ast.SignVoid}
	case "getPlainStartText":
		return []ast.Signature{ast.SignText}
	case "closeScreenWithPlainText":
		return []ast.Signature{ast.SignText}
	case "copyList":
		return []ast.Signature{ast.SignList}
	case "copyDict":
		return []ast.Signature{ast.SignDict}

	case "makeColor":
		return []ast.Signature{ast.SignNumb}
	case "splitColor":
		return []ast.Signature{ast.SignList}

	case "set":
		return []ast.Signature{ast.SignVoid}
	case "get":
		return []ast.Signature{ast.SignAny}
	case "call":
		return []ast.Signature{ast.SignAny}
	case "every":
		return []ast.Signature{ast.SignList}
	default:
		panic("Unknown function " + f.Name)
	}
}

func (f *FuncCall) everyComponent() ast.Block {
	compType, ok := f.Args[0].(*variables.Get)
	if !ok || compType.Global {
		f.Where.Error("Expected a component type for every() 1st argument!")
	}
	return ast.Block{
		Type:     "component_all_component_block",
		Mutation: &ast.Mutation{ComponentType: compType.Name},
		Fields:   []ast.Field{{Name: "COMPONENT_SELECTOR", Value: compType.Name}},
	}
}

func (f *FuncCall) genericCall() ast.Block {
	// arg[0] 	 compType
	// arg[1] 	 component (any object)
	// arg[2] 	 method name
	// arg[4->n] invoke args
	compType, ok := f.Args[0].(*variables.Get)
	if !ok || compType.Global {
		f.Where.Error("Expected a component type for call() 1st argument!")
	}
	vGet, ok := f.Args[2].(*variables.Get)
	if !ok || vGet.Global {
		f.Where.Error("Expected a method name for call() 3rd argument!")
	}
	return ast.Block{
		Type: "component_method",
		Mutation: &ast.Mutation{
			MethodName:    vGet.Name,
			IsGeneric:     true,
			ComponentType: compType.Name,
		},
		Values: ast.ValueArgsByPrefix(f.Args[1], "COMPONENT", "ARG", f.Args[3:]),
	}
}

func (f *FuncCall) genericGet() ast.Block {
	f.assertArgLen(3)
	compType, ok := f.Args[0].(*variables.Get)
	if !ok || compType.Global {
		f.Where.Error("Expected a component type for get() 1st argument!")
	}
	vGet, ok := f.Args[2].(*variables.Get)
	if !ok || vGet.Global {
		f.Where.Error("Expected a property type for get() 3rd argument!")
	}
	return ast.Block{
		Type: "component_set_get",
		Mutation: &ast.Mutation{
			SetOrGet:      "get",
			PropertyName:  vGet.Name,
			IsGeneric:     true,
			ComponentType: compType.Name,
		},
		Fields: []ast.Field{{Name: "PROP", Value: vGet.Name}},
		Values: []ast.Value{{Name: "COMPONENT", Block: f.Args[1].Blockly(false)}},
	}
}

func (f *FuncCall) genericSet() ast.Block {
	f.assertArgLen(4)
	compType, ok := f.Args[0].(*variables.Get)
	if !ok || compType.Global {
		f.Where.Error("Expected a component type for set() 1st argument!")
	}
	vGet, ok := f.Args[2].(*variables.Get)
	if !ok || vGet.Global {
		f.Where.Error("Expected a property type for set() 3rd argument!")
	}
	return ast.Block{
		Type: "component_set_get",
		Mutation: &ast.Mutation{
			SetOrGet:      "set",
			PropertyName:  vGet.Name,
			IsGeneric:     true,
			ComponentType: compType.Name,
		},
		Fields: []ast.Field{{Name: "PROP", Value: vGet.Name}},
		Values: ast.MakeValues([]ast.Expr{f.Args[1], f.Args[3]}, "COMPONENT", "VALUE"),
	}
}

func (f *FuncCall) splitColor() ast.Block {
	return ast.Block{
		Type:   "color_make_color",
		Values: ast.MakeValues(f.Args, "COLOR"),
	}
}

func (f *FuncCall) makeColor() ast.Block {
	return ast.Block{
		Type:   "color_make_color",
		Values: ast.MakeValues(f.Args, "COLORLIST"),
	}
}

func (f *FuncCall) copyDict() ast.Block {
	return ast.Block{
		Type:   "dictionaries_copy",
		Values: ast.MakeValues(f.Args, "DICT"),
	}
}

func (f *FuncCall) copyList() ast.Block {
	return ast.Block{
		Type:   "lists_copy",
		Values: ast.MakeValues(f.Args, "LIST"),
	}
}

func (f *FuncCall) ctrlSimpleBlock(blockType string) ast.Block {
	return ast.Block{Type: blockType}
}

func (f *FuncCall) closeScreenWithPlainText() ast.Block {
	f.assertArgLen(1)
	return ast.Block{
		Type:   "controls_closeScreenWithPlainText",
		Values: ast.MakeValues(f.Args, "TEXT"),
	}
}

func (f *FuncCall) closeScreenWithValue() ast.Block {
	f.assertArgLen(1)
	return ast.Block{
		Type:   "controls_closeScreenWithValue",
		Values: ast.MakeValues(f.Args, "SCREEN"),
	}
}

func (f *FuncCall) openScreenWithValue() ast.Block {
	f.assertArgLen(2)
	return ast.Block{
		Type:   "controls_openAnotherScreenWithStartValue",
		Values: ast.MakeValues(f.Args, "SCREENNAME", "STARTVALUE"),
	}
}

func (f *FuncCall) openScreen() ast.Block {
	f.assertArgLen(1)
	return ast.Block{
		Type:   "controls_openAnotherScreen",
		Values: ast.MakeValues(f.Args, "SCREEN"),
	}
}

func (f *FuncCall) println() ast.Block {
	f.assertArgLen(1)
	return ast.Block{Type: "controls_eval_but_ignore", Values: ast.MakeValues(f.Args, "VALUE")}
}

var mathFuncMap = map[string]string{
	"sqrt":     "ROOT",
	"abs":      "ABS",
	"neg":      "NEG",
	"log":      "LN",
	"exp":      "EXP",
	"round":    "ROUND",
	"ceil":     "CEILING",
	"floor":    "FLOOR",
	"sin":      "SIN",
	"cos":      "COS",
	"tan":      "TAN",
	"asin":     "ASIN",
	"acos":     "ACOS",
	"atan":     "ATAN",
	"degrees":  "RADIANS_TO_DEGREES",
	"radians":  "DEGREES_TO_RADIANS",
	"decToHex": "DEC_TO_HEX",
	"decToBin": "DEC_TO_BIN",
	"hexToDec": "HEX_TO_DEC",
	"binToDec": "BIN_TO_DEC",
}

func (f *FuncCall) mathConversions() ast.Block {
	f.assertArgLen(1)
	fieldOp, ok := mathFuncMap[f.Name]
	if !ok {
		f.Where.Error("Unknown Math Conversion %()", f.Name)
	}
	var blockType string
	switch fieldOp {
	case "SIN", "COS", "TAN", "ASIN", "ACOS", "ATAN":
		blockType = "math_trig"
	case "RADIANS_TO_DEGREES", "DEGREES_TO_RADIANS":
		blockType = "math_convert_angles"
	case "DEC_TO_HEX", "HEX_TO_DEC", "DEC_TO_BIN", "BIN_TO_DEC":
		blockType = "math_convert_number"
	default:
		blockType = "math_single"
	}
	return ast.Block{
		Type:   blockType,
		Fields: []ast.Field{{Name: "OP", Value: fieldOp}},
		Values: []ast.Value{{Name: "NUM", Block: f.Args[0].Blockly(false)}},
	}
}

func (f *FuncCall) formatDecimal() ast.Block {
	f.assertArgLen(2)
	return ast.Block{
		Type:   "math_format_as_decimal",
		Values: ast.MakeValues(f.Args, "NUM", "PLACES"),
	}
}

func (f *FuncCall) atan2() ast.Block {
	f.assertArgLen(2)
	return ast.Block{
		Type:   "math_atan2",
		Values: ast.MakeValues(f.Args, "Y", "X"),
	}
}

func (f *FuncCall) mathDivide() ast.Block {
	f.assertArgLen(2)
	var fieldOp string
	switch f.Name {
	case "mod":
		fieldOp = "MODULO"
	case "rem":
		fieldOp = "REMAINDER"
	case "quot":
		fieldOp = "QUOTIENT"
	}
	return ast.Block{
		Type:   "math_divide",
		Fields: []ast.Field{{Name: "OP", Value: fieldOp}},
		Values: ast.MakeValues(f.Args, "DIVIDEND", "DIVISOR"),
	}
}

func (f *FuncCall) modeOf() ast.Block {
	f.assertArgLen(1)
	return ast.Block{
		Type:   "math_mode_of_list",
		Values: ast.MakeValues(f.Args, "LIST"),
	}
}

func (f *FuncCall) mathOnList() ast.Block {
	f.assertArgLen(1)
	var fieldOp string
	switch f.Name {
	case "avgOf":
		fieldOp = "AVG"
	case "maxOf":
		fieldOp = "MAX"
	case "minOf":
		fieldOp = "MIN"
	case "geoMeanOf":
		fieldOp = "GM"
	case "stdDevOf":
		fieldOp = "SD"
	case "stdErrOf":
		fieldOp = "SE"
	}
	return ast.Block{
		Type:   "math_on_list2",
		Fields: []ast.Field{{Name: "OP", Value: fieldOp}},
		Values: ast.MakeValues(f.Args, "LIST"),
	}
}

func (f *FuncCall) minOrMax() ast.Block {
	argSize := len(f.Args)
	if argSize == 0 {
		f.Where.Error("No arguments provided for %()", f.Name)
	}
	var fieldOp string
	switch f.Name {
	case "min":
		fieldOp = "MIN"
	case "max":
		fieldOp = "MAX"
	}
	return ast.Block{
		Type:     "math_on_list",
		Fields:   []ast.Field{{Name: "OP", Value: fieldOp}},
		Mutation: &ast.Mutation{ItemCount: argSize},
		Values:   ast.ValuesByPrefix("NUM", f.Args),
	}
}

func (f *FuncCall) setRandSeed() ast.Block {
	f.assertArgLen(1)
	return ast.Block{
		Type:   "math_random_set_seed",
		Values: ast.MakeValues(f.Args, "NUM"),
	}
}

func (f *FuncCall) randFloat() ast.Block {
	f.assertArgLen(0)
	return ast.Block{Type: "math_random_float"}
}

func (f *FuncCall) randInt() ast.Block {
	f.assertArgLen(2)
	return ast.Block{
		Type:   "math_random_int",
		Values: ast.MakeValues(f.Args, "FROM", "TO"),
	}
}

func (f *FuncCall) mathRadix() ast.Block {
	f.assertArgLen(1)
	var fieldOp string
	switch f.Name {
	case "dec":
		fieldOp = "DEC"
	case "bin":
		fieldOp = "BIN"
	case "octal":
		fieldOp = "OCT"
	case "hexa":
		fieldOp = "HEX"
	}
	textExpr, ok := f.Args[0].(*fundamentals.Text)
	if !ok {
		f.Where.Error("Expected a numeric string argument for %()", f.Name)
	}
	return ast.Block{
		Type: "math_number_radix",
		Fields: []ast.Field{
			{Name: "OP", Value: fieldOp},
			{Name: "NUM", Value: textExpr.Content},
		},
	}
}

func (f *FuncCall) assertArgLen(expectLen int) {
	argsLen := len(f.Args)
	if argsLen != expectLen {
		f.Where.Error("Expected % argument for %() but got %", strconv.Itoa(expectLen), f.Name, strconv.Itoa(argsLen))
	}
}
