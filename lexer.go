package main

import "strconv"

// A Lexer is an interface that can be scanned into a slice of tokens
type Lexer interface {
	ScanTokens() []*Token
}

// LexScanner provides an implementation of Lexer that reads token from a string
// LexScanner.Init() MUST be called before a LexScanner object is used
type LexScanner struct {
	reserved             map[string]TokenType
	source               string
	start, current, line int
	tokens               []*Token
}

// ScanTokens gets a list of tokens from a Lex object
func (l *LexScanner) ScanTokens() []*Token {
	for !l.isAtEnd() {
		l.start = l.current
		l.scanToken()
	}
	// add EOF token
	l.addToken(EOF, nil)
	return l.tokens
}

// NewLexScanner is a simple factory function that
// creates LexScanner objects and returns pointers to them
func NewLexScanner(inputStr string) *LexScanner {
	// initialize the reserved word map
	m := map[string]TokenType{
		"and":    And,
		"class":  Class,
		"else":   Else,
		"false":  FalseTok,
		"for":    ForTok,
		"fun":    Fun,
		"if":     IfTok,
		"nil":    NilTok,
		"or":     OrTok,
		"print":  PrintTok,
		"return": ReturnTok,
		"super":  Super,
		"this":   ThisTok,
		"true":   TrueTok,
		"var":    VarTok,
		"while":  WhileTok,
	}
	return &LexScanner{line: 1, source: inputStr, reserved: m}
}

// Has our scanner class reached the end of source string ?
func (l *LexScanner) isAtEnd() bool {
	return l.current >= len(l.source)
}

// advance gets the next character (byte) from the source
func (l *LexScanner) advance() byte {
	l.current++
	return l.source[l.current-1]
}

// add a new token to the token list the substring of
// source from start:current is yanked and stored as the token's lexeme
func (l *LexScanner) addToken(tok TokenType, lit interface{}) {
	text := l.source[l.start:l.current]
	newtok := &Token{toktype: tok, literal: lit, lexeme: text, line: l.line}
	l.tokens = append(l.tokens, newtok)
}

// the "big switch" scans individual tokens. the string
// contained at source[start:current] is the current token
func (l *LexScanner) scanToken() {
	c := l.advance()
	switch c {
	// single character tokens
	case '(':
		l.addToken(LeftParen, nil)
	case ')':
		l.addToken(RightParen, nil)
	case '{':
		l.addToken(LeftBrace, nil)
	case '}':
		l.addToken(RightBrace, nil)
	case ',':
		l.addToken(Comma, nil)
	case '.':
		l.addToken(Dot, nil)
	case '-':
		l.addToken(Minus, nil)
	case '+':
		l.addToken(Plus, nil)
	case ';':
		l.addToken(Semicolon, nil)
	case '*':
		l.addToken(Star, nil)
	case '!':
		tmp := Bang
		// lookahead by one character
		if l.match('=') {
			tmp = BangEqual
		}
		l.addToken(TokenType(tmp), nil)
	case '=':
		tmp := Equal
		if l.match('=') {
			tmp = EqualEqual
		}
		l.addToken(TokenType(tmp), nil)
	case '<':
		tmp := Less
		if l.match('=') {
			tmp = LessEqual
		}
		l.addToken(TokenType(tmp), nil)
	case '>':
		tmp := Greater
		if l.match('=') {
			tmp = GreaterEqual
		}
		l.addToken(TokenType(tmp), nil)
	case '/':
		if l.match('/') {
			// this is comment, discard everything until a newline
			for l.peek() != '\n' && !l.isAtEnd() {
				l.advance()
			}
		} else {
			l.addToken(Slash, nil)
		}
	case '"':
		l.string()
	case '\n':
		l.line++
	case ' ':
	case '\r':
	case '\t':
	default:
		if isADigit(c) {
			l.number()
		} else if isAlphaNumeric(c) {
			l.identifier()
		} else {
			error(l.line, "Unexpected character.")
		}
	}
}

// identifer() scans an identifer from the input stream
func (l *LexScanner) identifier() {
	for isAlphaNumeric(l.peek()) {
		l.advance()
	}
	text := l.source[l.start:l.current]
	typ, prs := l.reserved[text]
	// if the selected identifer is NOT a reserved word, then its an identifier
	if !prs {
		typ = Identifier
	}
	l.addToken(typ, nil)
}

// isAlpha returns true if the given character is alphabetical OR an underscore... false otherwise
func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

// isAlphaNumeric returns true if the given character is alphabetical OR a digit... false otherwise
func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isADigit(c)
}

// number() scans a number from the input stream
func (l *LexScanner) number() {
	for isADigit(l.peek()) {
		l.advance()
	}
	if l.peek() == '.' && isADigit(l.peekNext()) {
		l.advance()
		for isADigit(l.peek()) {
			l.advance()
		}
	}
	f, err := strconv.ParseFloat(l.source[l.start:l.current], 64)
	if err != nil {
		error(l.line, "Error reading floating point value")
	}
	l.addToken(Number, f)
}

// isADigit
func isADigit(c byte) bool {
	return (c >= '0' && c <= '9')
}

// string() scans a string form the input stream input a token
func (l *LexScanner) string() {
	// move 'current' pointer across the string
	// while maintaining the line number correctly
	for l.peek() != '"' && !l.isAtEnd() {
		if l.peek() == '\n' {
			l.line++
		}
		l.advance()
	}
	if l.isAtEnd() {
		error(l.line, "Unterminated String.")
	}
	l.advance()
	// trim quotes + create token
	val := l.source[l.start+1 : l.current-1]
	l.addToken(StringTok, val)
}

// match is a simple lookahead method that consumes
// the next character iff it's the character we're expecting
// thanks to Bob Nystrom for the code !
func (l *LexScanner) match(expected byte) bool {
	if l.isAtEnd() {
		return false
	}
	if l.source[l.current] != expected {
		return false
	}
	// we found what we're looking for, advance current pointer
	l.current++
	return true
}

// look at the next character in the source stream
func (l *LexScanner) peek() byte {
	if l.isAtEnd() {
		// Go Quirk: \0 is illegal.... use \000 instead for all the Cstrings peeps out there
		return '\000'
	}
	return l.source[l.current]
}

func (l *LexScanner) peekNext() byte {
	// there is no next to peek
	if (l.current + 1) >= len(l.source) {
		return '\000'
	}
	return l.source[l.current+1]
}
