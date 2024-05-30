package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/samasno/little-compiler/pkg/compiler"
	"github.com/samasno/little-compiler/pkg/frontend/lexer"
	"github.com/samasno/little-compiler/pkg/frontend/parser"
	"github.com/samasno/little-compiler/pkg/vm"
  "github.com/samasno/little-compiler/pkg/frontend/object"
)

// need to add error handling to lexer and parser

func Run() {
  println("Starting repl for little-compiler")
	scanner := bufio.NewScanner(os.Stdin)
  constants := []object.Object{}
  symbolTable := compiler.NewSymbolTable()
  globals := make([]object.Object, vm.GlobalSize)
  io.WriteString(os.Stdout, ">>")
outer:
	for {
	inner:
		for scanner.Scan() {
			text := scanner.Text()
			text = strings.Replace(text, "\n", "", -1)
			switch text {
			case "quit":
				println("received quit command")
				break outer
			default:
				l := lexer.New(text)
        p := parser.New(l)
        prg := p.ParseProgram()
        comp := compiler.NewWithState(symbolTable, constants)
        err := comp.Compile(prg)
        if err != nil {
          fmt.Fprintf(os.Stdout, "Failed to compile: \n%s\n", err)
          continue
        }

        code := comp.Bytecode()
        constants = code.Constants

        machine := vm.NewWithGlobalStore(code, globals)
        machine.Run()
        o := machine.LastPoppedStackElement()
        io.WriteString(os.Stdout, o.Inspect())
        io.WriteString(os.Stdout, "\n>>")
				break inner
			}

		}
	}

	println("exiting repl")
}


