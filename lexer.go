package main

// A Lexer is an interface that can be scanned into a slice of tokens
type Lexer interface {
	scanTokens() []*Token
}

// Lex provides an implementation of Lexer that reads token from a string
type Lex struct {
	source string
	tokens []*Token
}

// Tokens gets a list of tokens from a Lex object
func (l *Lex) scanTokens() []*Token {
	return l.tokens
}
