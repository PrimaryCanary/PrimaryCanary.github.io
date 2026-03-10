package main

import (
	"bufio"
	"fmt"
	"loxogon/ast"
	"loxogon/interpreter"
	"loxogon/lexer"
	"loxogon/parser"
	"os"
	"strings"
)

func main() {
	interp := interpreter.New()
	if len(os.Args) > 2 {
		fmt.Println("Usage: loxogon [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		_, exit, err := runFile(os.Args[1], interp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "running script failed: %s\n", err)
			os.Exit(exit)
		} else {
			os.Exit(0)
		}
	} else {
		exit, err := runRepl(interp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "running repl failed: %s\n", err)
			os.Exit(exit)
		} else {
			os.Exit(0)
		}
	}
}

func runFile(file string, interp interpreter.Interpreter) (ast.LoxObject, int, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return ast.LoxObject{}, 1, fmt.Errorf("could not read file: %w", err)
	}
	return run(string(bytes), interp)
}

func run(code string, interp interpreter.Interpreter) (ast.LoxObject, int, error) {
	l := lexer.New(code)
	toks, err := l.ScanTokens()
	if err != nil {
		return ast.LoxObject{}, 1, err
	}

	parsed, err := parser.Parse(toks)
	if err != nil {
		return ast.LoxObject{}, 2, err
	}

	for _, stmt := range parsed {
		err = interp.EvaluateStmt(stmt)
		if err != nil {
			return ast.LoxObject{}, 3, err
		}
	}

	return interp.LastExpr, 0, nil
}

func runRepl(interp interpreter.Interpreter) (int, error) {
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
		result, _, err := run(string(trimmed), interp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		} else {
			fmt.Println(result)
		}
	}
}
