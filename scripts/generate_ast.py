import sys
import os

# Abstract Syntax Tree Generator for Go
#
# Input: list of strings specifying the different
#        types of AST node types in the following
#        format:
#           "TYPE STRUCT NAME : FIELD-NAME TYPE, FIELD-NAME TYPE, ..."
# Output: an autogenerated .go file (called 'ast_*.go') in the current directory


def main():
    sys.argv.pop(0)  # remove first item (current exec path) from input list
    if len(sys.argv) != 1:
        print("usage: .\python.exe generate_ast.py [output_directory]")
        sys.exit(64)
    output_dir = sys.argv[0]
    output_file_name = "ast_expr.go"
    current_path = os.path.dirname(__file__)  # get current directory
    output_path = os.path.relpath(
        output_dir, current_path) + "\\" + output_file_name
    print("generating AST struct types...")
    print(f"writing to file at {output_path}")
    ts = ["BinaryExpr   :  left Expr, op Token, right Expr",
          "Grouping     :  exp Expr",
          "Literal      :  val interface{}",
          "Unary        :  op Token, right Expr",
          "Variable     :  name Token"]
    with open(output_path, 'w+') as f:
        define_ast(f, ts)


def define_ast(outfile, typestrings):
    write_preamble(outfile)
    # write out Visitor interface
    outfile.write("type Visitor interface {\n")
    for typ in typestrings:
        class_name = typ.split(':')[0].strip()
        outfile.write(f"\tVisit{class_name}(c *{class_name})\n")
    outfile.write("}\n")
    # write out Expr interface
    outfile.write("\ntype Expr interface {\n")
    outfile.write("\taccept(Visitor)\n")
    outfile.write("}\n")
    # write out classes for each type of AST node
    for typ in typestrings:
        class_name = typ.split(':')[0].strip()
        field_list = typ.split(':')[1].strip().split(',')
        write_class(outfile, class_name, field_list)


# write a Go struct type out to a given file
def write_class(outfile, class_name, field_list):
    outfile.write(f"\n// {class_name} is a simple type of AST node")
    outfile.write("\ntype " + class_name + " struct {")
    for field in field_list:
        outfile.write(f"\n\t{field.lstrip()}")
    outfile.write("\n}\n")
    # write out visitor accept method for current class
    outfile.write(f"\n// accept method stub for {class_name}\n")
    outfile.write("func (c * " + class_name + ") accept(v Visitor) {\n")
    outfile.write(f"\tv.Visit{class_name}(c)\n")
    outfile.write("}\n")


def write_preamble(outfile):
    outfile.write("package main\n")
    outfile.write(
        "\n// -- AUTOGENERATED FILE -- (see scripts/generate_ast.py for details...)\n")
    outfile.write(
        "// This is a simple implementation of the Visitor pattern from OOP\n")


if __name__ == '__main__':
    main()
