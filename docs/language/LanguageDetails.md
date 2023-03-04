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

##### TODO set_follow replace_operations

```

## Typechecking

TODO
