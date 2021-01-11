package main

import (
	"fmt"
	"strings"
)

// ASTPrinter is an implementation of a visitor interface that "pretty-prints" AST nodes.
// Each Visit method generates the correct call to the parenthesize() method
type ASTPrinter struct {
	str string
}

func (a2 *ASTPrinter) VisitCall(c *CallExpr) {
	panic("implement me")
}

func (a2 *ASTPrinter) VisitLogical(l *LogicalExpr) {
	panic("implement me")
}

func (a2 *ASTPrinter) VisitAssign(a *AssignExpr) {
	panic("implement me" + a.name.lexeme)
}

func (a *ASTPrinter) VisitVariable(c *Variable) {
	panic("implement me" + c.name.lexeme)
}

// Print passes the ASTPrinter visitor to an Expr
func (a *ASTPrinter) Print(exp Expr) string {
	exp.accept(a)
	return a.String()
}

// VisitBinaryExpr pprints a binary expression
func (a *ASTPrinter) VisitBinaryExpr(b *BinaryExpr) {
	a.parenthesize(b.op.lexeme, b.left, b.right)
}

// VisitGrouping pprints a grouped expression
func (a *ASTPrinter) VisitGrouping(g *Grouping) {
	a.parenthesize("group", g.exp)
}

// VisitLiteral pprints a literal expr
func (a *ASTPrinter) VisitLiteral(l *Literal) {
	if l.val == nil {
		a.str = "nil"
	}
	switch lit := l.val.(type) {
	case float64:
		a.str = fmt.Sprintf("%f", lit)
	case string:
		a.str = lit
	}
}

// VisitUnary pprints a unary expression
func (a *ASTPrinter) VisitUnary(u *Unary) {
	a.parenthesize(u.op.lexeme, u.right)
}

// parenthesize prints the name of an AST node and pprints its expression operands
func (a *ASTPrinter) parenthesize(name string, exps ...Expr) {
	var build strings.Builder
	build.WriteByte('(')
	build.WriteString(name)
	for _, exp := range exps {
		build.WriteByte(' ')
		exp.accept(a)
		build.WriteString(a.String())
	}
	build.WriteByte(')')
	a.str = build.String()
}

// Get the string representation for the Expr to be printed
func (a *ASTPrinter) String() string {
	return a.str
}
