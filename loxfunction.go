package main

// LoxFunction is a wrapper around a FunctionStmt AST node that implements the LoxCaller interface.
// In other words, LoxFunction keeps the logic related to binding arguments and parameters out of the parser.
type LoxFunction FunctionStmt

// the call method allows a FunctionStmt body to be executed in a correctly configured environment.
func (l *LoxFunction) call(in *Interpreter, args []interface{}) interface{} {
	// create new environment from interpreter's global environment
	env := NewEnvironment(in.globals)
	// create mapping between parameters and arguments to function
	for i, param := range l.params {
		env.Define(param.lexeme, args[i])
	}
	// execute function body inside newly-created environment
	in.executeBlock(l.body, env)
	// TODO: implement return values from function calls
	return nil
}

// arity returns the required number of arguments needed to call the current LoxFunction
func (l *LoxFunction) arity() int {
	return len(l.params)
}

// simple String() representation
func (l *LoxFunction) String() string {
	return "<fn " + l.name.lexeme + ">"
}
