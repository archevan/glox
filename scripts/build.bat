@echo off
go clean
del /F /Q build\*
go build -o build\glx.exe main.go tokentypes.go environment.go parser.go lexer.go ast_expr.go ast_stmt.go ast_printer.go interpreter.go