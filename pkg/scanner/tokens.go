package scanner

type TokenType string

const (
	// Single-character tokens.

	LEFTPAREN  = TokenType("LEFTPAREN")
	RIGHTPAREN = TokenType("RIGHTPAREN")
	LEFTBRACE  = TokenType("LEFTBRACE")
	RIGHTBRACE = TokenType("RIGHTBRACE")
	COMMA      = TokenType("COMMA")
	DOT        = TokenType("DOT")
	MINUS      = TokenType("MINUS")
	PLUS       = TokenType("PLUS")
	SEMICOLON  = TokenType("SEMICOLON")
	SLASH      = TokenType("SLASH")
	STAR       = TokenType("STAR")

	// One or two character tokens.

	BANG         = TokenType("BANG")
	BANGEQUAL    = TokenType("BANGEQUAL")
	EQUAL        = TokenType("EQUAL")
	EQUALEQUAL   = TokenType("EQUALEQUAL")
	GREATER      = TokenType("GREATER")
	GREATEREQUAL = TokenType("GREATEREQUAL")
	LESS         = TokenType("LESS")
	LESSEQUAL    = TokenType("LESSEQUAL")

	// Literals.

	IDENTIFIER = TokenType("IDENTIFIER")
	STRING     = TokenType("STRING")
	NUMBER     = TokenType("NUMBER")

	// Keywords.

	AND    = TokenType("AND")
	CLASS  = TokenType("CLASS")
	ELSE   = TokenType("ELSE")
	FALSE  = TokenType("FALSE")
	FUN    = TokenType("FUN")
	FOR    = TokenType("FOR")
	IF     = TokenType("IF")
	NIL    = TokenType("NIL")
	OR     = TokenType("OR")
	PRINT  = TokenType("PRINT")
	RETURN = TokenType("RETURN")
	SUPER  = TokenType("SUPER")
	THIS   = TokenType("THIS")
	TRUE   = TokenType("TRUE")
	VAR    = TokenType("VAR")
	WHILE  = TokenType("WHILE")

	EOF = TokenType("EOF")
)

type Token struct {
	kind    TokenType
	lexeme  string
	literal any
	line    int
}

func NewToken(kind TokenType, lexeme string, literal any, line int) Token {
	return Token{
		kind:    kind,
		lexeme:  lexeme,
		literal: literal,
		line:    line,
	}
}

//func (t Token) String() string {
//	return string(t.kind) + " " + t.lexeme + " " + fmt.Sprint(t.literal)
//}

func (t Token) Lexeme() string {
	return t.lexeme
}

func (t Token) Kind() TokenType {
	return t.kind
}

func (t Token) Literal() any {
	return t.literal
}

func (t Token) Line() int {
	return t.line
}
