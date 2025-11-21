![image|690x338](https://community.appinventor.mit.edu/uploads/default/original/3X/e/d/ed3b9d22ddaefffb4fd5ab71964b8816c56c63a1.png)

Falcon is a language designed for App Inventor to enable syntax-based programming and for incorporating agenting coding abilities.

## Quirks
1. Falcon follows 1-based indexing.
2. Falcon is a type-blind language.
3. Lists and dictionaries are passed as references.
4. Falcon follows Kotlin's style of functional expressions.
5. Falcon does not have a return statement; the last expression in a body is returned.
6. Only single-line comments using double slash `//` are supported.
7. Do not use `_` in place of unused variables
8. Variables can never be in the uninitialized state.
9. Always keep the last expression for returning functions.
10. Always generate codes between ``` and code tag falcon.

## Data types
1. String `<block type="text">
    <field name="TEXT">Hello, world!</field>
  </block>`
2. Boolean `<block type="logic_boolean">
    <field name="BOOL">TRUE</field>
  </block>` and `<block type="logic_boolean">
    <field name="BOOL">FALSE</field>
  </block>`
3. Number `block type="math_number">
    <field name="NUM">123</field>
  </block>` and `<block type="math_number">
    <field name="NUM">3.14</field>
  </block>`
4. List `<block type="lists_create_with">
    <mutation items="4" elseif="0" else="0"></mutation>
    <value name="ADD0">
      <block type="math_number">
        <field name="NUM">1</field>
      </block>
    </value>
    <value name="ADD1">
      <block type="math_number">
        <field name="NUM">2</field>
      </block>
    </value>
    <value name="ADD2">
      <block type="math_number">
        <field name="NUM">3</field>
      </block>
    </value>
    <value name="ADD3">
      <block type="math_number">
        <field name="NUM">4</field>
      </block>
    </value>
  </block>`
5. Dictionary `<block type="dictionaries_create_with">
    <mutation items="2" elseif="0" else="0"></mutation>
    <value name="ADD0">
      <block type="pair">
        <value name="KEY">
          <block type="text">
            <field name="TEXT">Animal</field>
          </block>
        </value>
        <value name="VALUE">
          <block type="text">
            <field name="TEXT">Tiger</field>
          </block>
        </value>
      </block>
    </value>
    <value name="ADD1">
      <block type="pair">
        <value name="KEY">
          <block type="text">
            <field name="TEXT">Scientific Name</field>
          </block>
        </value>
        <value name="VALUE">
          <block type="text">
            <field name="TEXT">Panthera tigris</field>
          </block>
        </value>
      </block>
    </value>
  </block>`
6. Colour `<block type="color_white">
    <field name="COLOR">#FFFFFF</field>
  </block>`

## Operators

Adding two numbers
```xml
  <block type="math_add">
    <mutation items="2" elseif="0" else="0"></mutation>
    <value name="NUM0">
      <block type="math_number">
        <field name="NUM">1</field>
      </block>
    </value>
    <value name="NUM1">
      <block type="math_number">
        <field name="NUM">2</field>
      </block>
    </value>
  </block>
```

Subtracting two numbers
```xml
  <block type="math_subtract">
    <value name="A">
      <block type="math_number">
        <field name="NUM">1</field>
      </block>
    </value>
    <value name="B">
      <block type="math_number">
        <field name="NUM">2</field>
      </block>
    </value>
  </block>
```

Multiplying two numbers
```xml
  <block type="math_multiply">
    <mutation items="2" elseif="0" else="0"></mutation>
    <value name="NUM0">
      <block type="math_number">
        <field name="NUM">1</field>
      </block>
    </value>
    <value name="NUM1">
      <block type="math_number">
        <field name="NUM">2</field>
      </block>
    </value>
  </block>
```

Dividing two numbers
```xml
  <block type="math_division">
    <value name="A">
      <block type="math_number">
        <field name="NUM">1</field>
      </block>
    </value>
    <value name="B">
      <block type="math_number">
        <field name="NUM">2</field>
      </block>
    </value>
  </block>
```

Remainder operator
```xml
  <block type="math_divide">
    <field name="OP">REMAINDER</field>
    <value name="DIVIDEND">
      <block type="math_number">
        <field name="NUM">8</field>
      </block>
    </value>
    <value name="DIVISOR">
      <block type="math_number">
        <field name="NUM">2</field>
      </block>
    </value>
  </block>
```

Power operator
```xml
  <block type="math_power">
    <value name="A">
      <block type="math_number">
        <field name="NUM">2</field>
      </block>
    </value>
    <value name="B">
      <block type="math_number">
        <field name="NUM">3</field>
      </block>
    </value>
  </block>
```

Logical &&
```xml
  <block type="logic_operation">
    <mutation items="2" elseif="0" else="0"></mutation>
    <field name="OP">AND</field>
    <value name="A">
      <block type="logic_boolean">
        <field name="BOOL">TRUE</field>
      </block>
    </value>
    <value name="B">
      <block type="logic_boolean">
        <field name="BOOL">FALSE</field>
      </block>
    </value>
  </block>
```

Logical OR
```xml
  <block type="logic_operation">
    <mutation items="2" elseif="0" else="0"></mutation>
    <field name="OP">OR</field>
    <value name="A">
      <block type="logic_boolean">
        <field name="BOOL">TRUE</field>
      </block>
    </value>
    <value name="B">
      <block type="logic_boolean">
        <field name="BOOL">FALSE</field>
      </block>
    </value>
  </block>
```

Bitwise &
```xml
  <block type="math_bitwise">
    <mutation items="2" elseif="0" else="0"></mutation>
    <field name="OP">BITAND</field>
    <value name="NUM0">
        <block type="math_number">
            <field name="NUM">23</field>
        </block>
    </value>
    <value name="NUM1">
        <block type="math_number">
            <field name="NUM">21</field>
        </block>
    </value>
</block>
```
Bitwise |
```xml
  <block type="math_bitwise">
    <mutation items="2" elseif="0" else="0"></mutation>
    <field name="OP">BITOR</field>
    <value name="NUM0">
      <block type="math_number">
        <field name="NUM">23</field>
      </block>
    </value>
    <value name="NUM1">
      <block type="math_number">
        <field name="NUM">21</field>
      </block>
    </value>
  </block>
```

Equals
```xml
  <block type="logic_compare">
    <field name="OP">EQ</field>
    <value name="A">
      <block type="math_number">
        <field name="NUM">123</field>
      </block>
    </value>
    <value name="B">
      <block type="math_number">
        <field name="NUM">8</field>
      </block>
    </value>
  </block>
```

Not equals
```xml
  <block type="logic_compare">
    <field name="OP">NEQ</field>
    <value name="A">
      <block type="math_number">
        <field name="NUM">123</field>
      </block>
    </value>
    <value name="B">
      <block type="math_number">
        <field name="NUM">8</field>
      </block>
    </value>
  </block>
```

Less than <
```xml
  <block type="math_compare">
    <field name="OP">LT</field>
    <value name="A">
        <block type="math_number">
            <field name="NUM">123</field>
        </block>
    </value>
    <value name="B">
        <block type="math_number">
            <field name="NUM">8</field>
        </block>
    </value>
</block>
```

Greater than >
```xml
  <block type="math_compare">
    <field name="OP">GT</field>
    <value name="A">
      <block type="math_number">
        <field name="NUM">123</field>
      </block>
    </value>
    <value name="B">
      <block type="math_number">
        <field name="NUM">8</field>
      </block>
    </value>
  </block>
```

Less than or equals <=
```xml
  <block type="math_compare">
    <field name="OP">LT</field>
    <value name="A">
      <block type="math_number">
        <field name="NUM">123</field>
      </block>
    </value>
    <value name="B">
      <block type="math_number">
        <field name="NUM">8</field>
      </block>
    </value>
  </block>
```

Greater than or equals
```xml
  <block type="math_compare">
    <field name="OP">GTE</field>
    <value name="A">
      <block type="math_number">
        <field name="NUM">123</field>
      </block>
    </value>
    <value name="B">
      <block type="math_number">
        <field name="NUM">8</field>
      </block>
    </value>
  </block>
```

Text lexicographic equals ===
```xml
  <block type="text_compare">
    <field name="OP">EQ</field>
    <value name="TEXT1">
      <block type="text">
        <field name="TEXT">cat</field>
      </block>
    </value>
    <value name="TEXT2">
      <block type="text">
        <field name="TEXT">Cat</field>
      </block>
    </value>
  </block>
```

Text lexicographic not equals !==
```xml
  <block type="text_compare">
    <field name="OP">NEQ</field>
    <value name="TEXT1">
      <block type="text">
        <field name="TEXT">cat</field>
      </block>
    </value>
    <value name="TEXT2">
      <block type="text">
        <field name="TEXT">Cat</field>
      </block>
    </value>
  </block>
```

Text lexicographic less than <<
```xml
  <block type="text_compare">
    <field name="OP">LT</field>
    <value name="TEXT1">
      <block type="text">
        <field name="TEXT">cat</field>
      </block>
    </value>
    <value name="TEXT2">
      <block type="text">
        <field name="TEXT">Cat</field>
      </block>
    </value>
  </block>
```

Text lexicographic greater than >>
```xml
  <block type="text_compare">
    <field name="OP">GT</field>
    <value name="TEXT1">
      <block type="text">
        <field name="TEXT">cat</field>
      </block>
    </value>
    <value name="TEXT2">
      <block type="text">
        <field name="TEXT">Cat</field>
      </block>
    </value>
  </block>
```

Unary Not
```xml
  <block type="logic_negate">
    <value name="BOOL">
      <block type="logic_boolean">
        <field name="BOOL">TRUE</field>
      </block>
    </value>
  </block>
```

Math Negate
```xml
  <block type="math_single">
    <field name="OP">NEG</field>
    <value name="NUM">
      <block type="math_number">
        <field name="NUM">123</field>
      </block>
    </value>
  </block>
```

Join
```xml
  <block type="text_join">
    <mutation items="2" elseif="0" else="0"></mutation>
    <value name="ADD0">
      <block type="text">
        <field name="TEXT">Hello </field>
      </block>
    </value>
    <value name="ADD1">
      <block type="text">
        <field name="TEXT">World!</field>
      </block>
    </value>
  </block>
```

Pair
```xml
  <block type="pair">
    <value name="KEY">
      <block type="text">
        <field name="TEXT">Fruit</field>
      </block>
    </value>
    <value name="VALUE">
      <block type="text">
        <field name="TEXT">Mango</field>
      </block>
    </value>
  </block>
```

Operator precedence
The precedence of an operator dictates its parse order. E.g. `*` and `/` is parsed before `+` and `-`.

It is similar to that of Java. Below is a ranking from the lowest to the highest precedence.

1. Assignment `=`
2. Pair `:`
3. TextJoin `_`
4. LogicOr `||`
5. LogicAnd `&&`
6. BitwiseOr `|`
7. BitwiseAnd `&`
8. BitwiseXor `~`
9. Equality `==`, `!=`, `===`, and `!==`
10. Relational `<`, `<=`, `>`, `>=`, `<<`, and `>>`
11. Binary `+`, and `-`
12. BinaryL1 `*`, `/`, and `%`
13. BinaryL2 `^`


## Variables

### Global variable

A global variable is always declared at the root:
```
  <block type="global_declaration">
    <field name="NAME">name</field>
    <value name="VALUE">
      <block type="text">
        <field name="TEXT">Kumaraswamy B G</field>
      </block>
    </value>
  </block>
  <block type="controls_eval_but_ignore">
    <value name="VALUE">
      <block type="lexical_variable_get">
        <field name="VAR">global name</field>
      </block>
    </value>
  </block>
```

### Local variable

```
  <block type="local_declaration_statement">
    <mutation items="0" elseif="0" else="0">
      <localname name="age"></localname>
    </mutation>
    <field name="VAR0">age</field>
    <value name="DECL0">
      <block type="math_number">
        <field name="NUM">17</field>
      </block>
    </value>
    <statement name="STACK">
      <block type="controls_eval_but_ignore">
        <value name="VALUE">
          <block type="lexical_variable_get">
            <field name="VAR">age</field>
          </block>
        </value>
      </block>
    </statement>
  </block>
```

## If else

If-else can be a statement or an expression depending on the context.

```
  <block type="local_declaration_statement">
    <mutation items="0" elseif="0" else="0">
      <localname name="x"></localname>
      <localname name="y"></localname>
    </mutation>
    <field name="VAR0">x</field>
    <field name="VAR1">y</field>
    <value name="DECL0">
      <block type="math_number">
        <field name="NUM">8</field>
      </block>
    </value>
    <value name="DECL1">
      <block type="math_number">
        <field name="NUM">12</field>
      </block>
    </value>
    <statement name="STACK">
      <block type="controls_if">
        <mutation items="0" elseif="1" else="1"></mutation>
        <value name="IF0">
          <block type="math_compare">
            <field name="OP">GT</field>
            <value name="A">
              <block type="lexical_variable_get">
                <field name="VAR">x</field>
              </block>
            </value>
            <value name="B">
              <block type="lexical_variable_get">
                <field name="VAR">y</field>
              </block>
            </value>
          </block>
        </value>
        <value name="IF1">
          <block type="math_compare">
            <field name="OP">GT</field>
            <value name="A">
              <block type="lexical_variable_get">
                <field name="VAR">y</field>
              </block>
            </value>
            <value name="B">
              <block type="lexical_variable_get">
                <field name="VAR">x</field>
              </block>
            </value>
          </block>
        </value>
        <statement name="DO0">
          <block type="controls_eval_but_ignore">
            <value name="VALUE">
              <block type="text">
                <field name="TEXT">X is greater</field>
              </block>
            </value>
          </block>
        </statement>
        <statement name="DO1">
          <block type="controls_eval_but_ignore">
            <value name="VALUE">
              <block type="text">
                <field name="TEXT">Y is greater</field>
              </block>
            </value>
          </block>
        </statement>
        <statement name="ELSE">
          <block type="controls_eval_but_ignore">
            <value name="VALUE">
              <block type="text">
                <field name="TEXT">They both are equal!</field>
              </block>
            </value>
          </block>
        </statement>
      </block>
    </statement>
  </block>
```

Used an expression:

```
  <block type="controls_eval_but_ignore">
    <value name="VALUE">
      <block type="controls_choose">
        <value name="TEST">
          <block type="math_compare">
            <field name="OP">GT</field>
            <value name="A">
              <block type="lexical_variable_get">
                <field name="VAR">x</field>
              </block>
            </value>
            <value name="B">
              <block type="lexical_variable_get">
                <field name="VAR">y</field>
              </block>
            </value>
          </block>
        </value>
        <value name="THENRETURN">
          <block type="text">
            <field name="TEXT">X is greater</field>
          </block>
        </value>
        <value name="ELSERETURN">
          <block type="controls_choose">
            <value name="TEST">
              <block type="math_compare">
                <field name="OP">GT</field>
                <value name="A">
                  <block type="lexical_variable_get">
                    <field name="VAR">y</field>
                  </block>
                </value>
                <value name="B">
                  <block type="lexical_variable_get">
                    <field name="VAR">x</field>
                  </block>
                </value>
              </block>
            </value>
            <value name="THENRETURN">
              <block type="text">
                <field name="TEXT">Y is greater</field>
              </block>
            </value>
            <value name="ELSERETURN">
              <block type="text">
                <field name="TEXT">They both are equal!</field>
              </block>
            </value>
          </block>
        </value>
      </block>
    </value>
  </block>
```


## While loop

```
  <block type="local_declaration_statement">
    <mutation items="0" elseif="0" else="0">
      <localname name="x"></localname>
    </mutation>
    <field name="VAR0">x</field>
    <value name="DECL0">
      <block type="math_number">
        <field name="NUM">8</field>
      </block>
    </value>
    <statement name="STACK">
      <block type="controls_while">
        <value name="TEST">
          <block type="logic_boolean">
            <field name="BOOL">TRUE</field>
          </block>
        </value>
        <statement name="DO">
          <block type="lexical_variable_set">
            <field name="VAR">x</field>
            <value name="VALUE">
              <block type="math_add">
                <mutation items="2" elseif="0" else="0"></mutation>
                <value name="NUM0">
                  <block type="lexical_variable_get">
                    <field name="VAR">x</field>
                  </block>
                </value>
                <value name="NUM1">
                  <block type="math_number">
                    <field name="NUM">1</field>
                  </block>
                </value>
              </block>
            </value>
            <next>
              <block type="controls_if">
                <mutation items="0" elseif="0" else="0"></mutation>
                <value name="IF0">
                  <block type="logic_compare">
                    <field name="OP">EQ</field>
                    <value name="A">
                      <block type="lexical_variable_get">
                        <field name="VAR">x</field>
                      </block>
                    </value>
                    <value name="B">
                      <block type="math_number">
                        <field name="NUM">5</field>
                      </block>
                    </value>
                  </block>
                </value>
                <statement name="DO0">
                  <block type="controls_break"></block>
                </statement>
              </block>
            </next>
          </block>
        </statement>
      </block>
    </statement>
  </block>
```

## For n loop

```
  <block type="controls_forRange">
    <field name="VAR">i</field>
    <value name="START">
      <block type="math_number">
        <field name="NUM">1</field>
      </block>
    </value>
    <value name="END">
      <block type="math_number">
        <field name="NUM">10</field>
      </block>
    </value>
    <value name="STEP">
      <block type="math_number">
        <field name="NUM">2</field>
      </block>
    </value>
    <statement name="DO">
      <block type="controls_eval_but_ignore">
        <value name="VALUE">
          <block type="lexical_variable_get">
            <field name="VAR">i</field>
          </block>
        </value>
      </block>
    </statement>
  </block>
```

The `by` clause is optional and defaults to 1.


## Each loop

To iterate over a list:

```
  <block type="local_declaration_statement">
    <mutation items="0" elseif="0" else="0">
      <localname name="names"></localname>
    </mutation>
    <field name="VAR0">names</field>
    <value name="DECL0">
      <block type="lists_create_with">
        <mutation items="4" elseif="0" else="0"></mutation>
        <value name="ADD0">
          <block type="text">
            <field name="TEXT">India</field>
          </block>
        </value>
        <value name="ADD1">
          <block type="text">
            <field name="TEXT">Japan</field>
          </block>
        </value>
        <value name="ADD2">
          <block type="text">
            <field name="TEXT">Russia</field>
          </block>
        </value>
        <value name="ADD3">
          <block type="text">
            <field name="TEXT">Germany</field>
          </block>
        </value>
      </block>
    </value>
    <statement name="STACK">
      <block type="controls_forEach">
        <field name="VAR">country</field>
        <value name="LIST">
          <block type="lexical_variable_get">
            <field name="VAR">names</field>
          </block>
        </value>
        <statement name="DO">
          <block type="controls_eval_but_ignore">
            <value name="VALUE">
              <block type="lexical_variable_get">
                <field name="VAR">country</field>
              </block>
            </value>
          </block>
        </statement>
      </block>
    </statement>
  </block>
```

Or over a dictionary:

```
  <block type="local_declaration_statement">
    <mutation items="0" elseif="0" else="0">
      <localname name="animalInfo"></localname>
    </mutation>
    <field name="VAR0">animalInfo</field>
    <value name="DECL0">
      <block type="dictionaries_create_with">
        <mutation items="2" elseif="0" else="0"></mutation>
        <value name="ADD0">
          <block type="pair">
            <value name="KEY">
              <block type="text">
                <field name="TEXT">Animal</field>
              </block>
            </value>
            <value name="VALUE">
              <block type="text">
                <field name="TEXT">Tiger</field>
              </block>
            </value>
          </block>
        </value>
        <value name="ADD1">
          <block type="pair">
            <value name="KEY">
              <block type="text">
                <field name="TEXT">Scientific Name</field>
              </block>
            </value>
            <value name="VALUE">
              <block type="text">
                <field name="TEXT">Panthera tigris</field>
              </block>
            </value>
          </block>
        </value>
      </block>
    </value>
    <statement name="STACK">
      <block type="controls_for_each_dict">
        <field name="KEY">key</field>
        <field name="VALUE">value</field>
        <value name="DICT">
          <block type="lexical_variable_get">
            <field name="VAR">animalDetail</field>
          </block>
        </value>
        <statement name="DO">
          <block type="controls_eval_but_ignore">
            <value name="VALUE">
              <block type="text_join">
                <mutation items="3" elseif="0" else="0"></mutation>
                <value name="ADD0">
                  <block type="lexical_variable_get">
                    <field name="VAR">key</field>
                  </block>
                </value>
                <value name="ADD1">
                  <block type="text">
                    <field name="TEXT"> : </field>
                  </block>
                </value>
                <value name="ADD2">
                  <block type="lexical_variable_get">
                    <field name="VAR">value</field>
                  </block>
                </value>
              </block>
            </value>
          </block>
        </statement>
      </block>
    </statement>
  </block>
```


## Functions

### Void function


```
  <block type="procedures_defnoreturn">
    <mutation items="0" elseif="0" else="0">
      <arg name="x"></arg>
      <arg name="y"></arg>
    </mutation>
    <field name="VAR0">x</field>
    <field name="VAR1">y</field>
    <field name="NAME">fooBar</field>
    <statement name="STACK">
      <block type="controls_eval_but_ignore">
        <value name="VALUE">
          <block type="math_add">
            <mutation items="2" elseif="0" else="0"></mutation>
            <value name="NUM0">
              <block type="lexical_variable_get">
                <field name="VAR">x</field>
              </block>
            </value>
            <value name="NUM1">
              <block type="lexical_variable_get">
                <field name="VAR">y</field>
              </block>
            </value>
          </block>
        </value>
      </block>
    </statement>
  </block>
```

### Result function

```
  <block type="procedures_defreturn">
    <mutation items="0" elseif="0" else="0">
      <arg name="n"></arg>
    </mutation>
    <field name="VAR0">n</field>
    <field name="NAME">double</field>
    <value name="RETURN">
      <block type="math_multiply">
        <mutation items="2" elseif="0" else="0"></mutation>
        <value name="NUM0">
          <block type="lexical_variable_get">
            <field name="VAR">n</field>
          </block>
        </value>
        <value name="NUM1">
          <block type="math_number">
            <field name="NUM">2</field>
          </block>
        </value>
      </block>
    </value>
  </block>
```

Or multiple expressions:

```
  <block type="procedures_defreturn">
    <mutation items="0" elseif="0" else="0">
      <arg name="n"></arg>
    </mutation>
    <field name="VAR0">n</field>
    <field name="NAME">FibSum</field>
    <value name="RETURN">
      <block type="controls_choose">
        <value name="TEST">
          <block type="math_compare">
            <field name="OP">LT</field>
            <value name="A">
              <block type="lexical_variable_get">
                <field name="VAR">n</field>
              </block>
            </value>
            <value name="B">
              <block type="math_number">
                <field name="NUM">2</field>
              </block>
            </value>
          </block>
        </value>
        <value name="THENRETURN">
          <block type="lexical_variable_get">
            <field name="VAR">n</field>
          </block>
        </value>
        <value name="ELSERETURN">
          <block type="math_add">
            <mutation items="2" elseif="0" else="0"></mutation>
            <value name="NUM0">
              <block type="procedures_callreturn">
                <mutation items="0" elseif="0" else="0" name="FibSum">
                  <arg name="n"></arg>
                </mutation>
                <field name="PROCNAME">FibSum</field>
                <value name="ARG0">
                  <block type="math_subtract">
                    <value name="A">
                      <block type="lexical_variable_get">
                        <field name="VAR">n</field>
                      </block>
                    </value>
                    <value name="B">
                      <block type="math_number">
                        <field name="NUM">1</field>
                      </block>
                    </value>
                  </block>
                </value>
              </block>
            </value>
            <value name="NUM1">
              <block type="procedures_callreturn">
                <mutation items="0" elseif="0" else="0" name="FibSum">
                  <arg name="n"></arg>
                </mutation>
                <field name="PROCNAME">FibSum</field>
                <value name="ARG0">
                  <block type="math_subtract">
                    <value name="A">
                      <block type="lexical_variable_get">
                        <field name="VAR">n</field>
                      </block>
                    </value>
                    <value name="B">
                      <block type="math_number">
                        <field name="NUM">2</field>
                      </block>
                    </value>
                  </block>
                </value>
              </block>
            </value>
          </block>
        </value>
      </block>
    </value>
  </block>
```
Note that there is no `return` statement in Falcon. The last statement in a body is taken as the output of an expression.

## List access

```
  <block type="local_declaration_statement">
    <mutation items="0" elseif="0" else="0">
      <localname name="numbers"></localname>
    </mutation>
    <field name="VAR0">numbers</field>
    <value name="DECL0">
      <block type="lists_create_with">
        <mutation items="3" elseif="0" else="0"></mutation>
        <value name="ADD0">
          <block type="math_number">
            <field name="NUM">1</field>
          </block>
        </value>
        <value name="ADD1">
          <block type="math_number">
            <field name="NUM">2</field>
          </block>
        </value>
        <value name="ADD2">
          <block type="math_number">
            <field name="NUM">4</field>
          </block>
        </value>
      </block>
    </value>
    <statement name="STACK">
      <block type="controls_eval_but_ignore">
        <value name="VALUE">
          <block type="lists_select_item">
            <value name="LIST">
              <block type="lexical_variable_get">
                <field name="VAR">numbers</field>
              </block>
            </value>
            <value name="NUM">
              <block type="math_number">
                <field name="NUM">2</field>
              </block>
            </value>
          </block>
        </value>
        <next>
          <block type="lists_replace_item">
            <value name="LIST">
              <block type="lexical_variable_get">
                <field name="VAR">numbers</field>
              </block>
            </value>
            <value name="NUM">
              <block type="math_number">
                <field name="NUM">1</field>
              </block>
            </value>
            <value name="ITEM">
              <block type="math_number">
                <field name="NUM">8</field>
              </block>
            </value>
          </block>
        </next>
      </block>
    </statement>
  </block>
```

## Dictionary access

```
  <block type="local_declaration_statement">
    <mutation items="0" elseif="0" else="0">
      <localname name="animalInfo"></localname>
    </mutation>
    <field name="VAR0">animalInfo</field>
    <value name="DECL0">
      <block type="dictionaries_create_with">
        <mutation items="2" elseif="0" else="0"></mutation>
        <value name="ADD0">
          <block type="pair">
            <value name="KEY">
              <block type="text">
                <field name="TEXT">Animal</field>
              </block>
            </value>
            <value name="VALUE">
              <block type="text">
                <field name="TEXT">Tiger</field>
              </block>
            </value>
          </block>
        </value>
        <value name="ADD1">
          <block type="pair">
            <value name="KEY">
              <block type="text">
                <field name="TEXT">Scientific Name</field>
              </block>
            </value>
            <value name="VALUE">
              <block type="text">
                <field name="TEXT">Panthera tigris</field>
              </block>
            </value>
          </block>
        </value>
      </block>
    </value>
    <statement name="STACK">
      <block type="controls_eval_but_ignore">
        <value name="VALUE">
          <block type="dictionaries_lookup">
            <value name="DICT">
              <block type="lexical_variable_get">
                <field name="VAR">animalInfo</field>
              </block>
            </value>
            <value name="KEY">
              <block type="text">
                <field name="TEXT">Scientific Name</field>
              </block>
            </value>
            <value name="NOTFOUND">
              <block type="text">
                <field name="TEXT">Not found</field>
              </block>
            </value>
          </block>
        </value>
      </block>
    </statement>
  </block>
```

## List lambdas

Inspired by Kotlin, list lambdas allow for list manipulation.

### Map lambda

Maps each element of a list to a new value.

```
  <block type="local_declaration_statement">
    <mutation items="0" elseif="0" else="0">
      <localname name="numbers"></localname>
    </mutation>
    <field name="VAR0">numbers</field>
    <value name="DECL0">
      <block type="lists_create_with">
        <mutation items="3" elseif="0" else="0"></mutation>
        <value name="ADD0">
          <block type="math_number">
            <field name="NUM">1</field>
          </block>
        </value>
        <value name="ADD1">
          <block type="math_number">
            <field name="NUM">2</field>
          </block>
        </value>
        <value name="ADD2">
          <block type="math_number">
            <field name="NUM">3</field>
          </block>
        </value>
      </block>
    </value>
    <statement name="STACK">
      <block type="local_declaration_statement">
        <mutation items="0" elseif="0" else="0">
          <localname name="doubled"></localname>
        </mutation>
        <field name="VAR0">doubled</field>
        <value name="DECL0">
          <block type="lists_map">
            <field name="VAR">n</field>
            <value name="LIST">
              <block type="lexical_variable_get">
                <field name="VAR">numbers</field>
              </block>
            </value>
            <value name="TO">
              <block type="math_multiply">
                <mutation items="2" elseif="0" else="0"></mutation>
                <value name="NUM0">
                  <block type="lexical_variable_get">
                    <field name="VAR">n</field>
                  </block>
                </value>
                <value name="NUM1">
                  <block type="math_number">
                    <field name="NUM">2</field>
                  </block>
                </value>
              </block>
            </value>
          </block>
        </value>
        <statement name="STACK">
          <block type="controls_eval_but_ignore">
            <value name="VALUE">
              <block type="lexical_variable_get">
                <field name="VAR">doubled</field>
              </block>
            </value>
          </block>
        </statement>
      </block>
    </statement>
  </block>
```

### Filter lambda

Filters out unwanted elements.

```
  <block type="local_declaration_statement">
    <mutation items="0" elseif="0" else="0">
      <localname name="numbers"></localname>
    </mutation>
    <field name="VAR0">numbers</field>
    <value name="DECL0">
      <block type="lists_create_with">
        <mutation items="4" elseif="0" else="0"></mutation>
        <value name="ADD0">
          <block type="math_number">
            <field name="NUM">1</field>
          </block>
        </value>
        <value name="ADD1">
          <block type="math_number">
            <field name="NUM">2</field>
          </block>
        </value>
        <value name="ADD2">
          <block type="math_number">
            <field name="NUM">3</field>
          </block>
        </value>
        <value name="ADD3">
          <block type="math_number">
            <field name="NUM">4</field>
          </block>
        </value>
      </block>
    </value>
    <statement name="STACK">
      <block type="local_declaration_statement">
        <mutation items="0" elseif="0" else="0">
          <localname name="evens"></localname>
        </mutation>
        <field name="VAR0">evens</field>
        <value name="DECL0">
          <block type="lists_filter">
            <field name="VAR">n</field>
            <value name="LIST">
              <block type="lexical_variable_get">
                <field name="VAR">numbers</field>
              </block>
            </value>
            <value name="TEST">
              <block type="logic_compare">
                <field name="OP">EQ</field>
                <value name="A">
                  <block type="math_divide">
                    <field name="OP">REMAINDER</field>
                    <value name="DIVIDEND">
                      <block type="lexical_variable_get">
                        <field name="VAR">n</field>
                      </block>
                    </value>
                    <value name="DIVISOR">
                      <block type="math_number">
                        <field name="NUM">2</field>
                      </block>
                    </value>
                  </block>
                </value>
                <value name="B">
                  <block type="math_number">
                    <field name="NUM">0</field>
                  </block>
                </value>
              </block>
            </value>
          </block>
        </value>
        <statement name="STACK">
          <block type="controls_eval_but_ignore">
            <value name="VALUE">
              <block type="lexical_variable_get">
                <field name="VAR">evens</field>
              </block>
            </value>
          </block>
        </statement>
      </block>
    </statement>
  </block>
```

### Sort lambda

Helps to define a custom sort method.

```
  <block type="local_declaration_statement">
    <mutation items="0" elseif="0" else="0">
      <localname name="names"></localname>
    </mutation>
    <field name="VAR0">names</field>
    <value name="DECL0">
      <block type="lists_create_with">
        <mutation items="3" elseif="0" else="0"></mutation>
        <value name="ADD0">
          <block type="text">
            <field name="TEXT">Bob</field>
          </block>
        </value>
        <value name="ADD1">
          <block type="text">
            <field name="TEXT">Alice</field>
          </block>
        </value>
        <value name="ADD2">
          <block type="text">
            <field name="TEXT">John</field>
          </block>
        </value>
      </block>
    </value>
    <statement name="STACK">
      <block type="local_declaration_statement">
        <mutation items="0" elseif="0" else="0">
          <localname name="namesSorted"></localname>
        </mutation>
        <field name="VAR0">namesSorted</field>
        <value name="DECL0">
          <block type="lists_sort_comparator">
            <field name="VAR1">m</field>
            <field name="VAR2">n</field>
            <value name="LIST">
              <block type="lexical_variable_get">
                <field name="VAR">names</field>
              </block>
            </value>
            <value name="COMPARE">
              <block type="math_compare">
                <field name="OP">GT</field>
                <value name="A">
                  <block type="text_length">
                    <value name="VALUE">
                      <block type="lexical_variable_get">
                        <field name="VAR">m</field>
                      </block>
                    </value>
                  </block>
                </value>
                <value name="B">
                  <block type="text_length">
                    <value name="VALUE">
                      <block type="lexical_variable_get">
                        <field name="VAR">m</field>
                      </block>
                    </value>
                  </block>
                </value>
              </block>
            </value>
          </block>
        </value>
        <statement name="STACK">
          <block type="controls_eval_but_ignore">
            <value name="VALUE">
              <block type="lexical_variable_get">
                <field name="VAR">namesSorted</field>
              </block>
            </value>
          </block>
        </statement>
      </block>
    </statement>
  </block>
```

### Min and Max lambdas

Sorts the elements in a list and returns the maximum or minimum value.

```
  <block type="local_declaration_statement">
    <mutation items="0" elseif="0" else="0">
      <localname name="names"></localname>
    </mutation>
    <field name="VAR0">names</field>
    <value name="DECL0">
      <block type="lists_create_with">
        <mutation items="3" elseif="0" else="0"></mutation>
        <value name="ADD0">
          <block type="text">
            <field name="TEXT">Bob</field>
          </block>
        </value>
        <value name="ADD1">
          <block type="text">
            <field name="TEXT">Alice</field>
          </block>
        </value>
        <value name="ADD2">
          <block type="text">
            <field name="TEXT">John</field>
          </block>
        </value>
      </block>
    </value>
    <statement name="STACK">
      <block type="local_declaration_statement">
        <mutation items="0" elseif="0" else="0">
          <localname name="longestName"></localname>
        </mutation>
        <field name="VAR0">longestName</field>
        <value name="DECL0">
          <block type="lists_maximum_value">
            <field name="VAR1">m</field>
            <field name="VAR2">n</field>
            <value name="LIST">
              <block type="lexical_variable_get">
                <field name="VAR">names</field>
              </block>
            </value>
            <value name="COMPARE">
              <block type="math_compare">
                <field name="OP">GT</field>
                <value name="A">
                  <block type="text_length">
                    <value name="VALUE">
                      <block type="lexical_variable_get">
                        <field name="VAR">n</field>
                      </block>
                    </value>
                  </block>
                </value>
                <value name="B">
                  <block type="text_length">
                    <value name="VALUE">
                      <block type="lexical_variable_get">
                        <field name="VAR">m</field>
                      </block>
                    </value>
                  </block>
                </value>
              </block>
            </value>
          </block>
        </value>
        <statement name="STACK">
          <block type="controls_eval_but_ignore">
            <value name="VALUE">
              <block type="lexical_variable_get">
                <field name="VAR">longestName</field>
              </block>
            </value>
          </block>
        </statement>
      </block>
    </statement>
  </block>
```

### Reduce lambda

Reduce lambda reduces many elements to a single element.

```
  <block type="local_declaration_statement">
    <mutation items="0" elseif="0" else="0">
      <localname name="numbers"></localname>
    </mutation>
    <field name="VAR0">numbers</field>
    <value name="DECL0">
      <block type="lists_create_with">
        <mutation items="7" elseif="0" else="0"></mutation>
        <value name="ADD0">
          <block type="math_number">
            <field name="NUM">1</field>
          </block>
        </value>
        <value name="ADD1">
          <block type="math_number">
            <field name="NUM">2</field>
          </block>
        </value>
        <value name="ADD2">
          <block type="math_number">
            <field name="NUM">3</field>
          </block>
        </value>
        <value name="ADD3">
          <block type="math_number">
            <field name="NUM">4</field>
          </block>
        </value>
        <value name="ADD4">
          <block type="math_number">
            <field name="NUM">5</field>
          </block>
        </value>
        <value name="ADD5">
          <block type="math_number">
            <field name="NUM">6</field>
          </block>
        </value>
        <value name="ADD6">
          <block type="math_number">
            <field name="NUM">7</field>
          </block>
        </value>
      </block>
    </value>
    <statement name="STACK">
      <block type="local_declaration_statement">
        <mutation items="0" elseif="0" else="0">
          <localname name="numbersSum"></localname>
        </mutation>
        <field name="VAR0">numbersSum</field>
        <value name="DECL0">
          <block type="lists_reduce">
            <field name="VAR1">x</field>
            <field name="VAR2">valueSoFar</field>
            <value name="LIST">
              <block type="lexical_variable_get">
                <field name="VAR">numbers</field>
              </block>
            </value>
            <value name="INITANSWER">
              <block type="math_number">
                <field name="NUM">0</field>
              </block>
            </value>
            <value name="COMBINE">
              <block type="math_add">
                <mutation items="2" elseif="0" else="0"></mutation>
                <value name="NUM0">
                  <block type="lexical_variable_get">
                    <field name="VAR">x</field>
                  </block>
                </value>
                <value name="NUM1">
                  <block type="lexical_variable_get">
                    <field name="VAR">valueSoFar</field>
                  </block>
                </value>
              </block>
            </value>
          </block>
        </value>
        <statement name="STACK">
          <block type="controls_eval_but_ignore">
            <value name="VALUE">
              <block type="lexical_variable_get">
                <field name="VAR">numbersSum</field>
              </block>
            </value>
          </block>
        </statement>
      </block>
    </statement>
  </block>
```

### Example

For example, let’s say Bob has a list of lemons sold per day for the last week and he’d like to calculate his revenue for lemon priced at $2 each.

The days he missed are marked as "N/A"

```
  <block type="global_declaration">
    <field name="NAME">LemonadeSold</field>
    <value name="VALUE">
      <block type="lists_create_with">
        <mutation items="7" elseif="0" else="0"></mutation>
        <value name="ADD0">
          <block type="math_number">
            <field name="NUM">9</field>
          </block>
        </value>
        <value name="ADD1">
          <block type="math_number">
            <field name="NUM">12</field>
          </block>
        </value>
        <value name="ADD2">
          <block type="text">
            <field name="TEXT">N/A</field>
          </block>
        </value>
        <value name="ADD3">
          <block type="math_number">
            <field name="NUM">15</field>
          </block>
        </value>
        <value name="ADD4">
          <block type="math_number">
            <field name="NUM">18</field>
          </block>
        </value>
        <value name="ADD5">
          <block type="text">
            <field name="TEXT">N/A</field>
          </block>
        </value>
        <value name="ADD6">
          <block type="math_number">
            <field name="NUM">8</field>
          </block>
        </value>
      </block>
    </value>
  </block>
  ``` 

Then we create a function that calculates the total revenue using list lambdas:

```
  <block type="procedures_defreturn">
    <mutation items="0" elseif="0" else="0"></mutation>
    <field name="NAME">GetTotalRevenue</field>
    <value name="RETURN">
      <block type="lists_reduce">
        <field name="VAR1">x</field>
        <field name="VAR2">soFar</field>
        <value name="LIST">
          <block type="lists_map">
            <field name="VAR">n</field>
            <value name="LIST">
              <block type="lists_filter">
                <field name="VAR">n</field>
                <value name="LIST">
                  <block type="lexical_variable_get">
                    <field name="VAR">global LemonadeSold</field>
                  </block>
                </value>
                <value name="TEST">
                  <block type="math_is_a_number">
                    <field name="OP">NUMBER</field>
                    <value name="NUM">
                      <block type="lexical_variable_get">
                        <field name="VAR">n</field>
                      </block>
                    </value>
                  </block>
                </value>
              </block>
            </value>
            <value name="TO">
              <block type="math_multiply">
                <mutation items="2" elseif="0" else="0"></mutation>
                <value name="NUM0">
                  <block type="lexical_variable_get">
                    <field name="VAR">n</field>
                  </block>
                </value>
                <value name="NUM1">
                  <block type="math_number">
                    <field name="NUM">2</field>
                  </block>
                </value>
              </block>
            </value>
          </block>
        </value>
        <value name="INITANSWER">
          <block type="math_number">
            <field name="NUM">0</field>
          </block>
        </value>
        <value name="COMBINE">
          <block type="math_add">
            <mutation items="2" elseif="0" else="0"></mutation>
            <value name="NUM0">
              <block type="lexical_variable_get">
                <field name="VAR">x</field>
              </block>
            </value>
            <value name="NUM1">
              <block type="lexical_variable_get">
                <field name="VAR">soFar</field>
              </block>
            </value>
          </block>
        </value>
      </block>
    </value>
  </block>
```
