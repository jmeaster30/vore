# Syntax mapping from Regex to VORE

Regex Syntax | VORE Syntax
-------------|-------------
**Anchors** |
```^``` (beginning)| ```file start``` (start of file)
```$``` (end)      | ```file end``` (end of file)
```\b``` (boundary)| ```word start``` (start of word)
```\b``` (boundary)| ```word end``` (end of word)
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
```(?P<name>[abc])``` (subroutine) | ```{in 'a', 'b', 'c'} = name```
```(?P>name)``` (call subroutine) | ```name```
**Recursion**
```a(?R)?b``` (Recurses on the entire regex) | ```{"a" maybe mySub 'b'} = mySub``` (Recursion only within the subroutine)
```(a(?1)?b)``` (Recurses only on that capture group) | ```{"a" maybe mySub 'b'} = mySub``` (same as before)
**Quantifiers & Alternation**
```a+``` (plus) | ```at least 1 'a'``` (at least)
```a*``` (star) | ```at least 0 "a"``` (at least)
```a{3}``` (quantifier) | ```exactly 3 'a'``` (exactly)
```a{4,}``` (quantifier) | ```at least 4 'a'``` (at least)
```a{5,8}``` (quantifier) | ```between 5 and 8 'a'``` (between)
```a{0,4}``` (quantifier) | ```at most 4 'a'``` (at most)
```a?``` (optional) | ```maybe 'a'``` (maybe)
```a+?``` (lazy) | ```at least 1 'a' fewest``` (at least followed by fewest)
```a\|b``` (alternation) | ```'a' or "b"``` (or)
**Lookaround**
```(?=ABC)``` (positive lookahead) | TODO ```followed by "ABC"```
```(?!ABC)``` (negative lookahead) | TODO ```not followed by "ABC"```
```(?<=ABC)``` (positive lookbehind) | TODO ```preceded by "ABC"```
```(?<!ABC)``` (negative lookbehind) | TODO ```not preceded by "ABC"```
**SPECIAL**
```(?#This is a comment)``` (comment) | ```-- This is a comment``` (comment)
```(?#This would work as a block comment)``` (block comment) | ```--(This is a block comment)--``` (comment)
```(?'one'a)?(?('one')b\|c)``` (conditional) | TODO
```(?>regex)``` (atomic group) | MAYBE WONT DO - I don't see the use for it. Maybe it helps with performance?
```(?\|regex)``` (branch reset group) | WONT DO - Not useful since we don't use unnamed capture groups
```\K``` (Keep text out) | WONT DO - Doesn't seem useful with proper look around support
**_Replacement_**
all of them :) | There are a lot of things that would be useful but how I did replacement makes it irrelevant because you have a mini-imperative language to do processing on the match.

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
