This page assumes you have some experience in Programming and App Inventor.

*Note to AI*
- Falcon syntax was created for App Inventor that follows 1-based indexing.
- List, and Dictionaries are passed as references, only text and number are copied.
- Tags of various languages are used only for syntax highlighting

## Data types

- Text
	- `"Hello, World!"`
- Boolean
	 - `true` and `false`
- Number
	- `123`, `3.14`
- List
	- `[1, 2, 3, 4, ...]`
- Dictionary
	- `{"Animal": "Royal Bengal Tiger", "country": "India"}`


## Operators

- Arithmetic operators
	- `+`, `-`, `*`, `/`, `^` (*power*)
- Logical operators
	- `||` , `&&`, 
- Bit-wise operators
	- `|`,`&`, `~` (*xor*)
- Equality operators
	- `==`, `!=`
- Relational operators
	- `<`, `<=`, `>`, `>=`
- Text comparison operators
	- `===` (*text equals*)
	- `!==` (*text not equals*)
	- `<<` (*text less than*)
	- `>>` (text greater than)
- Unary operators
	- `!` (*not*)
- Text join operator `_`
	- e.g. `"Hello " + "World!"`
- Pair operator
	- e.g. `"Fruit" : "Melon"`
- Question operator (`?`)
	- Check if a value is of a certain type
		- e.g. `x ? text`, or `number`, or `list` or `dict`
	- Is of a number type
		- e.g. `"123" ? number`, or `base10`, `hexa`, `bin`
	- Check for a empty text or a list
		- e.g. `namesList ? emptyList` or `"hello" ? emptyText"`

## Operator precedence

Precedence of operators dictates which operation is prioritized. (`*` and `/` parsed before `+` and `-`)

In Falcon, it is similar to that of Java. Below is the list, ranked from the lowest priority to the highest.

