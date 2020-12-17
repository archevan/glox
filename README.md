# GLOX 🍣

### A _very_ **WIP** Lox interpreter written in Golang

This project is intended (almost) exclusively as a learning project.
The code is _heavily_ inspired by the Java source code presented in _Crafting Interpreters_ written by Bob Nystrom

#### usage

Run from file:

```
.\glx.exe [path-to-script]
```

Run the REPL:

```
.\glx.exe
```

#### misc. tool usage

Run the AST generator:
_From inside the scripts directory_

```
.\python.exe generate_ast.py [ouput-directory]
```

#### misc. notes

- The build script(s) in the scripts directory are for windows (sorry!)
- The build script(s) in the scripts directory are also FRAGILE!! A robust build process is something that I (intentionally) did not spend a long time working on. Use at your own risk.
- Run the build.bat script from **outside** from the main glox directory.
- The source is 100% Go so it should be pretty easy to build for other platforms
- I intend to keep up with the unit tests for the whole project to some extent in the files named '\*\_test.go'. A call to 'go test' should be all you need to invoke them.
