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
	resultVal interface{}
}

// RuntimeError is a wrapper around the "offending" token and its associated error message
type RuntimeError struct {
	tkn Token
	msg string
}

func (r RuntimeError) Error() string {
	return r.msg
}

// Interpret is the Interpreter type's public API that allows values to be interpreted
func (in *Interpreter) Interpret(e Expr) {
	val, err := in.evaluate(e)
	if err != nil {
		// detect any runtime errors
		// TODO: convert to a type switch when different types of errors are added
		if rtErr, ok := err.(RuntimeError); ok {
			runtimeError(rtErr)
		}
	} else {
		fmt.Println(in.stringify(val))
	}
}

// convert an evaluated Lox value into a string
func (in *Interpreter) stringify(val interface{}) string {
	if val == nil {
		return "nil"
	}
	if num, ok := val.(float64); ok {
		str := fmt.Sprintf("%.1f", num)
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

// isTruthy determines whether a given value will evalulate to true
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
