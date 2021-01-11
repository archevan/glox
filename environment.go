package main

// Environment DOES NOT have usable default values. Please initialize with a call to New()
type Environment struct {
	enclosing *Environment // pointer to enclosing scope
	bindings  map[string]interface{}
}

// NewEnvironment() returns a pointer to a properly initialized Environment
func NewEnvironment(enclosing *Environment) *Environment {
	env := &Environment{
		enclosing: enclosing,
		bindings:  make(map[string]interface{}),
	}
	return env
}

// Define() adds a new entry to the given environment bindings
func (e *Environment) Define(name string, val interface{}) {
	e.bindings[name] = val
}

// Get() searches the scope chain for a given name and throws an error if it's not found
func (e *Environment) Get(name Token) (interface{}, error) {
	if val, ok := e.bindings[name.lexeme]; ok {
		return val, nil
	}
	// name not found in innermost scope, check enclosing scopes
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	// variable not found
	return nil, RuntimeError{
		tkn: name,
		msg: "Undefined variable " + name.lexeme + ".",
	}
}

// Assign() attempts to change the value bound to 'name' in the scope chain, throws a RuntimeError if 'name' isn't present.
func (e *Environment) Assign(name Token, val interface{}) error {
	if _, ok := e.bindings[name.lexeme]; ok {
		e.bindings[name.lexeme] = val
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.Assign(name, val)
	}
	return RuntimeError{
		tkn: name,
		msg: "Undefined Variable " + name.lexeme + ".",
	}
}
