package lex

var Symbols = map[string]StaticToken{
	"+": staticOf(Plus, Binary, Operator),
	"-": staticOf(Dash, Binary, Operator, Unary),
	"*": staticOf(Times, BinaryL1, Operator, PreserveOrder),
	"/": staticOf(Slash, BinaryL1, Operator, PreserveOrder),
	"^": staticOf(Power, BinaryL2, Operator, PreserveOrder),

	"||": staticOf(LogicOr, LLogicOr, Operator),
	"&&": staticOf(LogicAnd, LLogicAnd, Operator),
	"|":  staticOf(BitwiseOr, BBitwiseOr, Operator),
	"&":  staticOf(BitwiseAnd, BBitwiseAnd, Operator),
	"~":  staticOf(BitwiseXor, BBitwiseXor, Operator),

	"==": staticOf(Equals, Equality, Operator),
	"!=": staticOf(NotEquals, Equality, Operator),

	"<":  staticOf(LessThan, Relational, Operator),
	"<=": staticOf(LessThanEqual, Relational, Operator),
	">":  staticOf(GreatThan, Relational, Operator),
	">=": staticOf(GreaterThanEqual, Relational, Operator),

	"(": staticOf(OpenCurve),
	")": staticOf(CloseCurve),
	"[": staticOf(OpenSquare),
	"]": staticOf(CloseSquare),
	"{": staticOf(OpenCurly),
	"}": staticOf(CloseCurly),

	"=": staticOf(Assign),
	".": staticOf(Dot),
	",": staticOf(Comma),
	"?": staticOf(Question),
	"!": staticOf(Not),
}

var Keywords = map[string]StaticToken{
	"true":  staticOf(True, Value),
	"false": staticOf(False, Value),

	"if":   staticOf(If),
	"elif": staticOf(Elif),
	"else": staticOf(Else),
}
