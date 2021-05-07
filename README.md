# VORE

VORE - **V**erb**O**se Regula**R** Expr**E**ssions

VORE is just a regular expression engine but with more readable syntax and some extra features that I felt would be useful from processing text.

This language is mostly a result of difficulties that I have had with regular expressions and I also took into account some things I have heard others have difficulty with. Some of the difficulties I have experienced with regular expressions are how its difficult to read the expressions and how difficult it is to remember the syntax when writing them. So, when coming up with the syntax for VORE I tried to make it feel like you can type out what you want like you are saying the rules and you'll have the proper regular expression. I took some syntax style inspiration from SQL as you will probably see.

Another thing I wanted from this language is to fully encompass original regular expressions to the point where I would be able to write a transpiler from VORE to regex and back. I don't know if I will actually write the transpiler but if I think it would be a fun quick project I probably will.

Here are some examples of the language...

>This example replaces all instances of "test - error" or "test - fail" with "test - success"
>
>``` replace all "test - " = @a "error" or "fail" with @a "success"```

In the above example you can see the functionality of replacing text but it also is an example of using variables in regular expressions. Original regular expressions had the ability to use numeric references and named capture groups but I feel this syntax is significantly easier.

>RegQL find statement: 
>
>```find all sol (at least zero any) = @myLine eol```
>
>Equivalent regular expression
>
>```/^(?<myLine>.*)$/g```

You can also use the variables to find that sequence again in the match. This next example matches the string "aabb"

>RegQL example:
>
>```find all "a" = $varA $varA "b" = $varB $varB```
>
>Equivalent regular expression (using numeric references)
>
>```/(a)\1(b)\2/g```

While regular expressions are shorter in length, you don't need to understand the meaning of these mysterious symbols to be able to read what the query is meant to do.

When coming up with the syntax for this language I basically translated every symbol that you can have in regular expressions and combined some that I felt were represented by the other constructs.

## Syntax mapping from Regex to VORE

Regex Syntax | VORE Syntax
-------------|-------------
**Anchors** | 
```^``` (beginning)| ```sof``` (start of file)
```$``` (end)      | ```eof``` (end of file)
N/A          | ```sol``` (start of line)*
N/A          | ```eol``` (end of line)**
**Character Classes**
```[\s\S]``` (any) | ```any```
```.``` (dot) | N/A
```\s``` (whitespace) | ```whitespace```
```\d``` (digit) | ```digit```
```\w``` (word) | N/A
```\p{}``` (unicode) | N/A ***
```\D``` (not a digit) | ```not digit``` (not operator works for every character class)
```[ABC]``` (character set) | ```[A, B, C]``` (character set)
```[^ABC]``` (negated set) | ```not ['A', "B", 'C']``` (the not operator again)
```[A-Z]``` (range) | ```['A'-'Z']``` (range)
**Escaped Characters** ****
```\n``` (escaped character) | ```"\n"```
```\@``` (escaped at sign) | ```'@'```
**Groups & References** |
```\1``` (numeric reference) | ```@1``` (variables)
```\k<name>``` (named back reference) | ```@name``` (variables)
```(?<name>ABC)``` (named capturing group) | ```"ABC" = @name``` (assigning variables)
```(ABC)``` (capturing group) (capturing behavior) | ```'ABC' = @1``` (assigning variables)
```(ABC)``` (capturing group) (subexpression behavior) | ```("ABC")``` (subexpression)
```(?:ABC)``` (non-capturing group) | ```("ABC")``` (subexpression)
**Subroutines** |
```(?P<name>[abc])``` (subroutine) | ```[abc] = $name```
```(?P>name)``` (call subroutine) | ```$name```

**Quantifiers & Alternation**
```a+``` (plus) | ```at least 1 'a'``` (at least)
```a*``` (start) | ```at least 0 "a"``` (at least)
```a{3}``` (quantifier) | ```exactly 3 'a'``` (exactly)
```a{4,}``` (quantifier) | ```at least 4 'a'``` (at least)
```a{5,8}``` (quantifier) | ```between 5 and 8 'a'``` (between)
```a?``` (optional) | ```at most 1 'a'``` (at most)
```a+?``` (lazy) | ```at least 1 'a' fewest``` (at least followed by fewest)
```a\|b``` (alternation) | ```'a' or "b"``` (or)

\* also matches start of file for the first line

\** also matches end if file for the last line

\*** I know I want to add this sometime in the future but I am just trying to get this basic language out. Also, many of the stuff marked "N/A" in VORE is something that could potentially be added.

\**** all characters in quotes are "escaped" but you need to follow usual escaping rules you have in any programming language with strings

## Notes on the RegEx modifiers
There are no equivalents to the expression modifiers in this language. If you don't know, regular expressions have modifiers that follow the regular expression and modify what the regular expression does (Ex. ```/(abc)+/gm``` the ```g``` and the ```m``` are the modifiers). In my experience, I have never even seen some of these modifiers used and other times have found myself confused as to why my regular expression wasn't working properly when I was either missing a modifier or had an extraneous one. Despite writing this paragraph about my issues with these modifiers I actually don't have THAT much of an issue with them because I really only ever used the ```g``` modifier. I had never needed to explore the other ones in any depth at all. So, here is a table of all of the regex modifiers and their place in the VORE language...

Modifier | Name | Behavior | Place in VORE
---------|------|----|--------------
```g```  | Global | Retains the end index of the last match allowing subsequent searches to find further matches in the text | Base assumption of search functionality
```u``` | Unicode | Allows using extended unicode escape characters | Unicode will be a base feature of VORE (once I actually get it working lol)
```i``` | Ignorecase | Basically what it sounds like, the case of the characters in the search are ignored when making a match | This is the only modifier that makes me question my stance on this. I may add a statement that sets let the character case to be ignored
```m``` | Multiline | Changes the behavior of ```^``` and ```$``` to match the start of a line and end of a line instead of the start of the file and end of the file. | This has no place in VORE I would rather have the behavior of an anchor be explicit so I split the behavior of ```^``` and ```$``` into ```sof```, ```sol```, ```eof```, and ```eol```.
```s``` | Dotall | Changes the behavior of ```.``` to match all characters including newlines instead of the original behavior of matching all characters except newline | Currently, I am assuming this behavior with the ```any``` character class. I do like the original behavior of ```.``` but I feel it would be confusing if something called ```any``` didn't match ALL characters. (Maybe we can add ```any*``` to replace the original ```.``` since it's "anything* (*but newlines)" lol)
```y``` | Sticky | Does not advance the last match index unless a match was found at that index | No place in VORE. The use cases for this modifier can be implemented without this modifier with not much extra programming overhead. If someone can come up with a way to change my mind on this I am open but I doubt I will be easily convinced.


## Todo

- [ ] Some idea of "functions" which allow for transforming atoms based on computations
  - These "functions" will be statements and almost definitely be pure functions
- [ ] Unicode support
- [ ] Lookahead and lookbehind regex equivalent. I have ideas for the syntax but I am not set on it
