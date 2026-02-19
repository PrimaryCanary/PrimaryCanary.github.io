package lexer

import (
	"fmt"
	"loxogon/token"
	"testing"
)

func TestLexer(t *testing.T) {
	lexerTests := []struct {
		input          string
		expectedTokens []struct {
			expectedKind    token.TokenKind
			expectedLiteral string
		}
	}{
		{
			input: `(( )){}`,
			expectedTokens: []struct {
				expectedKind    token.TokenKind
				expectedLiteral string
			}{
				{token.LEFT_PAREN, ""},
				{token.LEFT_PAREN, ""},
				{token.RIGHT_PAREN, ""},
				{token.RIGHT_PAREN, ""},
				{token.LEFT_BRACE, ""},
				{token.RIGHT_BRACE, ""},
			},
		},
		{
			input: `!*+-/=<> <= ==`,
			expectedTokens: []struct {
				expectedKind    token.TokenKind
				expectedLiteral string
			}{
				{token.BANG, ""},
				{token.STAR, ""},
				{token.PLUS, ""},
				{token.MINUS, ""},
				{token.SLASH, ""},
				{token.EQUAL, ""},
				{token.LESS, ""},
				{token.GREATER, ""},
				{token.LESS_EQUAL, ""},
				{token.EQUAL_EQUAL, ""},
			},
		},
		{
			input: `"string literal"`,
			expectedTokens: []struct {
				expectedKind    token.TokenKind
				expectedLiteral string
			}{
				{token.STRING, "string literal"},
			},
		},
		{
			input: `"multiline
  string literal"`,
			expectedTokens: []struct {
				expectedKind    token.TokenKind
				expectedLiteral string
			}{
				{token.STRING, "multiline\n  string literal"},
			},
		},
		{
			input: `01234 45.67`,
			expectedTokens: []struct {
				expectedKind    token.TokenKind
				expectedLiteral string
			}{
				{token.NUMBER, "1234"},
				{token.NUMBER, "45.67"},
			},
		},
		{
			input: `0123.`,
			expectedTokens: []struct {
				expectedKind    token.TokenKind
				expectedLiteral string
			}{
				{token.NUMBER, "123"},
				{token.DOT, ""},
			},
		},
		{
			input: `.9897`,
			expectedTokens: []struct {
				expectedKind    token.TokenKind
				expectedLiteral string
			}{
				{token.DOT, ""},
				{token.NUMBER, "9897"},
			},
		},
		{
			input: `id _ident _978indent _asd_ident`,
			expectedTokens: []struct {
				expectedKind    token.TokenKind
				expectedLiteral string
			}{
				{token.IDENTIFIER, "id"},
				{token.IDENTIFIER, "_ident"},
				{token.IDENTIFIER, "_978indent"},
				{token.IDENTIFIER, "_asd_ident"},
			},
		},
		{
			input: `and class else false fun for if nil or print return super this true var while`,
			expectedTokens: []struct {
				expectedKind    token.TokenKind
				expectedLiteral string
			}{
				{token.AND, ""},
				{token.CLASS, ""},
				{token.ELSE, ""},
				{token.FALSE, ""},
				{token.FUN, ""},
				{token.FOR, ""},
				{token.IF, ""},
				{token.NIL, ""},
				{token.OR, ""},
				{token.PRINT, ""},
				{token.RETURN, ""},
				{token.SUPER, ""},
				{token.THIS, ""},
				{token.TRUE, ""},
				{token.VAR, ""},
				{token.WHILE, ""},
			},
		},
		{
			input: `nile an _class`,
			expectedTokens: []struct {
				expectedKind    token.TokenKind
				expectedLiteral string
			}{
				{token.IDENTIFIER, "nile"},
				{token.IDENTIFIER, "an"},
				{token.IDENTIFIER, "_class"},
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

			if tokens[len(tokens)-1].Kind != token.EOF {
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
