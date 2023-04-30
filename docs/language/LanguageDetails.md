# Language Details

## Tokens and Lexing

The lexer is a hand-rolled lexer that works very similar to a lexer you would find in any other language. An interesting part of this lexer is because the majority of our tokens are words it is really just gathering a list of word characters and comparing them in a big switch to get a token type.

As of writing this (3/4/2023), we don't have many aliases but I really want to start adding aliases as I feel the language would benefit from allowing people to just write code with the intent in mind instead of battling with syntax. So, in the future the raw syntax would be a lot more free.

### Whitespace

We have 2 tokens for whitespace. One is WHITESPACE this is meant to be for the character class that matches whitespace characters in the vore language. The other is WS which is meant to be the whitespace characters in the source code.

Whitespace characters in the source code is not necessary anywhere but it's use is to separate words so you can't type `findallwhitespace` since the lexer will assume you want the literal word "findallwhitespace" instead of `find all whitespace`.

### Lexing Rules in Regex and Vore

| Token Type | Regex | Vore |
|------------|-------|------|
| WS | `\s` | `whitespace` |
| COMMENT (single line) | `\-\-.*` | `'--' at least 0 any fewest line end` |
| COMMENT (block) | `\-\-\([\s\S]*?\)\-\-` | `'--(' at least 0 any fewest ')--'` |
| IDENTIFIER | `[a-zA-Z][a-zA-Z0-9]*` | `letter at least 0 (letter or digit)` |
| NUMBER | | |
| STRING | `('\|")[\s\S]*?\1` | `("'" or '"') = quote at least 0 any fewest quote` |
| EQUAL | `=` | `'='` |
| COLONEQ | `:=` | `':='` |
| COMMA | `,` | `','` |
| OPENPAREN | `\(` | `'='` |
| CLOSEPAREN | `\)` | `'('` |
| OPENCURLY | `{` | `'{'` |
| CLOSECURLY | `}` | `'}'` |
| PLUS | `\+` | `'+'` |
| MINUS | `-` | `'-'` |
| MULT | `\*` | `'*'` |
| DIV | `\/` | `'/'` |
| LESS | `<` | `'<'` |
| GREATER | `>` | `'>'` |
| LESSEQ | `<=` | `'<='` |
| GREATEREQ | `>=` | `'>='` |
| DEQUAL | `==` | `'=='` |
| NEQUAL | `!=` | `'!='` |
| MOD | `%` | `'%'` |
| FIND | `find` | `'find'` |
| REPLACE | `replace` | `'replace'` |
| WITH | `with` | `'with'` |
| SET | `set` | `'set'` |
| TO | `to` | `'to'` |
| PATTERN | `pattern` | `'pattern'` |
| MATCHES | `matches` | `'matches'` |
| TRANSFORM | `transform\|function` | `'transform' or 'function'` |
| ALL | `all` | `'all'` |
| SKIP | `skip` | `'skip'` |
| TAKE | `take` | `'take'` |
| TOP | `top` | `'top'` |
| LAST | `last` | `'last'` |
| ANY | `any` | `'any'` |
| WHITESPACE | `whitespace` | `'whitespace'` |
| DIGIT | `digit` | `'digit'` |
| UPPER | `upper` | `'upper'` |
| LOWER | `lower` | `'lower'` |
| LETTER | `letter` | `'letter'` |
| WHOLE | `whole` | `'whole'` |
| LINE | `line` | `'line'` |
| FILE | `file` | `'file'` |
| WORD | `word` | `'word'` |
| START | `start` | `'start'` |
| END | `end` | `'end'` |
| BEGIN | `begin` | `'begin'` |
| NOT | `not` | `'not'` |
| AT | `at` | `'at'` |
| LEAST | `least` | `'least'` |
| MOST | `most` | `'most'` |
| BETWEEN | `between` | `'between'` |
| AND | `and` | `'and'` |
| EXACTLY | `exactly` | `'exactly'` |
| MAYBE | `maybe` | `'maybe'` |
| FEWEST | `fewest` | `'fewest'` |
| NAMED | `named` | `'named'` |
| IN | `in` | `'in'` |
| OR | `or` | `'or'` |
| IF | `if` | `'if'` |
| THEN | `then` | `'then'` |
| ELSE | `else` | `'else'` |
| DEBUG | `debug` | `'debug'` |
| RETURN | `return` | `'return'` |
| HEAD | `head` | `'head'` |
| TAIL | `tail` | `'tail'` |
| LOOP | `loop` | `'loop'` |
| BREAK | `break` | `'break'` |
| CONTINUE | `continue` | `'continue'` |
| TRUE | `true` | `'true'` |
| FALSE | `false` | `'false'` |

Going through this made me realize that some of these are unused. There are also plans for more features that may change this list but I will work on keeping it up-to-date.

## Grammar and Parsing

The parser used here is a hand-rolled parser. I was trying to do it fairly freeform so there was no formalism. The main part of the parser ended up being LL(1) I believe and then for expressions in the transforms and predicates are parsed with a really basic pratt parser.

In the following grammar, grammar rules are lowercase with underscores and tokens are uppercase like they are in the table above.

