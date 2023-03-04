# VORE - **V**erb**O**se **R**egular **E**xpressions

VORE is just a regular expression engine but with more readable syntax and some extra features that I felt would be useful from processing text. With this project, I wanted to learn about language design and implementation with a strong enough constraint that would keep me from expanding the scope to an unmanageable level.

[![Tests](https://github.com/jmeaster30/vore/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/jmeaster30/vore/actions/workflows/go.yml) [![codecov](https://codecov.io/gh/jmeaster30/vore/branch/main/graph/badge.svg?token=8NFDH5ALID)](https://codecov.io/gh/jmeaster30/vore)

[![codecov](https://codecov.io/gh/jmeaster30/vore/branch/main/graphs/tree.svg?token=8NFDH5ALID)](https://codecov.io/gh/jmeaster30/vore)

## Documentation

### [Docs Home](docs/DocumentationHome.md)

### [Getting Started](docs/GettingStarted.md)

### [Examples](docs/examples/)

### [Regex Comparison](docs/language/RegexComparison.md)

## Project Structure

### root

The root of this repository is the CLI app which is all in the main.go file. It just adds a command line interface over `libvore`.

### libvore

The core library where the regular expression engine is implemented.

### libvorejs

The Javascript wrapper for `libvore` it compiles the Go code into WASM and uses webpack and npm to package everything together.

### libvore-syntax-highlighter

This is a really basic syntax highlighter extension for VSCode so you can look at nice colors while writing out Vore code.

## About VORE

This language is mostly a result of difficulties that I have had with regular expressions and I also took into account some things I have heard others have difficulty with. Some of the difficulties I have experienced with regular expressions are how its difficult to read the expressions and how difficult it is to remember the syntax when writing them. So, when coming up with the syntax for VORE I tried to make it feel like you can type out what you want like you are saying the rules and you'll have the proper regular expression. I took some syntax style inspiration from SQL as you will probably see.

Another thing I wanted from this language is to fully encompass original regular expressions to the point where I would be able to write a transpiler from VORE to regex and back. However, there are some features and semantics that I am allowing that would make transpiling any arbitrary VORE code into a valid regular expression difficult if not impossible. I don't know if I will actually write the transpiler but if I think it would be a fun quick project I probably will.

Here are some examples of the language...

>This example replaces all instances of "test - error" or "test - fail" with "test - success"
>
>```replace all "test - " = prefix "error" or "fail" with prefix "success"```

In the above example you can see the functionality of replacing text but it also is an example of using variables in regular expressions. Original regular expressions had the ability to use numeric references and named capture groups but I feel this syntax is significantly easier.

>Vore find statement:
>
>```find all line start (at least zero any) = myLine line end```
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
