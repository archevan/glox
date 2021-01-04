package main

// -- AUTOGENERATED FILE -- (see scripts/generate_ast.py for details...)
// This is a simple implementation of the Visitor pattern from OOP
type ExprVisitor interface {
	VisitBinaryExpr(c *BinaryExpr)
	VisitGrouping(c *Grouping)
	VisitLiteral(c *Literal)
	VisitUnary(c *Unary)
	VisitVariable(c *Variable)
	VisitAssign(a *AssignExpr)
}

type Expr interface {
	accept(ExprVisitor)
}

// AssignExpr is a simple AST node
type AssignExpr struct {
	name Token
	val  Expr
}

// accept method stub for AssignExpr
func (a *AssignExpr) accept(v ExprVisitor) {
	v.VisitAssign(a)
}

// BinaryExpr is a simple type of AST node
type BinaryExpr struct {
	left  Expr
	op    Token
	right Expr
}

// accept method stub for BinaryExpr
func (c *BinaryExpr) accept(v ExprVisitor) {
	v.VisitBinaryExpr(c)
}

// Grouping is a simple type of AST node
type Grouping struct {
	exp Expr
}

// accept method stub for Grouping
func (c *Grouping) accept(v ExprVisitor) {
	v.VisitGrouping(c)
}

// Literal is a simple type of AST node
type Literal struct {
	val interface{}
}

// accept method stub for Literal
func (c *Literal) accept(v ExprVisitor) {
	v.VisitLiteral(c)
}

// Unary is a simple type of AST node
type Unary struct {
	op    Token
	right Expr
}

// accept method stub for Unary
func (c *Unary) accept(v ExprVisitor) {
	v.VisitUnary(c)
}

// Variable is a simple type of AST node
type Variable struct {
	name Token
}

// accept method stub for Variable
func (c *Variable) accept(v ExprVisitor) {
	v.VisitVariable(c)
}
