package scanner

import (
	"fmt"
	"strconv"

	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/errors"
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Scanner struct {
	sources []rune
	tokens  []Token
	start   int
	current int
	line    int

	errReporter errors.Reporter
}

func NewScanner(sources string, errReporter errors.Reporter) *Scanner {
	return &Scanner{
		sources: []rune(sources),
		tokens:  nil,
		start:   0,
		current: 0,
		line:    1,

		errReporter: errReporter,
	}
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))
	return s.tokens
}

// scanToken is main token scanning function.
// gocyclo considers this function too difficult, but scanners are written always
// in this way.
//
//gocyclo:ignore
func (s *Scanner) scanToken() {
	currRune := s.advance()
	switch currRune {
	case '(':
		s.addNoLiteralToken(LEFTPAREN)
	case ')':
		s.addNoLiteralToken(RIGHTPAREN)
	case '{':
		s.addNoLiteralToken(LEFTBRACE)
	case '}':
		s.addNoLiteralToken(RIGHTBRACE)
	case ',':
		s.addNoLiteralToken(COMMA)
	case '.':
		s.addNoLiteralToken(DOT)
	case '-':
		s.addNoLiteralToken(MINUS)
	case '+':
		s.addNoLiteralToken(PLUS)
	case ';':
		s.addNoLiteralToken(SEMICOLON)
	case '*':
		s.addNoLiteralToken(STAR)
	case '!':
		if s.match('=') {
			s.addNoLiteralToken(BANGEQUAL)
		} else {
			s.addNoLiteralToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addNoLiteralToken(EQUALEQUAL)
		} else {
			s.addNoLiteralToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addNoLiteralToken(LESSEQUAL)
		} else {
			s.addNoLiteralToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addNoLiteralToken(GREATEREQUAL)
		} else {
			s.addNoLiteralToken(GREATER)
		}
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addNoLiteralToken(SLASH)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace.
	case '\n':
		s.line++
	case '"':
		s.str()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.number()
	default:
		if isAlpha(currRune) {
			s.identifier()
		} else {
			s.errReporter(
				s.line,
				fmt.Sprintf("Unexpected character: %s.", string(currRune)),
			)
		}
	}
}

func (s *Scanner) advance() rune {
	currRune := s.sources[s.current]
	s.current++
	return currRune
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if expected != s.sources[s.current] {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}
	return s.sources[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.sources) {
		return rune(0)
	}
	return s.sources[s.current+1]
}

func (s *Scanner) str() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		s.errReporter(s.line, "Unterminated string.")
		return
	}

	// The closing ".
	s.advance()

	// Trim the surrounding quotes.
	value := s.sources[s.start+1 : s.current-1]
	s.addToken(STRING, string(value))
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	// Look for fractional part.
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the ".".
		s.advance()

		// Repeat the same as for integer part.
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	value, err := strconv.ParseFloat(
		string(s.sources[s.start:s.current]),
		64,
	)
	if err != nil {
		s.errReporter(s.line, "Cannot parse float.")
	}
	s.addToken(NUMBER, value)
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := string(s.sources[s.start:s.current])
	kind, ok := keywords[text]
	if !ok {
		kind = IDENTIFIER
	}

	s.addNoLiteralToken(kind)
}

func (s *Scanner) addNoLiteralToken(kind TokenType) {
	s.addToken(kind, nil)
}

func (s *Scanner) addToken(kind TokenType, literal any) {
	runes := s.sources[s.start:s.current]
	token := NewToken(kind, string(runes), literal, s.line)
	s.tokens = append(s.tokens, token)
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.sources)
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') ||
		c == '_'
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}
