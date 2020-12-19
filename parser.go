package main

import "os"

/*
The simple expression grammar for Lox is as follows (left-factored & unambiguous):
expression     → equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")" ;
*/

// Parser is a recursive descent parser
type Parser struct {
	inputTokens []*Token
	current     int
}

// NewParser is a factory function that creates a new Parser struct from a Lexer implementation
func NewParser(l Lexer) Parser {
	p := Parser{inputTokens: l.ScanTokens()}
	return p
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	exp := p.comparison()
	// left-associatively group equality expressions
	for p.match(BangEqual, EqualEqual) {
		op := p.previous()
		right := p.comparison()
		exp = &BinaryExpr{
			left:  exp,
			op:    *op,
			right: right,
		}
	}
	return exp
}

func (p *Parser) comparison() Expr {
	exp := p.term()
	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		op := p.previous()
		right := p.term()
		exp = &BinaryExpr{
			left:  exp,
			op:    *op,
			right: right,
		}
	}
	return exp
}

func (p *Parser) term() Expr {
	exp := p.factor()
	for p.match(Plus, Minus) {
		op := p.previous()
		right := p.factor()
		exp = &BinaryExpr{
			left:  exp,
			op:    *op,
			right: right,
		}
	}
	return exp
}

func (p *Parser) factor() Expr {
	exp := p.unary()
	for p.match(Star, Slash) {
		op := p.previous()
		right := p.unary()
		exp = &BinaryExpr{
			left:  exp,
			op:    *op,
			right: right,
		}
	}
	return exp
}

func (p *Parser) unary() Expr {
	if p.match(Bang, Minus) {
		op := p.previous()
		right := p.unary()
		return &Unary{
			op:    *op,
			right: right,
		}
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	// match a number of different types of literals
	switch {
	case p.match(FalseTok):
		return &Literal{val: false}
	case p.match(TrueTok):
		return &Literal{val: true}
	case p.match(NilTok):
		return &Literal{val: nil}
	case p.match(Number, StringTok):
		return &Literal{p.previous().literal}
	}
	// enforce matching parens
	if p.match(LeftParen) {
		exp := p.expression()
		p.consume(RightParen, "Expect ')' after expression")
		return &Grouping{exp: exp}
	}
	// TODO: implement below
	return &BinaryExpr{}
}

// consume matches the given token type or panic
func (p *Parser) consume(typ TokenType, fails string) {
	if p.check(typ) {
		p.advance()
		return
	}
	// error out and panic if we don't find a matching paren
	error(p.inputTokens[p.current].line, fails)
	os.Exit(64)
}

// match consumes the next token in the input stream if and only if
// it's equal to ANY of the given token types
func (p *Parser) match(typs ...TokenType) bool {
	for _, typ := range typs {
		if p.check(typ) {
			p.advance()
			return true
		}
	}
	return false
}

// advance consumes the next token in the input stream and returns it
func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// check compares the next token in the input stream to a given token type
func (p *Parser) check(typ TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.Peek().toktype == typ
}

// isAtEnd returns true if the next token is EOF
func (p *Parser) isAtEnd() bool {
	return p.Peek().toktype == EOF
}

// previous returns a pointer to the token we just consumed
func (p *Parser) previous() *Token {
	return p.inputTokens[p.current-1]
}

// Peek the next token from the input token string
func (p *Parser) Peek() *Token {
	return p.inputTokens[p.current]
}
