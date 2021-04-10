# RegQL

RegQL - **Reg**ular Expression **Q**uery **L**anguage

RegQL is just a regular expression engine but with more readable syntax and some extra features that I felt would be useful from processing text.

I like how natural it is to read SQL queries so I tried to design the syntax similarly.

Like regular expressions, you can find matches but I also added the ability to transform the text. (I think you may be able to do substitutions in regular expressions but I don't know how they work and I don't know if they actually transform the text)

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

## Syntax mapping from Regex to RegQL

Regex Syntax | RegQL Syntax
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

\*** I know I want to add this sometime in the future but I am just trying to get this basic language out. Also, many of the stuff marked "N/A" in RegQL is something that could potentially be added.

\**** all characters in quotes are "escaped" but you need to follow usual escaping rules you have in any programming language with strings

## Todo

- [ ] Conditionals
- [ ] Recursion (we have subroutines with subexpressions and variables but I don't think recursion would work out-of-the-box)
- [ ] Unicode support
