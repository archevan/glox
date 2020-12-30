package main

// -- AUTOGENERATED FILE -- (see scripts/generate_ast.py for details...)
// This is a simple implementation of the Visitor pattern from OOP
type StmtVisitor interface {
	VisitExprStmt(c *ExprStmt)
	VisitPrintStmt(c *PrintStmt)
}

type Stmt interface {
	accept(StmtVisitor)
}

// ExprStmt is a simple type of AST node
type ExprStmt struct {
	exp Expr
}

// accept method stub for ExprStmt
func (c ExprStmt) accept(v StmtVisitor) {
	v.VisitExprStmt(&c)
}

// PrintStmt is a simple type of AST node
type PrintStmt struct {
	exp Expr
}

// accept method stub for PrintStmt
func (c PrintStmt) accept(v StmtVisitor) {
	v.VisitPrintStmt(&c)
}
