package token

type TokenType string

const (
	ILLEGAL = "KHÔNG_HỢP_LỆ"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "ĐỊNH_DANH" // add, foobar, x, y, ...
	INT    = "SỐ_NGUYÊN" // 1343456
	STRING = "CHUỖI"     // "foobar"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "HÀM"
	LET      = "ĐẶT"
	TRUE     = "ĐÚNG"
	FALSE    = "SAI"
	IF       = "NẾU"
	ELSE     = "NGƯỢC_LẠI"
	RETURN   = "TRẢ_VỀ"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"hàm":       FUNCTION,
	"đặt":       LET,
	"đúng":      TRUE,
	"sai":       FALSE,
	"nếu":       IF,
	"ngược_lại": ELSE,
	"trả_về":    RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}