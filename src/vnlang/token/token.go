package token

import "vnlang/scanner"

type TokenType string

const (
	NULL    = "RỖNG"
	ILLEGAL = "KHÔNG_HỢP_LỆ"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "ĐỊNH_DANH" // add, foobar, x, y, ...
	INT    = "SỐ_NGUYÊN" // 1343456
	STRING = "XÂU"       // "foobar"
	FLOAT  = "SỐ_THỰC"   // "3.14"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	MOD      = "%"
	DOT      = "."

	AND = "&&"
	OR  = "||"

	LT = "<"
	GT = ">"

	LE = "<="
	GE = ">="

	EQ     = "=="
	NOT_EQ = "!="

	IN = "THUỘC"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	NEWLINE   = "\n"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "HÀM"
	LET      = "ĐẶT"
	MUT      = "BIẾN"
	CONST    = "HẰNG"
	TRUE     = "ĐÚNG"
	FALSE    = "SAI"
	IF       = "NẾU"
	ELSEIF   = "CÒN_NẾU"
	ELSE     = "NGƯỢC_LẠI"
	RETURN   = "TRẢ_VỀ"

	WHILE    = "KHI"
	BREAK    = "NGẮT"
	CONTINUE = "TIẾP"
)

type Token struct {
	Type    TokenType
	Literal string
	Pos     scanner.Position
}

var keywords = map[string]TokenType{
	"hàm":       FUNCTION,
	"đặt":       LET,
	"biến":      MUT,
	"hằng":      CONST,
	"đúng":      TRUE,
	"sai":       FALSE,
	"nếu":       IF,
	"còn_nếu":   ELSEIF,
	"ngược_lại": ELSE,
	"trả_về":    RETURN,
	"khi":       WHILE,
	"ngắt":      BREAK,
	"tiếp":      CONTINUE,
	"thuộc":     IN,
}

func GetNullToken() Token {
	return Token{
		Type:    NULL,
		Literal: "",
	}
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
