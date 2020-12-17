package main

import (
	"reflect"
	"testing"
)

// TestNewLexScanner tests NewLexScanner factory function
func TestNewLexScanner(t *testing.T) {
	l := NewLexScanner("test")
	if l.line != 1 || l.current != 0 || l.start != 0 {
		t.Errorf("NewLexScanner() init failed (line == 1, current == 0, start == 0): (got %v, got %v, got %v)\n", l.line, l.current, l.start)
	}
}

// compareTokenSlices is a helper that uses the 'reflect' library to compare two Token pointer slices
func compareTokenSlices(a, b []*Token) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		// is the corresponding element in b == v
		if !reflect.DeepEqual(b[i], v) {
			return false
		}
	}
	return true
}

// Test the ouput of an empty lexer
func TestEmptyScanToken(t *testing.T) {
	expected := []*Token{&Token{toktype: EOF, line: 1, lexeme: "END OF FILE"}}
	emptyLex := NewLexScanner("")
	emptyLex.ScanTokens()
	if !compareTokenSlices(emptyLex.tokens, expected) {
		t.Errorf("Empty lexer scanned incorrect tokens. %v, %v\n", expected[0], emptyLex.tokens[0])
	}
}

// Test the ouput of an empty lexer
func TestArithScanToken(t *testing.T) {
	expected := []*Token{
		// NUMBER tokens literals are *always* floating point values
		&Token{toktype: Number, line: 1, lexeme: "2", literal: 2.0},
		&Token{toktype: Plus, line: 1, lexeme: "+"},
		&Token{toktype: Number, line: 1, lexeme: "4", literal: 4.0},
		&Token{toktype: EOF, line: 1, lexeme: "END OF FILE"},
	}
	arithLex := NewLexScanner("2 + 4")
	arithLex.ScanTokens()
	if !compareTokenSlices(arithLex.tokens, expected) {
		t.Errorf("Arithmetic lexer scanned incorrect tokens.\nWanted: %v\nGot: %v\n", expected, arithLex.tokens)
	}
}
