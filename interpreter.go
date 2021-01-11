package main

import (
	"fmt"
	"reflect"
	"strings"
)

// Interpreter is an implementation of the Visitor interface to recursively
// walk the syntax tree generated by the parser. Tree-walk interpreter.
type Interpreter struct {
	// Lox return values are represented with an empty interface
	resultVal    interface{}
	globals, env *Environment
}

// RuntimeError is a wrapper around the "offending" token and its associated error message
type RuntimeError struct {
	tkn Token
	msg string
}

// NewInterpreter returns a properly initialized interpreter structure
func NewInterpreter() *Interpreter {
	newEnv := NewEnvironment(nil)
	newInt := &Interpreter{
		globals: newEnv,
		env:     newEnv,
	}
	// define native functions in the new interpreter's global environment
	newInt.globals.Define("clock", GlobalFunctionClock("<native clock fn>"))
	return newInt
}

// Interpret is the Interpreter type's public API that allows values to be interpreted
func (in *Interpreter) Interpret(stmtList []Stmt) {
	for _, stmt := range stmtList {
		err := in.execute(stmt)
		if err != nil {
			// catch error type
			switch errtyp := err.(type) {
			case RuntimeError:
				runtimeError(errtyp)
				return
			}
		}
	}
}

// execute() is the equivalent of evaluate() for statements
func (in *Interpreter) execute(s Stmt) error {
	s.accept(in)
	if err, ok := in.resultVal.(error); ok {
		return err
	}
	return nil
}

// convert an evaluated Lox value into a string
func (in *Interpreter) stringify(val interface{}) string {
	if val == nil {
		return "nil"
	}
	if num, ok := val.(float64); ok {
		str := fmt.Sprintf("%.1f", num)
		// strip decimal from int floats
		if strings.HasSuffix(str, ".0") {
			str = str[:len(str)-2]
		}
		return str
	}
	return fmt.Sprintf("%v", val)
}

// allow a given expression to call the correct Visit method for its type
func (in *Interpreter) evaluate(e Expr) (interface{}, error) {
	// each expression "accepts" the interpreter struct (which implements the Visitor interface)
	e.accept(in)
	// catch any runtime errors
	if err, ok := in.resultVal.(error); ok {
		// DO SOMETHING WITH THE ERROR
		return nil, err
	}
	return in.resultVal, nil
}

// VisitCall executes a call structure in the input AST
func (in *Interpreter) VisitCall(c *CallExpr) {
	callee, err := in.evaluate(c.callee)
	if err != nil {
		in.resultVal = err
		return
	}
	// eval args
	evalArgs := make([]interface{}, 0)
	for _, arg := range c.arguments {
		evalArg, err := in.evaluate(arg)
		if err != nil {
			in.resultVal = err
			return
		}
		evalArgs = append(evalArgs, evalArg)
	}
	// callee MUST BE callable
	function, ok := callee.(LoxFunction)
	if !ok {
		// throw a RuntimeError
		in.resultVal = &RuntimeError{
			tkn: c.paren,
			msg: "Can only call functions and classes.",
		}
		return
	}
	// correct number of arguments MUST BE given
	if len(evalArgs) != function.arity() {
		in.resultVal = &RuntimeError{
			tkn: c.paren,
			msg: fmt.Sprintf("Expected %d arguments but got %d.", function.arity(), len(evalArgs)),
		}
		return
	}
	// call the given function without
	in.resultVal = function.call(in, evalArgs)
}

// VisitFunctionStmt creates a binding in the interpreter's current environment between the function's name
// and its corresponding LoxFunction values when a variable declaration is encountered. This creates a "callable"
// interface (LoxFunction) for the given FunctionStmt node that can be invoked using the call() method later in the tree-walk.
func (in *Interpreter) VisitFunctionStmt(f *FunctionStmt) {
	function := LoxFunction(*f)
	in.env.Define(f.name.lexeme, function)
}

// VisitVariable evaluates a variable expression to its corresponding value in the symbol table
func (in *Interpreter) VisitAssign(a *AssignExpr) {
	val, err := in.evaluate(a.val)
	if err != nil {
		in.resultVal = err
		return
	}
	err = in.env.Assign(a.name, val)
	if err != nil {
		in.resultVal = err
	} else {
		in.resultVal = val
	}
}

