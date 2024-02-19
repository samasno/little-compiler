package tokens

type Token struct {
	Type    string
	Literal string
}

const (
	IDENTIFIER = "IDENTIFIER"
	ILLEGAL    = "ILLEGAL"
	EOF        = ""

	SPACE     = ` `
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"
	NEWLINE   = "\n"
	TAB       = "\t"

	// OPERATORS
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	MULTIPLY = "*"
	DIVIDE   = "/"

	// COMPARATORS
	EQUALTO = "=="
	LTE     = "<="
	GTE     = ">="
	NOT     = "!"

	// KEYWORDS
	STRING   = "str"
	INT      = "int" // type declarations
	INTEGER  = "INTEGER"
	FUNCTION = "fn"
	LET      = "let"
	RETURN   = "return"
)

func IsMathOperator(t Token) bool {
	switch t.Type {
	case PLUS,
		MINUS,
		MULTIPLY,
		DIVIDE:
		return true
	default:
		return false
	}
}
