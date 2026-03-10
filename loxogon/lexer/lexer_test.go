package lexer

import (
	"fmt"
	"loxogon/ast"
	"testing"
)

func TestLexer(t *testing.T) {
	lexerTests := []struct {
		input          string
		expectedTokens []struct {
			expectedKind    ast.TokenKind
			expectedLiteral string
		}
	}{
		{
			input: `(( )){}`,
			expectedTokens: []struct {
				expectedKind    ast.TokenKind
				expectedLiteral string
			}{
				{ast.LEFT_PAREN, ""},
				{ast.LEFT_PAREN, ""},
				{ast.RIGHT_PAREN, ""},
				{ast.RIGHT_PAREN, ""},
				{ast.LEFT_BRACE, ""},
				{ast.RIGHT_BRACE, ""},
			},
		},
		{
			input: `!*+-/=<> <= ==`,
			expectedTokens: []struct {
				expectedKind    ast.TokenKind
				expectedLiteral string
			}{
				{ast.BANG, ""},
				{ast.STAR, ""},
				{ast.PLUS, ""},
				{ast.MINUS, ""},
				{ast.SLASH, ""},
				{ast.EQUAL, ""},
				{ast.LESS, ""},
				{ast.GREATER, ""},
				{ast.LESS_EQUAL, ""},
				{ast.EQUAL_EQUAL, ""},
			},
		},
		{
			input: `"string literal"`,
			expectedTokens: []struct {
				expectedKind    ast.TokenKind
				expectedLiteral string
			}{
				{ast.STRING, "string literal"},
			},
		},
		{
			input: `"multiline
  string literal"`,
			expectedTokens: []struct {
				expectedKind    ast.TokenKind
				expectedLiteral string
			}{
				{ast.STRING, "multiline\n  string literal"},
			},
		},
		{
			input: `01234 45.67`,
			expectedTokens: []struct {
				expectedKind    ast.TokenKind
				expectedLiteral string
			}{
				{ast.NUMBER, "1234"},
				{ast.NUMBER, "45.67"},
			},
		},
		{
			input: `0123.`,
			expectedTokens: []struct {
				expectedKind    ast.TokenKind
				expectedLiteral string
			}{
				{ast.NUMBER, "123"},
				{ast.DOT, ""},
			},
		},
		{
			input: `.9897`,
			expectedTokens: []struct {
				expectedKind    ast.TokenKind
				expectedLiteral string
			}{
				{ast.DOT, ""},
				{ast.NUMBER, "9897"},
			},
		},
		{
			input: `id _ident _978indent _asd_ident`,
			expectedTokens: []struct {
				expectedKind    ast.TokenKind
				expectedLiteral string
			}{
				{ast.IDENTIFIER, "id"},
				{ast.IDENTIFIER, "_ident"},
				{ast.IDENTIFIER, "_978indent"},
				{ast.IDENTIFIER, "_asd_ident"},
			},
		},
		{
			input: `and class else false fun for if nil or print return super this true var while`,
			expectedTokens: []struct {
				expectedKind    ast.TokenKind
				expectedLiteral string
			}{
				{ast.AND, ""},
				{ast.CLASS, ""},
				{ast.ELSE, ""},
				{ast.FALSE, ""},
				{ast.FUN, ""},
				{ast.FOR, ""},
				{ast.IF_TOK, ""},
				{ast.NIL, ""},
				{ast.OR, ""},
				{ast.PRINT_TOK, ""},
				{ast.RETURN, ""},
				{ast.SUPER, ""},
				{ast.THIS, ""},
				{ast.TRUE, ""},
				{ast.VAR_TOK, ""},
				{ast.WHILE_TOK, ""},
			},
		},
		{
			input: `nile an _class`,
			expectedTokens: []struct {
				expectedKind    ast.TokenKind
				expectedLiteral string
			}{
				{ast.IDENTIFIER, "nile"},
				{ast.IDENTIFIER, "an"},
				{ast.IDENTIFIER, "_class"},
			},
		},
	}

	for _, lt := range lexerTests {
		t.Run(lt.input, func(t *testing.T) {
			l := New(lt.input)
			tokens, err := l.ScanTokens()
			if err != nil {
				t.Fatalf("ScanTokens failed: %v", err)
			}

			if tokens[len(tokens)-1].Kind != ast.EOF {
				t.Fatalf("missing EOF at end of input. input=%v", lt.input)
			}
			tokens = tokens[:len(tokens)-1]
			if len(tokens) != len(lt.expectedTokens) {
				t.Logf("%v", tokens)
				t.Fatalf("wrong number of tokens. expected=%d, got=%d",
					len(lt.expectedTokens), len(tokens))
			}

			for i, tt := range lt.expectedTokens {
				tok := tokens[i]
				tokLit := fmt.Sprintf("%v", tok.Literal)

				if tok.Kind != tt.expectedKind {
					t.Fatalf("tests[%d]: token kind wrong. expected=%v, got=%v",
						i, tt.expectedKind, tok.Kind)
				}

				if tok.Literal != nil && (tokLit != tt.expectedLiteral) {
					t.Fatalf("tests[%d]: literal wrong. expected=%v, got=%v",
						i, tt.expectedLiteral, tokLit)
				}
			}
		})
	}
}

func TestLexerErrors(t *testing.T) {
	inputs := []string{"\"asdf", "asd@asdf", "a\\na"}
	for i, in := range inputs {
		t.Run(in, func(t *testing.T) {
			l := New(in)
			toks, err := l.ScanTokens()
			if err == nil {
				t.Fatalf("tests[%d]: lexing should fail. got=%v", i, toks)
			}
		})
	}
}