// VisitWhileStmt executes a while statement in the input syntax tree
// this is a thin wrapper around Go's for loop
func (in *Interpreter) VisitWhileStmt(w *WhileStmt) {
	condition, err := in.evaluate(w.condition)
	if err != nil {
		in.resultVal = err
		return
	}
	for in.isTruthy(condition) {
		err = in.execute(w.statement)
		if err != nil {
			in.resultVal = err
			return
		}
		// check condition again
		condition, err = in.evaluate(w.condition)
		if err != nil {
			in.resultVal = err
			return
		}
	}
	in.resultVal = nil
}

// VisitVariable evaluates a variable expression to its corresponding value in the symbol table
func (in *Interpreter) VisitVariable(v *Variable) {
	val, err := in.env.Get(v.name)
	if err != nil {
		in.resultVal = err
		return
	}
	in.resultVal = val
}

// VisitIfStmt interprets an if statement
func (in *Interpreter) VisitIfStmt(i *IfStmt) {
	condition, err := in.evaluate(i.exp)
	if err != nil {
		in.resultVal = err
		return
	}
	if in.isTruthy(condition) {
		if err = in.execute(i.thenPart); err != nil {
			in.resultVal = err
			return
		}
	} else if i.elsePart != nil {
		// execute the else statement if it exists
		if err = in.execute(i.elsePart); err != nil {
			in.resultVal = err
			return
		}
	}
	in.resultVal = nil
}

// VisitLogical() interprets the expressions given as
// arguments to a logical expression (short-circuiting if necessary)
// a value with appropriate "truthy-ness" will be returned
func (in *Interpreter) VisitLogical(l *LogicalExpr) {
	left, err := in.evaluate(l.left)
	if err != nil {
		in.resultVal = err
		return
	}
	// the following conditional block allows logical operators to "short circuit"
	if l.op.toktype == OrTok {
		// OR token with true left expr
		if in.isTruthy(left) {
			in.resultVal = left
			return
		}
	} else {
		// AND token with false left expr
		if !in.isTruthy(left) {
			in.resultVal = left
			return
		}
	}
	right, err := in.evaluate(l.right)
	if err != nil {
		in.resultVal = err
		return
	}
	in.resultVal = right
}

// VisitBlockStmt evaluates the statements inside of a lexical block
func (in *Interpreter) VisitBlockStmt(b *BlockStmt) {
	// execute block statements in a new environment
	in.executeBlock(b.statements, NewEnvironment(in.env))
}

// execute a given list of statements in the given environment
func (in *Interpreter) executeBlock(stmts []Stmt, newEnv *Environment) {
	// "push" the given environment onto the top of the scope chain
	newEnv.enclosing = in.env
	in.env = newEnv
	for _, statement := range stmts {
		err := in.execute(statement)
		if err != nil {
			in.resultVal = err
			in.env = in.env.enclosing
			return
		}
	}
	// pop the innermost scope off of the "scope chain"
	in.env = in.env.enclosing
}

// VisitVarStmt inserts a variable binding into the current environment
func (in *Interpreter) VisitVarStmt(v *VarStmt) {
	var val interface{}
	var err error
	if v.init != nil {
		val, err = in.evaluate(v.init)
		if err != nil {
			in.resultVal = RuntimeError{
				tkn: *v.name,
				msg: "Can't evaluate variable init expression.",
			}
			return
		}
	}
	// add new binding to current environment
	in.env.Define(v.name.lexeme, val)
}

