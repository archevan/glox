@echo off
go clean
del /F /Q build\*
go build -o build\glox.exe main.go tokentypes.go lexer.go