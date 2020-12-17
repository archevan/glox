import sys
import os

# Abstract Syntax Tree Generator for Go
#
# Input: list of strings specifying the different
#        types of AST node types in the following
#        format:
#           "TYPE STRUCT NAME : FIELD-NAME TYPE, FIELD-NAME TYPE, ..."
# Output: an autogenerated .go file (called 'ast.go') in the current directory


def main():
    sys.argv.pop(0)  # remove first item (current exec path) from input list
    if len(sys.argv) != 1:
        print("usage: .\python.exe generate_ast.py [output_directory]")
        sys.exit(64)
    output_dir = sys.argv[0]
    output_file_name = "ast.go"
    current_path = os.path.dirname(__file__)  # get current directory
    output_path = os.path.relpath(
        output_dir, current_path) + "\\" + output_file_name
    print("generating AST struct types...")
    print(f"writing to file at {output_path}")
    ts = ["BinaryExpr   :  left Expr, op Token, right Expr",
          "Grouping     :  exp Expr",
          "Literal      :  val interface{}",
          "Unary        :  op Token, right Expr"]
    with open(output_path, 'w+') as f:
        define_ast(f, ts)


def define_ast(outfile, typestrings):
    write_preamble(outfile)
    outfile.write(
        "\n// simple empty struct as stand-in for an abstract class.\n// let me know if there's a better way (or submit a PR!).\n")
    outfile.write("type Expr struct{}\n")
    for typ in typestrings:
        class_name = typ.split(':')[0].strip()
        field_list = typ.split(':')[1].strip().split(',')
        write_class(outfile, class_name, field_list)


# write a Go struct type out to a given file
def write_class(outfile, class_name, field_list):
    outfile.write("\ntype " + class_name + " struct {")
    outfile.write("\n\tExpr")
    for field in field_list:
        outfile.write(f"\n\t{field.lstrip()}")
    outfile.write("\n}\n")


def write_preamble(outfile):
    outfile.write("package main\n")
    outfile.write(
        "\n// -- AUTOGENERATED FILE -- (see scripts/generate_ast.py for details...)\n")


if __name__ == '__main__':
    main()
