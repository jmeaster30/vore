# VORE - **V**erb**O**se **R**egular **E**xpressions

VORE is just a regular expression engine but with more readable syntax and some extra features that I felt would be useful from processing text. With this project, I wanted to learn about language design and implementation with a strong enough constraint that would keep me from expanding the scope to an unmanageable level.

## Getting Started

This project uses Go so you will need the Go compiler installed and you can just run `go run .` from the root of the repository.

## About VORE

This language is mostly a result of difficulties that I have had with regular expressions and I also took into account some things I have heard others have difficulty with. Some of the difficulties I have experienced with regular expressions are how its difficult to read the expressions and how difficult it is to remember the syntax when writing them. So, when coming up with the syntax for VORE I tried to make it feel like you can type out what you want like you are saying the rules and you'll have the proper regular expression. I took some syntax style inspiration from SQL as you will probably see.

Another thing I wanted from this language is to fully encompass original regular expressions to the point where I would be able to write a transpiler from VORE to regex and back. However, there are some features and semantics that I am allowing that would make transpiling any arbitrary VORE code into a valid regular expression difficult if not impossible. I don't know if I will actually write the transpiler but if I think it would be a fun quick project I probably will.

Here are some examples of the language...

>This example replaces all instances of "test - error" or "test - fail" with "test - success"
>
>``` replace all "test - " = prefix "error" or "fail" with prefix "success"```

In the above example you can see the functionality of replacing text but it also is an example of using variables in regular expressions. Original regular expressions had the ability to use numeric references and named capture groups but I feel this syntax is significantly easier.

>Vore find statement: 
>
>```find all sol (at least zero any) = myLine eol```
>
>Equivalent regular expression
>
>```/^(?<myLine>.*)$/g```

You can also use the variables to find that sequence again in the match. This next example matches the string "aabb"

>Vore example:
>
>```find all "a" = varA varA "b" = varB varB```
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
```^``` (beginning)| ```file start``` (start of file)
```$``` (end)      | ```file end``` (end of file)
N/A          | ```line start``` (start of line)*
N/A          | ```line end``` (end of line)**
**Character Classes**
```[\s\S]``` (any) | ```any```
```.``` (dot) | N/A
```\s``` (whitespace) | ```whitespace```
```\d``` (digit) | ```digit```
```\w``` (word) | instead of this I added ```letter``` which is equivalent to ```[a-zA-Z]``` in regex
```\p{}``` (unicode) | N/A ***
```\D``` (not a digit) | ```not digit``` (not operator works for every character class)
```[ABC]``` (character set) | ```in 'A', 'B', "C"``` (in set)
```[^ABC]``` (negated set) | ```not in 'A', "B", 'C'``` (not in set)
```[A-Z]``` (range) | ```'A' to 'Z'``` (range)
**Escaped Characters** ****
```\n``` (escaped character) | ```"\n"```
```\@``` (escaped at sign) | ```'@'```
**Groups & References** |
```\1``` (numeric reference) | ```myVar``` (variables)
```\k<name>``` (named back reference) | ```name``` (variables)
```(?<name>ABC)``` (named capturing group) | ```"ABC" = name``` (assigning variables)
```(ABC)``` (capturing group) (capturing behavior) | ```'ABC' = myVar``` (assigning variables)
```(ABC)``` (capturing group) (subexpression behavior) | ```("ABC")``` (subexpression)
```(?:ABC)``` (non-capturing group) | ```("ABC")``` (subexpression)
**Subroutines** |
```(?P<name>[abc])``` (subroutine) | ```('a', 'b', 'c') := name```
```(?P>name)``` (call subroutine) | ```name```
**Recursion**
```a(?R)?b``` (Recurses on the entire regex) | ```("a" maybe mySub 'b') := mySub``` (Recursion only within the subroutine)
```(a(?1)?b)``` (Recurses only on that capture group) | ```("a" maybe mySub 'b') := mySub``` (same as before)
**Quantifiers & Alternation**
```a+``` (plus) | ```at least 1 'a'``` (at least)
```a*``` (start) | ```at least 0 "a"``` (at least)
```a{3}``` (quantifier) | ```exactly 3 'a'``` (exactly)
```a{4,}``` (quantifier) | ```at least 4 'a'``` (at least)
```a{5,8}``` (quantifier) | ```between 5 and 8 'a'``` (between)
```a{0,4}``` (quantifier) | ```at most 4 'a'``` (at most)
```a?``` (optional) | ```maybe 'a'``` (maybe)
```a+?``` (lazy) | ```at least 1 'a' fewest``` (at least followed by fewest)
```a\|b``` (alternation) | ```'a' or "b"``` (or)

\* also matches start of file for the first line

\** also matches end if file for the last line

\*** I know I want to add this sometime in the future but I am just trying to get this basic language out. Also, many of the stuff marked "N/A" in VORE is something that could potentially be added.

\**** all characters in quotes are "escaped" so far I added basic C-style escape characters but will probably expand on it more when I add Unicode characters.

## Notes on the RegEx modifiers
I do not like the idea of these modifiers in regular expressions. They are just kinda weird and personally I rarely used them (besides global and multiline), so I figured I wouldn't add them or add language contructs to replace them. Here is a table of the regex modifiers and their equivalents in VORE

Modifier | Name | Behavior | VORE Equivalent
---------|------|----|--------------
```g```  | Global | Retains the end index of the last match allowing subsequent searches to find further matches in the text | Base assumption of search functionality
```u``` | Unicode | Allows using extended unicode escape characters | Unicode will be a base feature of VORE (once I actually get it working lol)
```i``` | Ignorecase | Basically what it sounds like, the case of the characters in the search are ignored when making a match | The grave quotes around a string would perform a match ignoring the case of the string.
```m``` | Multiline | Changes the behavior of ```^``` and ```$``` to match the start of a line and end of a line instead of the start of the file and end of the file. | Instead of having a modifier, I split the behavior of ```^``` and ```$``` into ```sof```, ```sol```, ```eof```, and ```eol```.
```s``` | Dotall | Changes the behavior of ```.``` to match all characters including newlines instead of the original behavior of matching all characters except newline | Currently, I am assuming this behavior with the ```any``` character class. I do like the original behavior of ```.``` but I feel it would be confusing if something called ```any``` didn't match ALL characters. (Maybe we can add ```any*``` to replace the original ```.``` since it's "anything* (*but newlines)" lol)
```y``` | Sticky | Does not advance the last match index unless a match was found at that index | I have never used this in my life and even seeing what it is used for I couldn't think of a use for this
