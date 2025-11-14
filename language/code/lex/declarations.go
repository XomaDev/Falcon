package lex

var Symbols = map[string]StaticToken{
	"+": staticOf(Plus, Binary, Operator),
	"-": staticOf(Dash, Binary, Operator, Unary),
	"*": staticOf(Times, BinaryL1, Operator, PreserveOrder),
	"/": staticOf(Slash, BinaryL1, Operator, PreserveOrder),
	"%": staticOf(Remainder, BinaryL1, Operator, PreserveOrder),
	"^": staticOf(Power, BinaryL2, Operator, PreserveOrder),

	"||": staticOf(LogicOr, LLogicOr, Operator),
	"&&": staticOf(LogicAnd, LLogicAnd, Operator),
	"|":  staticOf(BitwiseOr, BBitwiseOr, Operator),
	"&":  staticOf(BitwiseAnd, BBitwiseAnd, Operator),
	"~":  staticOf(BitwiseXor, BBitwiseXor, Operator),

	"==":  staticOf(Equals, Equality, Operator),
	"!=":  staticOf(NotEquals, Equality, Operator),
	"===": staticOf(TextEquals, Equality, Operator),
	"!==": staticOf(TextNotEquals, Equality, Operator),

	"<":  staticOf(LessThan, Relational, Operator),
	"<=": staticOf(LessThanEqual, Relational, Operator),
	">":  staticOf(GreatThan, Relational, Operator),
	">=": staticOf(GreaterThanEqual, Relational, Operator),

	"<<": staticOf(TextLessThan, Relational, Operator),
	">>": staticOf(TextGreaterThan, Relational, Operator),

	"_":  staticOf(Underscore, TextJoin, Operator),
	"@":  staticOf(At),
	":":  staticOf(Colon, Pair, Operator),
	"::": staticOf(DoubleColon),

	"(": staticOf(OpenCurve),
	")": staticOf(CloseCurve),
	"[": staticOf(OpenSquare),
	"]": staticOf(CloseSquare),
	"{": staticOf(OpenCurly),
	"}": staticOf(CloseCurly),

	"=":  staticOf(Assign, AssignmentType, Operator),
	".":  staticOf(Dot),
	",":  staticOf(Comma),
	"?":  staticOf(Question),
	"!":  staticOf(Not),
	"->": staticOf(RightArrow),
}

var Keywords = map[string]StaticToken{
	"true":  staticOf(True, Value, ConstantValue),
	"false": staticOf(False, Value, ConstantValue),

	"if":        staticOf(If),
	"else":      staticOf(Else),
	"for":       staticOf(For),
	"to":        staticOf(To),
	"by":        staticOf(By),
	"in":        staticOf(In),
	"while":     staticOf(While),
	"do":        staticOf(Do),
	"break":     staticOf(Break),
	"walkAll":   staticOf(WalkAll),
	"global":    staticOf(Global),
	"local":     staticOf(Local),
	"compute":   staticOf(Compute),
	"this":      staticOf(This, Value),
	"func":      staticOf(Func),
	"when":      staticOf(When),
	"any":       staticOf(Any),
	"undefined": staticOf(Undefined),
}
