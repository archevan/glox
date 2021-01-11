package main

import "time"

/*
Native functions should be defined as types that implement that LoxCaller interface
*/

// LoxCaller encompasses any type that supported being called with arguments
type LoxCaller interface {
	arity() int
	call(in Interpreter, args []interface{}) interface{}
}

// GlobalFunctionClock is a native function wrapper that exposes clock() which returns a Unix time
type GlobalFunctionClock string

func (g *GlobalFunctionClock) arity() int {
	return 0
}

func (g *GlobalFunctionClock) String() string {
	return g.String()
}

func (g *GlobalFunctionClock) call(in *Interpreter, args []interface{}) interface{} {
	return time.Now().Unix()
}

func (r RuntimeError) Error() string {
	return r.msg
}
