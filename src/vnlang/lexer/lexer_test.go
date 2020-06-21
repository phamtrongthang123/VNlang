package lexer

import (
	"fmt"
	"strings"
	"testing"
	"vnlang/token"
)

// At line 9 of this code, I put it */ instead /* because /* will trigger comment
func TestNextToken(t *testing.T) {
	input := `đặt năm = 5;
đặt mười = 10;

đặt cộng = hàm(x, y) {
  x + y;
};

đặt result = cộng(năm, mười);
!-*/5;
5 < 10 > 5;

nếu (5 < 10) {
	trả_về đúng;
} ngược_lại {
	trả_về sai;
}

10 == 10;
10 != 9;
"foobar"
"foo bar"
[1, 2];
{"foo": "bar"}
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "đặt"},
		{token.IDENT, "năm"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "đặt"},
		{token.IDENT, "mười"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "đặt"},
		{token.IDENT, "cộng"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "hàm"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "đặt"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "cộng"},
		{token.LPAREN, "("},
		{token.IDENT, "năm"},
		{token.COMMA, ","},
		{token.IDENT, "mười"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "nếu"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "trả_về"},
		{token.TRUE, "đúng"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "ngược_lại"},
		{token.LBRACE, "{"},
		{token.RETURN, "trả_về"},
		{token.FALSE, "sai"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(strings.NewReader(input))

	var tok token.Token
	for i, tt := range tests {
		for {
			tok = l.NextToken()
			if tok.Type != token.NEWLINE {
				break
			}
		}
		fmt.Printf("*** %q *** \n", tok.Type)
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
