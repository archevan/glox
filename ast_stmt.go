package main

// -- AUTOGENERATED FILE -- (see scripts/generate_ast.py for details...)
// This is a simple implementation of the Visitor pattern from OOP
type StmtVisitor interface {
	VisitPrintStmt(c *PrintStmt)
	VisitExprStmt(c *ExprStmt)
	VisitVarStmt(c *VarStmt)
	VisitBlockStmt(b *BlockStmt)
	VisitIfStmt(i *IfStmt)
	VisitWhileStmt(w *WhileStmt)
	VisitFunctionStmt(f *FunctionStmt)
}

// IfStmt represents a branch with an optional else
type IfStmt struct {
	thenPart, elsePart Stmt
	exp                Expr
}

// accept method stub for an if statement
func (i *IfStmt) accept(v StmtVisitor) {
	v.VisitIfStmt(i)
}

// FunctionStmt represents a function declaration in the AST
type FunctionStmt struct {
	name   Token
	params []Token
	body   []Stmt
}

// accept method stub for an if statement
func (f *FunctionStmt) accept(v StmtVisitor) {
	v.VisitFunctionStmt(f)
}

// WhileStmt represents a simple loop structure in the AST
type WhileStmt struct {
	condition Expr
	statement Stmt
}

// accept method stub for an if statement
func (w *WhileStmt) accept(v StmtVisitor) {
	v.VisitWhileStmt(w)
}

// BlockStmt is a node that represents a list of statements
type BlockStmt struct {
	statements []Stmt
}

// accept method stub for BlockStmt
func (b *BlockStmt) accept(v StmtVisitor) {
	v.VisitBlockStmt(b)
}

type Stmt interface {
	accept(v StmtVisitor)
}

// PrintStmt is a simple type of AST node
type PrintStmt struct {
	exp Expr
}

// accept method stub for PrintStmt
func (c *PrintStmt) accept(v StmtVisitor) {
	v.VisitPrintStmt(c)
}

// ExprStmt is a simple type of AST node
type ExprStmt struct {
	exp Expr
}

// accept method stub for ExprStmt
func (c *ExprStmt) accept(v StmtVisitor) {
	v.VisitExprStmt(c)
}

// VarStmt is a simple type of AST node
type VarStmt struct {
	name *Token
	init Expr
}

// accept method stub for VarStmt
func (c *VarStmt) accept(v StmtVisitor) {
	v.VisitVarStmt(c)
}
