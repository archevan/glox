package main

// Expr is a base empty expression struct
// I'm using this empty struct type as a stand-in for an abstract class. Change my mind (or open a PR).
type Expr struct{}

// BinaryExpr is an expression with two operands and an operator
type BinaryExpr struct {
	Expr
	operator    Token
	left, right Expr
}