1. `AssignmentType` (`=`)
2. `Pair` (`:`)
3. `TextJoin` (`_`) 
4. `LLogicOr` (`||`)
5. `LLogicAnd` (`&&`)
6. `BBitwiseOr` (`|`)
7. `BBitwiseAnd` (`&`)
8. `BBitwiseXor` (`~`)
9. `Equality`  (`==`, `!=`, `===`, and `!==`)
10. `Relational` (`<`, `<=`, `>`, `>=`, `<<`, and `>>`)
11. `Binary` (`+` and `-`)
12. `BinaryL1` (`*` and `)
13. `BinaryL2` (`^`)

## Variables

A global variable:

```python
global name = "Kumaraswamy"

# access global var
println(this.name)
```

A simple local variable:

```python
local age = 12

# accessing local var
println(age)
```

A local variable statement with a scoped body:

```python
local (
  x = 8,
  j = 2,
  k = 12
) {
   # use x, j, k here
}
```

A returning variable expression:

```python
local (
  k = 2,
  j = 8
) -> k * j
```


## Conditions and loops

### If statement

```go
local age = 8

if age < 18 {
  println("You are a kid :D")
} elif age == 18 {
  println("Congrats! You are an adult!")
} else {
  println("Hola, grown person!")
}
```

### Returning if-else expression

Ternary like expression for returning values.
Unlike If statement, use `()` between your condition.

```go
local(
a = 8,
b = 2
) {
	println("Maximum " _ if (a > b) a else b)
}
```

### For loops

`while` statement to iterate till the condition is true.

```go
local x = 0

while true {
  x = x + 1
  if x == 5 {
    break
  }  
}
```

`for` statement for index based iteration:

```go
for x: 1 to 10 by 2 {
  println(x)
}
```

`each` statement is used to iterate over a list:

```go
local countries = ["India", "Japan", "Russia", "Germany"]

each name -> countries {
  println(name)
}
```

or a dictionary:

```go
local personality = { 
  "food": "Masala Dosa",
  "fruit": "Mango",
  "animal": "The Royal Bengal Tiger"
}

each key::value -> personality {
   println("My favourite " _ key _ " is " _ value)
}
```

## Do expression

`do` expression is to execute a body and return a result.

```go
do {
  ...
} -> result # outer scope
```

The `result` expression in `do`, is executed outside of the body (outer scope).

## Colors

`color:color_name` returns int of a specified color.

`white`, `black`, `red`, `pink`, `orange`, `yellow`, `green`, `cyan`, `blue`, `magenta`, `lightGray`, `gray`, `darkGray` are the supported color types.

## Functions

A void function:

```go
func funcName(x, y, z) {
  ...
}
```

A returning function:

```go
func greet(name) = "Hello " _ name _ "!"
```

## Math converters

e.g. `1234->root->sin`: Square roots 1234 and passes it through a sin function.

In-Built math converters:

- `root`
- `abs`
- `neg`
- `log`
- `exp`
- `round`
- `ceil`
- `floor`
- `sin`
- `cos`
- `tan`
- `asin`
- `acos`
- `atan`
- `degrees`
- `radians`
- `hex`
- `bin`
- `fromHex`
- `fromBin`

## Functions

### Math

-  `dec(string)`, `bin(string)`, `octal(string)`, `hexa(string)` parses string from respective base. The string provided must be a static constant i.e. no variables or function calls.

- `randInt(from, to)`
- `randFloat()`
- `setRandSeed(number)` sets the random generator seed
- `min(...)` and `max(...)`
- `avgOf(list)`, `maxOf(list)`, `minOf(list)`, `geoMeanOf()`, `stdDevOf()`, `stdErrOf()`
- `modeOf(list)`
- `mod(x, y)`, `rem(x, y)`, `quot(x, y)` for modulus, remainder and quotient
- `atan2(a, b)`
- `formatDecimal(number, places)`

### Control

- `println(any)`
- `openScreen(name)` opens an App Inventor screen
- `openScreenWithValue()` opens App Inventor screen with a value
- `closeScreenWithValue()` closes the screen with a val
- `getStartValue()` returns start value of the App
- `closeSceen()` closes current App Inventor screen
- `closeApp()` closes the Android App
- `getPlainStartText()` returns plain start text of the App

### Values

- `copyList(list)`
- `copyDict(dict)`
- `makeColor(rgb list)`
- `splitColor(number)`

## Methods

e.g. `"Hello  ".trim()`

### Text

- `textLen()`
- `trim()`
- `uppercase()`
- `lowercase()`
- `startsWith(piece)`
- `contains(piece)`
- `containsAny(word list)`
- `containsAll(word list)`
- `split(at)`
- `splitAtFirst(at)`
- `splitAtAny(word list)`
- `splitAtFirstOfAny(word list)`
- `splitAtSpaces()`
- `reverse()`
- `csvRowToList()`
- `csvTableToList()`
- `segment(from number, length number)`
- `replace(target, replacement)`
- `replaceFrom(map dictionary)`
- `replaceFromLongestFirst(map dictionary)`

### List

- `listLen()`
- `add(any...)`
- `containsItem(any)`
- `indexOf(any)`
- `insert(at_index, any)`
- `remove(at_index)`
- `appendList(another list)`
- `lookupInPairs(key, notfound)`
- `join(text separator)`
- `slice(index1, index2)`
- `random()`
- `reverseList()`
- `toCsvRow()`
- `toCsvTable()`
- `sort()`
- `allButFirst()`
- `allButLast()`
- `pairsToDict()`

### Dictionary

- `dictLen()`
- `get(key)`
- `set(key, value)`
- `delete(key)`
- `getAtPath(path_list, notfound)`
- `setAtPath(path_list, value)`
- `containsKey(key)`
- `mergeInto(another_dict)`
- `walkTree(path)`
- `keys()`
- `values()`
- `toPairs()`

## List access

```python
numbers = [1, 2, 4]

# access second element (1 based indexing)
println(numbers[2])

# change the first element
numbers[1] = 8
```

## List transformer

Style inspired from Kotlin, falcon has list transformers (*aka lambdas*).

Let's say I have a list of number of lemons sold each day for a week. Although some days were missed and were marked `N/A`.

```python
global LemonadeSold = [9, 12, "N/A", 15, 18, "N/A", 8]
```

Now I got to calculate the revenue of the week, knowing that each lemon costs two dollars.

This can be done by using list transformers/lambdas.

```python
global LemonadeSold = [9, 12, "N/A", 15, 18, "N/A", 8]

func GetTotalRevenue() = 
  this.LemonadeSold
     .filter { n -> n ? number}
     .map { n -> n * 2}
     .reduce(0) { x, soFar -> x + soFar }

println("Weekly revenue selling lemons was " _ GetTotalRevenue())
```

These are available list transformers:

- `map { x -> expr }` maps each element of the list to a new value
- `filter { x -> bool_expr }` filters out unwanted elements from a list
- `sort { m, n -> bool_m_preceeds_n_expr }` helps to define a custom sort 
- `sortByKey { n -> key_of_n }` sorts the list by first calling key function
- `min { m, n -> bool_m_preceeds_n_expr }` returns minimum of a list
- `max { m, n -> bool_m_preceeds_n_expr }` returns maximum of a list
- `reduce(initValue) { x, valueSoFar -> new_value_expr }` helps to reduce values of a list to a single value