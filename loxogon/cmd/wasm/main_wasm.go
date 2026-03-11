package main

import (
	"loxogon/interpreter"
	"loxogon/lexer"
	"loxogon/parser"
	"strings"
	"syscall/js"
)

func main() {
	js.Global().Set("runLox", js.FuncOf(runInWasm))
	select {}
}

func run(code string, interp interpreter.Interpreter) (interpreter.LoxObject, int, error) {
	l := lexer.New(code)
	toks, err := l.ScanTokens()
	if err != nil {
		return interpreter.LoxObject{}, 1, err
	}

	parsed, err := parser.Parse(toks)
	if err != nil {
		return interpreter.LoxObject{}, 2, err
	}

	for _, stmt := range parsed {
		err = interp.EvaluateStmt(stmt)
		if err != nil {
			return interpreter.LoxObject{}, 3, err
		}
	}

	return interp.LastExpr, 0, nil
}

func runInWasm(this js.Value, args []js.Value) any {
	source := args[0].String()
	var sb strings.Builder
	i := interpreter.NewWithWriter(&sb)
	result, exitCode, err := run(source, i)
	// _, _, _ = run(source, i)
	retVal := make(map[string]any)
	retVal["stdout"] = sb.String()
	retVal["lastExpr"] = result.String()
	retVal["exitCode"] = exitCode
	if err != nil {
		retVal["error"] = err.Error()
	} else {
		retVal["error"] = nil
	}
	return retVal
}