// VisitBinaryExpr interprets any given binary expression
func (in *Interpreter) VisitBinaryExpr(b *BinaryExpr) {
	// evaluate left and right operand expressions, passing errors up the call stack as needed
	left, lerr := in.evaluate(b.left)
	if lerr != nil {
		in.resultVal = lerr
		return
	}
	right, rerr := in.evaluate(b.right)
	if rerr != nil {
		in.resultVal = rerr
		return
	}
	switch b.op.toktype {
	case Greater:
		in.checkNumberOperands(b.op, left, right)
		if _, ok := in.resultVal.(error); ok {
			return
		}
		leftd := left.(float64)
		rightd := right.(float64)
		in.resultVal = leftd > rightd
	case GreaterEqual:
		in.checkNumberOperands(b.op, left, right)
		if _, ok := in.resultVal.(error); ok {
			return
		}
		leftd := left.(float64)
		rightd := right.(float64)
		in.resultVal = leftd >= rightd
	case Less:
		in.checkNumberOperands(b.op, left, right)
		if _, ok := in.resultVal.(error); ok {
			return
		}
		leftd := left.(float64)
		rightd := right.(float64)
		in.resultVal = leftd < rightd
	case LessEqual:
		in.checkNumberOperands(b.op, left, right)
		if _, ok := in.resultVal.(error); ok {
			return
		}
		leftd := left.(float64)
		rightd := right.(float64)
		in.resultVal = leftd <= rightd
	case Minus:
		in.checkNumberOperands(b.op, left, right)
		if _, ok := in.resultVal.(error); ok {
			return
		}
		leftd := left.(float64)
		rightd := right.(float64)
		in.resultVal = leftd - rightd
	case Slash:
		in.checkNumberOperands(b.op, left, right)
		if _, ok := in.resultVal.(error); ok {
			return
		}
		leftd := left.(float64)
		rightd := right.(float64)
		in.resultVal = leftd / rightd
	case Star:
		in.checkNumberOperands(b.op, left, right)
		if _, ok := in.resultVal.(error); ok {
			return
		}
		leftd := left.(float64)
		rightd := right.(float64)
		in.resultVal = leftd * rightd
	case Plus:
		// plus can be applied to both numbers (doubles) and strings
		// this solution only looks at the type of the expression's left operand
		leftd, lOk := left.(float64)
		rightd, rOk := right.(float64)
		leftstr, lStrOk := left.(string)
		rightstr, rStrOk := right.(string)
		switch {
		case lOk && rOk:
			in.resultVal = leftd + rightd
		case lStrOk && rStrOk:
			in.resultVal = leftstr + rightstr
		default:
			in.resultVal = RuntimeError{
				tkn: b.op,
				msg: "Addition operands must both be numbers or strings",
			}
		}
	}
	// TODO: implement more binary operations
}

// isEqual checks whether two given values are equal.
// behavior is similar to Go's == but has support for nil values
func (in *Interpreter) isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	// same as Go's == for strings, booleans, and doubles (float64)
	return reflect.DeepEqual(a, b)
}

// VisitGrouping interprets any given Grouping expression
func (in *Interpreter) VisitGrouping(g *Grouping) {
	exp, err := in.evaluate(g.exp)
	if err != nil {
		in.resultVal = err
		return
	}
	in.resultVal = exp
}

// VisitLiteral interprets any given Literal expression
func (in *Interpreter) VisitLiteral(l *Literal) {
	in.resultVal = l.val
}

// VisitExprStmt interprets an expression-statement
func (in *Interpreter) VisitExprStmt(estmt *ExprStmt) {
	val, err := in.evaluate(estmt.exp)
	if err != nil {
		in.resultVal = RuntimeError{
			msg: "Can't eval expr statement",
		}
	}
	in.resultVal = val
}

// VisitPrintStmt interprets an print statement
func (in *Interpreter) VisitPrintStmt(pstmt *PrintStmt) {
	val, err := in.evaluate(pstmt.exp)
	if err != nil {
		in.resultVal = RuntimeError{
			msg: "Can't eval print statement",
		}
		return
	}
	fmt.Println(in.stringify(val))
}

// isTruthy determines whether a given value will evaluate to true
// nil and false both eval to false, everything else evaluates to true
func (in *Interpreter) isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	}
	return true
}

// VisitUnary interprets any given Unary expression
func (in *Interpreter) VisitUnary(u *Unary) {
	right, err := in.evaluate(u.right)
	if err != nil {
		in.resultVal = err
		return
	}
	switch u.op.toktype {
	case Minus:
		in.checkNumberOperand(u.op, right)
		if _, ok := in.resultVal.(error); ok {
			// if result value if an error val, unwind
			return
		}
		in.resultVal = -right.(float64)
	case Bang:
		in.resultVal = !in.isTruthy(right)
	}
}

// checkNumberOperand sets the result value of the current expression when operand is NaN, otherwise NOP
func (in *Interpreter) checkNumberOperand(op Token, operand interface{}) {
	if _, ok := operand.(float64); ok {
		return
	}
	in.resultVal = RuntimeError{
		tkn: op,
		msg: "operand must be a number",
	}
}

// checkNumberOperands sets the result value of the current expression to an error value if either operand is NaN, otherwise NOP
func (in *Interpreter) checkNumberOperands(op Token, left, right interface{}) {
	_, lok := left.(float64)
	_, rok := right.(float64)
	if lok && rok {
		return
	}
	in.resultVal = RuntimeError{
		tkn: op,
		msg: "both operands must be numbers",
	}
}
