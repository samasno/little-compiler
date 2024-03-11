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
	DUBQ      = `"`

	// OPERATORS
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	MULTIPLY = "*"
	DIVIDE   = "/"

	// COMPARATORS
	EQUALTO  = "=="
	LT       = "<"
	GT       = ">"
	LTE      = "<="
	GTE      = ">="
	NEQUALTO = "!="
	NOT      = "!"
	INC      = "++"
	DEC      = "--"

	// KEYWORDS
	STRING   = "STRING"
	INT      = "int" // type declarations
	INTEGER  = "INTEGER"
	FUNCTION = "fn"
	LET      = "let"
	RETURN   = "return"
	TRUE     = "true"
	FALSE    = "false"
	IF       = "if"
	ELSE     = "else"
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
