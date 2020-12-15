/*
Package main implements a simple driver program to accept
command line args and run the rest of the compiler */
package main

// TODO: implement OS-specific constants

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	version = "v0.0.1"
)

// Global flag is set on error
var hasError bool

// Run a given string of code input could be entire script or a single line
func run(script string) {
	var lexer Lexer = NewLexScanner(script)
	toks := lexer.ScanTokens()
	fmt.Println("Token Stream:")
	// For now, just print each token
	for index, tok := range toks {
		fmt.Printf("%v: %v\n", index, tok)
	}
}

// report an error on a line number 'line'
func error(line int, msg string) {
	report(line, "", msg)
}

// Report an error at a given line number
func report(line int, where, msg string) {
	fmt.Printf("[line %d] Error %v: %v\n", line, where, msg)
	hasError = true
}

// Read a given lox file at 'path' into a string and execute it
func runFile(path string) {
	fcontents, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Can't open file at [%v].\n", path)
	}
	fstring := string(fcontents)
	// execute the resulting string
	run(fstring)
	// did we find an error along the way
	if hasError {
		os.Exit(65)
	}
}

// Trim the last 'num' character from 'str'
func trimSuffix(str string, num int) string {
	return str[:len(str)-num]
}

// simple REPL implementation, input is executed line-by-line
func runPrompt() {
	fmt.Println("Hey. Lox Interpreter", version, "(type 'exit' to leave)")
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := r.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading line.")
		}
		// remove newline '\r\n' (windows) from input
		line = trimSuffix(line, 2)
		if line == "exit" {
			fmt.Println("Bye bye.")
			break
		}
		if line != "" {
			run(line)
			hasError = false // reset error flag in interactive mode
		}
	}
}

// Application entry point
func main() {
	// accept an input script
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Println("usage: glox.exe [script]")
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}
