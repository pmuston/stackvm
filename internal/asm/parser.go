package asm

import (
	"fmt"
	"strconv"
)

// StatementType represents the type of a statement.
type StatementType int

const (
	StmtLabel StatementType = iota
	StmtInstruction
)

// Statement represents a parsed assembly statement.
type Statement struct {
	Type     StatementType
	Label    string      // For StmtLabel
	Opcode   string      // For StmtInstruction
	Operand  *Operand    // For StmtInstruction (optional)
	Line     int
	Column   int
}

// OperandType represents the type of an instruction operand.
type OperandType int

const (
	OperandNumber OperandType = iota
	OperandLabel
)

// Operand represents an instruction operand.
type Operand struct {
	Type       OperandType
	Number     int64   // For OperandNumber
	FloatValue float64 // For OperandNumber (if float)
	IsFloat    bool    // True if float, false if int
	Label      string  // For OperandLabel
}

// Parser parses tokens into an AST.
type Parser struct {
	tokens  []Token
	current int
}

// NewParser creates a new parser for the given tokens.
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

// Parse converts tokens into a list of statements.
func (p *Parser) Parse() ([]Statement, error) {
	statements := make([]Statement, 0)

	// Skip initial newlines
	p.skipNewlines()

	for !p.isAtEnd() {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		if stmt != nil {
			statements = append(statements, *stmt)
		}

		// Skip trailing newlines
		p.skipNewlines()
	}

	return statements, nil
}

func (p *Parser) parseStatement() (*Statement, error) {
	token := p.peek()

	switch token.Type {
	case TokenLabel:
		return p.parseLabelDef()
	case TokenIdent:
		return p.parseInstruction()
	case TokenNewline:
		p.advance()
		return nil, nil
	case TokenEOF:
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected token %s at %d:%d", token.Type, token.Line, token.Column)
	}
}

func (p *Parser) parseLabelDef() (*Statement, error) {
	token := p.expect(TokenLabel)
	if token == nil {
		return nil, fmt.Errorf("expected label")
	}

	stmt := &Statement{
		Type:   StmtLabel,
		Label:  token.Value,
		Line:   token.Line,
		Column: token.Column,
	}

	// Label definitions can be followed by a newline or another statement
	if p.peek().Type == TokenNewline {
		p.advance()
	}

	return stmt, nil
}

func (p *Parser) parseInstruction() (*Statement, error) {
	token := p.expect(TokenIdent)
	if token == nil {
		return nil, fmt.Errorf("expected instruction")
	}

	stmt := &Statement{
		Type:   StmtInstruction,
		Opcode: token.Value,
		Line:   token.Line,
		Column: token.Column,
	}

	// Check for operand
	if !p.isAtEnd() && p.peek().Type != TokenNewline && p.peek().Type != TokenEOF {
		operand, err := p.parseOperand()
		if err != nil {
			return nil, err
		}
		stmt.Operand = operand
	}

	// Consume newline if present
	if p.peek().Type == TokenNewline {
		p.advance()
	}

	return stmt, nil
}

func (p *Parser) parseOperand() (*Operand, error) {
	token := p.peek()

	switch token.Type {
	case TokenNumber:
		p.advance()
		// Try parsing as integer first
		if intVal, err := strconv.ParseInt(token.Value, 10, 64); err == nil {
			return &Operand{
				Type:    OperandNumber,
				Number:  intVal,
				IsFloat: false,
			}, nil
		}
		// Parse as float
		floatVal, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number '%s' at %d:%d: %v", token.Value, token.Line, token.Column, err)
		}
		return &Operand{
			Type:       OperandNumber,
			FloatValue: floatVal,
			IsFloat:    true,
		}, nil

	case TokenIdent:
		p.advance()
		return &Operand{
			Type:  OperandLabel,
			Label: token.Value,
		}, nil

	default:
		return nil, fmt.Errorf("expected operand (number or label) at %d:%d, got %s", token.Line, token.Column, token.Type)
	}
}

func (p *Parser) peek() Token {
	if p.current >= len(p.tokens) {
		return Token{Type: TokenEOF}
	}
	return p.tokens[p.current]
}

func (p *Parser) advance() Token {
	if p.current < len(p.tokens) {
		p.current++
	}
	return p.tokens[p.current-1]
}

func (p *Parser) expect(typ TokenType) *Token {
	if p.peek().Type == typ {
		token := p.advance()
		return &token
	}
	return nil
}

func (p *Parser) isAtEnd() bool {
	if p.current >= len(p.tokens) {
		return true
	}
	return p.tokens[p.current].Type == TokenEOF
}

func (p *Parser) skipNewlines() {
	for !p.isAtEnd() && p.peek().Type == TokenNewline {
		p.advance()
	}
}
