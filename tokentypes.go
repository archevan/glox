/*
Package main stores token information */
package main

import "fmt"

// Each token type is assigned a unique int value
const (
	// single character tokens
	LeftParen = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	// one or two character tokens
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	// literals
	Identifier
	StringTok
	Number

	// keywords
	And
	Class
	Else
	FalseTok
	Fun
	ForTok
	IfTok
	NilTok
	OrTok
	PrintTok
	ReturnTok
	Super
	ThisTok
	TrueTok
	VarTok
	WhileTok

	// End of File
	EOF
)

// TokenType is an "enum-like" wrapper for the constants above
type TokenType int

/*
Token is a simple class to hold information about each encountered token
in the input stream. Line number, literal value, and tokentype are stored (among others) */
type Token struct {
	toktype TokenType
	lexeme  string
	literal interface{}
	line    int
}

// simple string representation for a token
func (t *Token) String() string {
	if t.toktype == EOF {
		t.lexeme = "END OF FILE"
	}
	return fmt.Sprintf("[TOKEN: %5v, %12s, %5v]", t.toktype, t.lexeme, t.line)
}
