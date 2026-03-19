# Loxogon Playground
For my application, I wrote a programming language called Loxogon and a website that executes Loxogon code then displays the results.

# Program Structure
Loxogon's implementation lives in the [loxogon/](loxogon/) folder. It is written in Go and broken up into three phases:

1) The [lexer](loxogon/lexer/) turns source code text into a stream of data that be efficiently operated on.
2) The [parser](loxogon/parser/) turns the token stream into an abstract syntax tree.
3) The [interpreter](loxogon/interpreter/) executes the syntax tree directly.

Loxogon is generally intended to be run locally as a CLI tool. Scripts can be run with `loxogon script1.lox script2.lox`. Loxogon also ships with a local read-eval-print loop (REPL) where expressions and code can be explored. The REPL is available by running `loxogon` with no arguments.

The online playground mimics the features of the local REPL, but provides a much nicer UX. It lives in the [playground/](playground/) folder. Most notably, it does not contain a reimplementation of Loxogon. Instead, Loxogon is compiled to WebAssembly and driven from Javascript. Make sure your browser has WebAssembly enabled.

A Makefile is provided to help build the code locally. Despite not written to be general, it should be clear what commands to run.