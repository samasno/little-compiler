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
)

// need to add error handling to lexer and parser

func Run() {
  println("Starting repl for little-compiler")
	scanner := bufio.NewScanner(os.Stdin)
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
        comp := compiler.New()

        err := comp.Compile(prg)
        if err != nil {
          fmt.Fprintf(os.Stdout, "Failed to compile: \n%s\n", err)
          continue
        }

        machine := vm.New(comp.Bytecode())
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


