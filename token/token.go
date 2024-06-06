package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 123456
	STRING = "STRING" // "foobar"
	FLOAT  = "FLOAT"  // 123.456

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	MODULUS  = "%"
	POW      = "**"

	LT = "<"
	GT = ">"
	LE = "<="
	GE = ">="

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
	FUNCTION = "FUNCTION"
	LET      = "LET"
	CONST    = "CONST"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	WHILE    = "WHILE"
	FOR      = "FOR"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
	IMPORT   = "IMPORT"

	EQ     = "=="
	NOT_EQ = "!="

	AND   = "&"
	OR    = "|"
	XOR   = "^"
	TILDE = "~"
	SHR   = ">>"
	SHL   = "<<"

	INCREMENT = "++"
	DECREMENT = "--"

	PLUS_ASSIGN     = "+="
	MINUS_ASSIGN    = "-="
	ASTERISK_ASSIGN = "*="
	SLASH_ASSIGN    = "/="
	MODULUS_ASSIGN  = "%="
	AND_ASSIGN      = "&="
	OR_ASSIGN       = "|="
	XOR_ASSIGN      = "^="
	SHL_ASSIGN      = "<<="
	SHR_ASSIGN      = ">>="

	AND_AND = "&&"
	OR_OR   = "||"

	DOT = "."
)

var keywords = map[string]TokenType{
	"fn":       FUNCTION,
	"let":      LET,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"const":    CONST,
	"while":    WHILE,
	"for":      FOR,
	"break":    BREAK,
	"continue": CONTINUE,
	"import":   IMPORT,
}

// LookupIdent checks if the given identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
