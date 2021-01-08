package main

import (
	"errors"
)

/*
The simple statement grammar for Lox:
program		   → declaration* EOF ;
declaration	   → varDecl | statement ;
varDecl		   → "var" IDENTIFIER ( "=" expression )? ";" ;
statement	   → exprStmt | printStmt | ifstmt | block;
block          → "{" declaration* "}" ;
ifstmt         → "if" "(" expression ")" statement ("else" statement)? ;

The simple expression grammar for Lox is as follows (left-factored & unambiguous):
expression     → assignment ;
assignment     → IDENTIFIER "=" assignment
			   | logic_or;
logic_of	   → logic_and ("or" logic_and)* ;
logic_and	   → equality ("and" equality)* ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil"
               | IDENTIFIER
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

// Parse parses and returns a syntax tree (as a statement slice) for the given token stream
func (p *Parser) Parse() []Stmt {
	stmtList := make([]Stmt, 0)
	for !p.isAtEnd() {
		stmt := p.declaration()
		stmtList = append(stmtList, stmt)
	}
	return stmtList
}

// declaration parses a declaration from the token struct.
// ParseErrors are caught and handled here.
func (p *Parser) declaration() Stmt {
	if p.match(VarTok) {
		stmt, err := p.varDeclaration()
		if err != nil {
			p.synchronize()
			return nil
		}
		return stmt
	}
	stmt, err := p.statement()
	if err != nil {
		p.synchronize()
		return nil
	}
	return stmt
}

// varDeclaration parses a variable declaration with an optional initializer expression
func (p *Parser) varDeclaration() (Stmt, error) {
	var init Expr = nil
	err := p.consume(Identifier, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	name := p.previous()
	if p.match(Equal) {
		init, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	err = p.consume(Semicolon, "Expect semicolon after variable declaration.")
	if err != nil {
		return nil, err
	}
	return &VarStmt{
		name: name,
		init: init,
	}, nil
}

// statement() parses a sequence of tokens from the input stream that corresponds to a statement
func (p *Parser) statement() (Stmt, error) {
	if p.match(IfTok) {
		ifStmt, err := p.ifStatement()
		if err != nil {
			return nil, err
		}
		return ifStmt, nil
	}
	if p.match(PrintTok) {
		stmt, err := p.printStmt()
		if err != nil {
			return nil, err
		}
		return stmt, nil
	}
	if p.match(LeftBrace) {
		block, err := p.block()
		if err != nil {
			return nil, err
		}
		return &BlockStmt{statements: block}, nil
	}
	estmt, expErr := p.exprStmt()
	if expErr != nil {
		return nil, expErr
	}
	return estmt, nil
}

// ifStatement() parses an if statement structure from the token stream
// each call to ifStatement() parses an else structure which disambiguate the dangling else
func (p *Parser) ifStatement() (Stmt, error) {
	// parse if condition expression
	err := p.consume(LeftParen, "Expect '(' after 'if'")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	err = p.consume(RightParen, "Expect ')' after if condition")
	// parse 'then' part and pass along errors (if any)
	thenPart, err := p.statement()
	if err != nil {
		return nil, err
	}
	// parse 'else' part (if exists) and pass along errors (if any)
	var elsePart Stmt
	if p.match(Else) {
		elsePart, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return &IfStmt{
		thenPart: thenPart,
		elsePart: elsePart,
		exp:      condition,
	}, nil
}

// block() parses any number of statements inside of a lexical block from the token stream
func (p *Parser) block() ([]Stmt, error) {
	statements := make([]Stmt, 0)
	for !p.check(RightBrace) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	err := p.consume(RightBrace, "Expect '}' after block")
	if err != nil {
		return nil, err
	}
	return statements, nil
}

// printStmt() extracts a statement of the form PRINT <expression> from the token stream
func (p *Parser) printStmt() (Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return nil, err
	}
	semicolonMatchErr := p.consume(Semicolon, "Expect ';' after value")
	if semicolonMatchErr != nil {
		return nil, semicolonMatchErr
	}
	return &PrintStmt{
		exp: val,
	}, nil
}

// exprStmt() extracts an expression-statement from the input token stream
func (p *Parser) exprStmt() (Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return nil, err
	}
	semicolonMatchErr := p.consume(Semicolon, "Expect ';' after value")
	if semicolonMatchErr != nil {
		return nil, semicolonMatchErr
	}
	return &ExprStmt{
		exp: val,
	}, nil
}

func (p *Parser) expression() (Expr, error) {
	asg, err := p.assignment()
	if err != nil {
		return nil, err
	}
	return asg, nil
}

// assignment generates a Assign token for an assignment expr
// the return value is the expression that represents the assignment target
func (p *Parser) assignment() (Expr, error) {
	orRes, err := p.or()
	if err != nil {
		return nil, err
	}
	if p.match(Equal) {
		var val Expr
		eqtok := p.previous()
		val, err = p.assignment()
		if err != nil {
			return nil, err
		}
		if varTok, ok := orRes.(*Variable); ok {
			return &AssignExpr{
				name: varTok.name,
				val:  val,
			}, nil
		} else {
			errorTok(*eqtok, "Invalid assignment target")
		}
	}
	return orRes, nil
}

// or() parses any number of logical OR expressions
func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.match(OrTok) {
		op := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = &LogicalExpr{
			left:  expr,
			right: right,
			op:    *op,
		}
	}
	return expr, nil
}

func (p *Parser) and() (Expr, error) {
	eq, err := p.equality()
	if err != nil {
		return nil, err
	}
	for p.match(And) {
		op := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		eq = &LogicalExpr{
			left:  eq,
			right: right,
			op:    *op,
		}
	}
	return eq, nil
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
	// check for a variable usage
	if p.match(Identifier) {
		return &Variable{name: *p.previous()}, nil
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
