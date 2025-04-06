package main

import (
	"bufio"
	"fmt"
	"loxogon/bad"
	"loxogon/lexer"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: loxogon [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		exit, err := runFile(os.Args[1])
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

func runFile(file string) (int, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return 1, fmt.Errorf("could not read file: %w", err)
	}
	return run(string(bytes))
}

func run(code string) (int, error) {
	l := lexer.New(code)
	toks := l.ScanTokens()

	for _, t := range toks {
		fmt.Println(t)
	}
	return 0, nil
}

func runRepl() (int, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, _, err := reader.ReadLine()
		if err != nil {
			return 1, fmt.Errorf("could not read line: %w", err)
		}
		if len(line) == 0 {
			return 0, nil
		}
		_, err = run(string(line))
		if bad.Ok() {
			//fmt.Fprintf(os.Stderr, "error running code: %s\n", err)
			bad.Reset()
		}
	}
}
