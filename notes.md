# Introduction
## Challenges
1) There are at least six DSLs used in publishing Crafting Interpreters. What are they?
Make, bash for wrapper scripts, IML (not sure what that is, something for IntelliJ?), yaml, some Dart that looks like parsing, Sass and Jinja for generating the site, not sure if HTML and CSS count
2) Get a Java environment running (or Go for me)
Committed.
3) Do the same thing for C. To get some practice with pointers, define a doubly linked list of heap-allocated strings. Write functions to insert, find, and delete items from it. Test them.
This sucked ass. Committed

# A Map of the Territory
This talked about the various types of interpreters and compilers (tree-walk, bytecode, JIT). All at a high level. Not much interesting here except that one sidenote about a dude dying in a biker bar.
## Challenges
1) Download the source of a popular programming language and find it's lexer and parser.
Go: It's handwritten lexers and parsers. The tokens, as of 1647896aa227d8546de3dbe70a5049eecee964e3, live in `src/go/token/token.go`. It does this really funny thing where it uses `iota` to mark where certain token types begin and end then compares against those integers. It's lexer is in `src/scanner/scanner.go`. It's parser is in `src/parser/parser.go`.
2) JIT compilation tends to be the fastest way to interpret languages. Why doesn't every language use one?
Because they're wildly complex, take a million engineering-hours, and you get security issues when you screw it up. The fundamental idea is also taking some flak these days with research that says they never reach a steady state.
3) Most Lisp languages that compile to C also have a Lisp interpreter included. Why?
Lisp has a bunch of runtime code generation with its macros that need an on-demand interpreter.

# The Lox Language
It has all the normal stuff: a print statement, numbers (all numbers are floats, which didn't used to be normal but hey, Javascript), arithmetic operators. No implicit conversion might not be normal but it is great. Ayy closures. 

I learned from a sidenote that Smalltalk doesn't have built in conditionals. It relies on dynamic dispatch. That's crazy. And sounds rediculously slow.

I'm looking forward to this classes vs prototypes description. I've read it. I'm not sure it's helpful. It seems like the distinction is where the functions are located??? In class inheritance, the instances only have data and a reference to the class (which I assume looks basically like a vtable). My understanding is that prototypal inheritance, objects (synonym for instances?) have data and functions directly. It seems ad hoc. Idk. Doesn't seem that important.

## Challenges
TODO do these. I don't want to right now. Seems boring.
TODO do these. I don't want to right now. Seems boring.