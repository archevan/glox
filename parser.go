package main

import (
	"errors"
	"fmt"
)

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
// Error handling is implemented using a "synchronization" technique
type Parser struct {
	inputTokens []*Token
	current     int
}

// NewParser is a factory function that creates a new Parser struct from a Lexer implementation
func NewParser(l Lexer) Parser {
	p := Parser{inputTokens: l.ScanTokens()}
	return p
}

// Parse parses and returns a syntax tree for the given token stream
func (p *Parser) Parse() Expr {
	exp, err := p.expression()
	if err != nil {
		return nil
	}
	return exp
}

func (p *Parser) expression() (Expr, error) {
	eq, err := p.equality()
	if err != nil {
		return nil, err
	}
	return eq, nil
}

func test() {
	var i int = 12
	test1(i)
}

func test1(i interface{}) {
	fmt.Println(i)
}

// equality() parses an equality structure from the input token stream
func (p *Parser) equality() (Expr, error) {
	exp, err := p.comparison()
	if err != nil {
		return nil, err
	}
	// left-associatively group equality expressions
	for p.match(BangEqual, EqualEqual) {
		op := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		exp = &BinaryExpr{
			left:  exp,
			op:    *op,
			right: right,
		}
	}
	return exp, nil
}

// comparison() parses a "comparison" structure from the input stream
func (p *Parser) comparison() (Expr, error) {
	exp, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		op := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		exp = &BinaryExpr{
			left:  exp,
			op:    *op,
			right: right,
		}
	}
	return exp, nil
}

// term() parses a "term" structure from the input token stream
func (p *Parser) term() (Expr, error) {
	exp, err := p.factor()
	if err != nil {
		// pass the buck
		return nil, err
	}
	for p.match(Plus, Minus) {
		op := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		exp = &BinaryExpr{
			left:  exp,
			op:    *op,
			right: right,
		}
	}
	return exp, nil
}

// factor() parses a "factor" structure from the input token stream
func (p *Parser) factor() (Expr, error) {
	exp, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(Star, Slash) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		exp = &BinaryExpr{
			left:  exp,
			op:    *op,
			right: right,
		}
	}
	return exp, nil
}

// unary() parses a unary op
func (p *Parser) unary() (Expr, error) {
	if p.match(Bang, Minus) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &Unary{
			op:    *op,
			right: right,
		}, nil
	}
	exp, err := p.primary()
	if err != nil {
		// pass the buck
		return nil, err
	}
	return exp, nil
}

func (p *Parser) primary() (Expr, error) {
	// match a number of different types of literals
	switch {
	case p.match(FalseTok):
		return &Literal{val: false}, nil
	case p.match(TrueTok):
		return &Literal{val: true}, nil
	case p.match(NilTok):
		return &Literal{val: nil}, nil
	case p.match(Number, StringTok):
		return &Literal{p.previous().literal}, nil
	}
	// enforce matching parens
	if p.match(LeftParen) {
		exp, err := p.expression()
		if err != nil {
			return nil, err
		}
		err = p.consume(RightParen, "Expect ')' after expression")
		if err != nil {
			// catch error thrown from consume
			return nil, err
		}
		return &Grouping{exp: exp}, nil
	}
	// current token can not be used to start an expression
	return nil, getError(*p.Peek(), "Expected expression.")
}

// consume matches the given token type or panic
// the error return type is similar to Java's throw it seems
func (p *Parser) consume(typ TokenType, fails string) error {
	if p.check(typ) {
		// expected char found, no error
		p.advance()
		return nil
	}
	return getError(*p.Peek(), fails)
}

// synchronize discard tokens from the parsers' input token steam
// until the beginning of a new statement is reached.
func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().toktype == Semicolon {
			return
		}

		switch p.Peek().toktype {
		case Class:
			return
		case Fun:
			return
		case VarTok:
			return
		case ForTok:
			return
		case IfTok:
			return
		case WhileTok:
			return
		case PrintTok:
			return
		case ReturnTok:
			return
		}
		// otherwise, discard current token.
		p.advance()
	}
}

// getError generates an error
func getError(tok Token, msg string) error {
	errorTok(tok, msg) // record invalid token
	return errors.New(msg)
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
