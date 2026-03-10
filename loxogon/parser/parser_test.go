package parser

import (
	"loxogon/ast"
	"loxogon/lexer"
	"reflect"
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantStr string
		wantErr error
	}{
		{
			name:    "expressions",
			source:  "(5 - (3 - 1)) + -1;",
			wantStr: "(+ (group (- 5 (group (- 3 1)))) (- 1))",
			wantErr: nil,
		},
		{
			name:    "expressions with decimals",
			source:  "(5.6 / (3 - 1.1)) + -1;",
			wantStr: "(+ (group (/ 5.6 (group (- 3 1.1)))) (- 1))",
			wantErr: nil,
		},
		{
			name:    "Uninitialized variable",
			source:  "var a;",
			wantStr: "(var a)",
			wantErr: nil,
		},
		{
			name:    "Initialized variable",
			source:  "var b=1;",
			wantStr: "(var b=1)",
			wantErr: nil,
		},
		{
			name:    "List of statements",
			source:  "var a; var b;",
			wantStr: "(var a) (var b)",
			wantErr: nil,
		},
		{
			name:    "Print statement",
			source:  "print true;",
			wantStr: "(print true)",
			wantErr: nil,
		},
		{
			name:    "Print expressions",
			source:  "var a=1; var b=3; print a > b;",
			wantStr: "(var a=1) (var b=3) (print (> a b))",
			wantErr: nil,
		},
		{
			name:    "Variable assignment",
			source:  "var a; a=1;",
			wantStr: "(var a) (a = 1)",
			wantErr: nil,
		}, {
			name:    "Block",
			source:  "{var a; var b; print 5;}",
			wantStr: "{(var a) (var b) (print 5)}",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.source)
			toks, err := l.ScanTokens()
			if err != nil {
				t.Errorf("improperly lexed source: %v", err)
			}

			gotStmts, gotErrs := Parse(toks)
			// TODO fix DeepEqual on errors
			if !reflect.DeepEqual(gotErrs, tt.wantErr) {
				t.Errorf("Parse() gotErrs = %v, want %v", gotErrs, tt.wantErr)
			}
			gotStr := ast.StmtsToString(gotStmts)
			if !(strings.Compare(gotStr, tt.wantStr) == 0) {
				t.Errorf("Parse() gotExpr = \n%v\n, want \n%v\n", gotStr, tt.wantStr)
			}
		})
	}
}
