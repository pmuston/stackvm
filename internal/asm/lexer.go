package asm

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// TokenType represents the type of a token.
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenNewline
	TokenIdent      // Identifier (opcode or label reference)
	TokenLabel      // Label definition (ends with :)
	TokenNumber     // Numeric literal
	TokenComment    // Comment
)

// Token represents a lexical token.
type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q) at %d:%d", t.Type, t.Value, t.Line, t.Column)
}

func (tt TokenType) String() string {
	switch tt {
	case TokenEOF:
		return "EOF"
	case TokenNewline:
		return "NEWLINE"
	case TokenIdent:
		return "IDENT"
	case TokenLabel:
		return "LABEL"
	case TokenNumber:
		return "NUMBER"
	case TokenComment:
		return "COMMENT"
	default:
		return fmt.Sprintf("TokenType(%d)", tt)
	}
}

// Lexer tokenizes assembly source code.
type Lexer struct {
	source  string
	pos     int
	line    int
	column  int
	tokens  []Token
	current int
}

// NewLexer creates a new lexer for the given source.
func NewLexer(source string) *Lexer {
	return &Lexer{
		source: source,
		pos:    0,
		line:   1,
		column: 1,
	}
}

// Tokenize converts the source into tokens.
func (l *Lexer) Tokenize() ([]Token, error) {
	l.tokens = make([]Token, 0)

	for l.pos < len(l.source) {
		if err := l.scanToken(); err != nil {
			return nil, err
		}
	}

	// Add EOF token
	l.tokens = append(l.tokens, Token{
		Type:   TokenEOF,
		Line:   l.line,
		Column: l.column,
	})

	return l.tokens, nil
}

func (l *Lexer) scanToken() error {
	ch := l.peek()

	// Skip whitespace (except newlines)
	if ch == ' ' || ch == '\t' || ch == '\r' {
		l.advance()
		return nil
	}

	// Newline
	if ch == '\n' {
		l.emitToken(TokenNewline, "\n")
		l.advance()
		l.line++
		l.column = 1
		return nil
	}

	// Comments
	if ch == ';' || ch == '#' {
		l.scanComment()
		return nil
	}

	// Numbers (including negative)
	if unicode.IsDigit(rune(ch)) || (ch == '-' && l.pos+1 < len(l.source) && unicode.IsDigit(rune(l.source[l.pos+1]))) {
		return l.scanNumber()
	}

	// Identifiers and labels
	if unicode.IsLetter(rune(ch)) || ch == '_' {
		return l.scanIdentOrLabel()
	}

	return fmt.Errorf("unexpected character '%c' at %d:%d", ch, l.line, l.column)
}

func (l *Lexer) scanComment() {
	// Skip to end of line
	for l.pos < len(l.source) && l.source[l.pos] != '\n' {
		l.pos++
		l.column++
	}
}

func (l *Lexer) scanNumber() error {
	start := l.pos
	startCol := l.column

	// Handle negative sign
	if l.peek() == '-' {
		l.advance()
	}

	// Scan digits
	for l.pos < len(l.source) && (unicode.IsDigit(rune(l.peek())) || l.peek() == '.') {
		l.advance()
	}

	value := l.source[start:l.pos]

	// Validate number
	if strings.Contains(value, ".") {
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid float '%s' at %d:%d: %v", value, l.line, startCol, err)
		}
	} else {
		_, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer '%s' at %d:%d: %v", value, l.line, startCol, err)
		}
	}

	l.emitTokenAt(TokenNumber, value, l.line, startCol)
	return nil
}

func (l *Lexer) scanIdentOrLabel() error {
	start := l.pos
	startCol := l.column

	// Scan identifier characters
	for l.pos < len(l.source) {
		ch := l.peek()
		if unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_' {
			l.advance()
		} else {
			break
		}
	}

	value := l.source[start:l.pos]

	// Check if it's a label definition (followed by :)
	if l.pos < len(l.source) && l.peek() == ':' {
		l.advance() // consume ':'
		l.emitTokenAt(TokenLabel, value, l.line, startCol)
	} else {
		l.emitTokenAt(TokenIdent, value, l.line, startCol)
	}

	return nil
}

func (l *Lexer) peek() byte {
	if l.pos >= len(l.source) {
		return 0
	}
	return l.source[l.pos]
}

func (l *Lexer) advance() {
	if l.pos < len(l.source) {
		l.pos++
		l.column++
	}
}

func (l *Lexer) emitToken(typ TokenType, value string) {
	l.emitTokenAt(typ, value, l.line, l.column)
}

func (l *Lexer) emitTokenAt(typ TokenType, value string, line, col int) {
	l.tokens = append(l.tokens, Token{
		Type:   typ,
		Value:  value,
		Line:   line,
		Column: col,
	})
}
