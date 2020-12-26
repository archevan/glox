@echo off
go clean
del /F /Q build\*
go build -o build\glx.exe main.go tokentypes.go parser.go lexer.go ast.go ast_printer.go