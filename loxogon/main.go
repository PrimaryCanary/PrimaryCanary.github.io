package main

import (
	"bufio"
	"fmt"
	"loxogon/interpreter"
	"loxogon/lexer"
	"loxogon/parser"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: loxogon [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		_, exit, err := runFile(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "running script failed: %s\n", err)
			os.Exit(exit)
		} else {
			os.Exit(0)
		}
	} else {
		exit, err := runRepl()
		if err != nil {
			fmt.Fprintf(os.Stderr, "running repl failed: %s\n", err)
			os.Exit(exit)
		} else {
			os.Exit(0)
		}
	}
}

func runFile(file string) (interpreter.LoxObject, int, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return interpreter.LoxObject{}, 1, fmt.Errorf("could not read file: %w", err)
	}
	return run(string(bytes))
}

func run(code string) (interpreter.LoxObject, int, error) {
	l := lexer.New(code)
	toks, err := l.ScanTokens()
	if err != nil {
		return interpreter.LoxObject{}, 1, err
	}

	ast, err := parser.Parse(toks)
	if err != nil {
		return interpreter.LoxObject{}, 2, err
	}

	result, err := interpreter.Evaluate(ast)
	if err != nil {
		return interpreter.LoxObject{}, 3, err
	}

	return result, 0, nil
}

func runRepl() (int, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, _, err := reader.ReadLine()
		if err != nil {
			return 1, fmt.Errorf("could not read line: %w", err)
		}
		trimmed := strings.TrimSpace(string(line))
		if len(trimmed) == 0 {
			return 0, nil
		}
		result, _, err := run(string(trimmed))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		fmt.Println(result)
	}
}
