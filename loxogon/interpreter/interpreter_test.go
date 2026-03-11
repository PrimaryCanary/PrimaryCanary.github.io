package interpreter

import (
	"loxogon/lexer"
	"loxogon/parser"
	"reflect"
	"strings"
	"testing"
)

func TestParser_Evaluate(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		wantResult any
		wantStdout string
		wantErr    error
	}{
		{
			name:       "expressions",
			source:     "(5 - (3 - 1)) + -1;",
			wantResult: float64(2),
			wantStdout: "",
			wantErr:    nil,
		},
		{
			name:       "expressions with decimals",
			source:     "(5.6 / (3 - 1.1)) + -1;",
			wantResult: 1.9473684210526314,
			wantStdout: "",
			wantErr:    nil,
		},
		{
			name:       "Overloaded '+'",
			source:     "\"asdf\" + \"fdsa\";",
			wantResult: "asdffdsa",
			wantStdout: "",
			wantErr:    nil,
		},
		{
			name:       "Initialized variable",
			source:     "var b=1;",
			wantResult: nil,
			wantStdout: "",
			wantErr:    nil,
		},
		{
			name:       "Print statement",
			source:     "print true;",
			wantResult: nil,
			wantStdout: "true\n",
			wantErr:    nil,
		},
		{
			name:       "Print expressions",
			source:     "var a=1; var b=3; print a > b;",
			wantResult: nil,
			wantStdout: "false\n",
			wantErr:    nil,
		},
		{
			name:       "Variable assignment",
			source:     "var a; a=1;",
			wantResult: float64(1),
			wantStdout: "",
			wantErr:    nil,
		},
		{
			name: "Blocks with scopes",
			source: `var a = "global a";
					 var b = "global b";
					 var c = "global c";
					 {
					 	var a = "outer a";
					 	var b = "outer b";
					 	{
							var a = "inner a";
							print a;
							print b;
							print c;
						}
						print a;
						print b;
						print c;
					 }
					 print a;
					 print b;
					 print c;`,
			wantResult: nil,
			wantStdout: `inner a
outer b
global c
outer a
outer b
global c
global a
global b
global c
`,
			wantErr: nil,
		},
		{
			name:       "If statement",
			source:     "if ((1<2)) {var a=1; var b=2; print a+b;}",
			wantResult: nil,
			wantStdout: "3\n",
			wantErr:    nil,
		},
		{
			name:       "If else statement",
			source:     "if (!(1<2)) {var a=1; var b=2; print a+b;} else {var c=\"foo\"; var d=\"bar\"; print c+d;}",
			wantResult: nil,
			wantStdout: "foobar\n",
			wantErr:    nil,
		},
		{
			name:       "Logical or short circuit",
			source:     "2 or 3;",
			wantResult: float64(2),
			wantStdout: "",
			wantErr:    nil,
		},
		{
			name:       "Logical or long circuit",
			source:     "nil or 3;",
			wantResult: float64(3),
			wantStdout: "",
			wantErr:    nil,
		},
		{
			name:       "Logical and short circuit",
			source:     "nil and 3;",
			wantResult: nil,
			wantStdout: "",
			wantErr:    nil,
		},
		{
			name:       "Logical and long circuit",
			source:     "2 and 3;",
			wantResult: float64(3),
			wantStdout: "",
			wantErr:    nil,
		},
		{
			name:       "While loop",
			source:     "var a=0; while(a < 5) { print a; a = a+1; }",
			wantResult: float64(5),
			wantStdout: "0\n1\n2\n3\n4\n",
			wantErr:    nil,
		},
		{
			name:       "For loop only condition",
			source:     "for (; false;) print 5;",
			wantResult: nil,
			wantStdout: "",
			wantErr:    nil,
		},
		{
			name:       "For loop condition and incr",
			source:     "var a=0; for (; a < 3; a=a+1) print a;",
			wantResult: float64(3),
			wantStdout: "0\n1\n2\n",
			wantErr:    nil,
		},
		{
			name:       "For loop all three",
			source:     "for (var a=0; a < 3; a=a+1) print a;",
			wantResult: float64(3),
			wantStdout: "0\n1\n2\n",
			wantErr:    nil,
		},
		{
			name: "Fibonacci numbers",
			source: `var a = 0;
var temp;

for (var b = 1; a < 10000; b = temp + b) {
  print a;
  temp = a;
  a = b;
}`,
			wantResult: float64(17711),
			wantStdout: "0\n1\n1\n2\n3\n5\n8\n13\n21\n34\n55\n89\n" +
				"144\n233\n377\n610\n987\n1597\n2584\n4181\n6765\n",
			wantErr: nil,
		},
		{
			name:       "Recursive function declaration and call",
			source:     "fun count(n) { if (n > 1) {count(n-1);} print n; } count(3);",
			wantResult: nil,
			wantStdout: "1\n2\n3\n",
			wantErr:    nil,
		},
		{
			name: "sayHi function declaration",
			source: `
fun sayHi(first, last) {
  print "Hi, " + first + " " + last + "!";
}

sayHi("Dear", "Reader");`,
			wantResult: nil,
			wantStdout: "Hi, Dear Reader!\n",
			wantErr:    nil,
		},
		{
			name:       "Fibonnacci",
			source:     "fun fib(n) { if (n <= 1) return n; return fib(n-1) + fib(n-2); } fib(10);",
			wantResult: float64(55),
			wantStdout: "",
			wantErr:    nil,
		},
		{
			name:       "Recursive countdown",
			source:     "fun countdown(from, to) { if (from < to) return; print from; countdown(from - 1, to); } countdown(10, 5);",
			wantResult: nil,
			wantStdout: "10\n9\n8\n7\n6\n5\n",
			wantErr:    nil,
		},
	}
	// TODO add tests with errors
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			l := lexer.New(tt.source)
			toks, err := l.ScanTokens()
			if err != nil {
				t.Errorf("improperly lexed source: %v", err)
			}

			stmts, err := parser.Parse(toks)
			if err != nil {
				t.Errorf("improperly parsed source: %v", err)
			}

			var stdout strings.Builder
			var gotErrs error
			i := NewWithWriter(&stdout)
			for _, st := range stmts {
				gotErrs = i.EvaluateStmt(st)
				if !reflect.DeepEqual(gotErrs, tt.wantErr) {
					t.Errorf("%v gotErrs = %v, want %v", tt.name, gotErrs, tt.wantErr)
				}
			}
			if !reflect.DeepEqual(tt.wantResult, i.LastExpr.Value) {
				// if want == gotResult.Value {
				t.Errorf("%v gotResult = \n%v\n, want \n%v\n", tt.name, i.LastExpr.Value, tt.wantResult)
			}
			gotStdout := stdout.String()
			if strings.Compare(gotStdout, tt.wantStdout) != 0 {
				t.Errorf("%v gotStdout = \n%v\n, want \n%v\n", tt.name, gotStdout, tt.wantStdout)
			}
		})
	}
}
