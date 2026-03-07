package parser

import (
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
			source:  "(5 - (3 - 1)) + -1",
			wantStr: "(+ (group (- 5 (group (- 3 1)))) (- 1))",
			wantErr: nil,
		},
		{
			name:    "expressions with decimals",
			source:  "(5.6 / (3 - 1.1)) + -1",
			wantStr: "(+ (group (/ 5.6 (group (- 3 1.1)))) (- 1))",
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

			gotExpr, gotErrs := Parse(toks)
			gotStr := gotExpr.String()
			if !(strings.Compare(gotStr, tt.wantStr) == 0) {
				t.Errorf("Parse() gotExpr = \n%v\n, want \n%v\n", gotStr, tt.wantStr)
			}
			if !reflect.DeepEqual(gotErrs, tt.wantErr) {
				t.Errorf("Parse() gotErrs = %v, want %v", gotErrs, tt.wantErr)
			}
		})
	}
}