```text
program -> command program 
        |  EOF
        .

command -> FIND amount search_operations
        |  REPLACE amount search_operations WITH replace_operations
        |  SET IDENTIFIER TO set_follow
        .

amount -> all
       | skip NUMBER amount_follow
       | top NUMBER
       | take NUMBER
       | last NUMBER
       .

amount_follow -> take NUMBER 
              | 
              .

search_operations -> search_operation search_operations
                  |  search_operation
                  .

search_operation -> AT LEAST NUMBER search_operation
                 |  AT MOST NUMBER search_operation
                 |  BETWEEN NUMBER AND NUMBER search_operation
                 |  EXACTLY NUMBER search_operation
                 |  MAYBE search_operation
                 |  IN list
                 |  NOT follow_not
                 |  subroutine
                 |  primary
                 |  @/ regexp /
                 .

list -> list_item follow_list .

follow_list -> COMMA list
            |
            .

follow_not -> IN list
           |  character_class_anchor follow_primary
           |  STRING follow_primary
           .

subroutine -> OPENCURLY search_operations CLOSECURLY EQUAL IDENTIFIER .

primary -> literal follow_primary .

follow_primary -> EQUAL IDENTIFIER
               |  OR literal_not_literal follow_or
               |  
               .

follow_or -> OR literal_not_literal follow_or
          |  
          .

literal_not_literal -> literal
                    |  NOT character_class_anchor
                    |  NOT STRING
                    .

literal -> OPENPAREN search_operations CLOSEDPAREN
        |  character_class_anchor
        |  STRING
        |  IDENTIFIER
        .

character_class_anchor -> ANY
                       |  WHITESPACE
                       |  DIGIT
                       |  UPPER
                       |  LOWER
                       |  LETTER
                       |  LINE START
                       |  LINE END
                       |  FILE START
                       |  FILE END
                       |  WORD START
                       |  WORD END
                       |  WHOLE LINE
                       |  WHOLE FILE
                       |  WHOLE WORD
                       .

regexp -> <https://262.ecma-international.org/13.0/#sec-patterns>
        .

set_follow -> PATTERN set_pattern
           |  MATCHES command
           |  TRANSFORM set_function
           |  FUNCTION set_function
           . 

set_pattern -> search_operations
            |  search_operations BEGIN process_statements END
            .

set_function -> BEGIN process_statements END
             |  process_statements END
             .

replace_operations -> replace_operation replace_operations
                   |  
                   .

replace_operation -> STRING
                  |  IDENTIFIER
                  .

process_statements -> process_statement process_statements
                   |  
                   .

process_statement -> SET IDENTIFIER TO 
                  |  IF process_expression THEN process_statements END
                  |  IF process_expression THEN process_statements ELSE process_statements END
                  |  RETURN process_expression
                  |  DEBUG process_expression
                  |  LOOP process_statements END
                  |  BREAK
                  |  CONTINUE
                  .

process_expression -> <pratt parser>
                   .

```

## Typechecking

There are 3 types in Vore: strings, numbers, and booleans. The goal is you don't really need to think about types so Vore uses type inferrence and type coersion in order to make the type system nearly invisible.

### Type Coersion

Type coersion will coerse types to make a sensible result. As stated above the idea is that the resulting types makes the most sense.

In the following table, italicized types are coerced.

| LH Operand   | Operator | RH Operand     | Result |
|--------------|----------|----------------|--------|
| string       | +        | **_string_**   | string |
| string       | ==       | **_string_**   | bool   |
| string       | !=       | **_string_**   | bool   |
| string       | <        | **_string_**   | bool   |
| string       | >        | **_string_**   | bool   |
| string       | <=       | **_string_**   | bool   |
| string       | >=       | **_string_**   | bool   |
|              | head     | string         | string |
|              | tail     | string         | string |
|              | not      | bool           | bool   |
| bool         | and      | **_bool_**     | bool   |
| bool         | or       | **_bool_**     | bool   |
| bool         | ==       | **_bool_**     | bool   |
| bool         | !=       | **_bool_**     | bool   |
| bool         | <        | **_bool_**     | bool   |
| bool         | >        | **_bool_**     | bool   |
| bool         | <=       | **_bool_**     | bool   |
| bool         | >=       | **_bool_**     | bool   |
| number       | ==       | **_number_**   | bool   |
| number       | !=       | **_number_**   | bool   |
| number       | <        | **_number_**   | bool   |
| number       | >        | **_number_**   | bool   |
| number       | <=       | **_number_**   | bool   |
| number       | >=       | **_number_**   | bool   |
| number       | +        | **_number_**   | number |
| number       | -        | **_number_**   | number |
| number       | *        | **_number_**   | number |
| number       | /        | **_number_**   | number |
| number       | %        | **_number_**   | number |
| **_number_** | -        | number         | number |
| **_number_** | *        | number         | number |
| **_number_** | /        | number         | number |
| **_number_** | %        | number         | number |

If an operand/type combination is not shown in this then a semantic error occurs saying that the operator is not defined for the provided type.

The coersion logic that transforms the value from one type to another was also designed to be as sensible as possible. This is really similar to other programming language's coersion style and I chose to mimic this style since it makes a lot of sense to me.

|        | string | number | bool |
|--------|--------|--------|------|
| string |        | `strconv.Itoa(number)` | `if bool then return "true" else return "false"` |
| number | `strconv.Atoi(string) (on error returns 0)` | | `if bool then return 1 else return 0` |
| bool   | `len(string) != 0` | `number != 0` | |

### Statement Type Requirements

Some of the statements that require an expression have some constraints on what the result can be that are not coerced. I may decide to change this behavior in the future.

| Statement | Part              | Required Type    |
|-----------|-------------------|------------------|
| if        | condition         | bool             |
| return    | function          | string or number |
| return    | pattern predicate | bool             |

Every other use of the expressions will have their type inferred or coerced.


