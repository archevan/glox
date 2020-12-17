@echo off
go clean
del /F /Q build\*
go build -o build\glx.exe main.go tokentypes.go lexer.go ast.go