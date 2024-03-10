package repl

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/samasno/little-compiler/pkg/ast"
	"github.com/samasno/little-compiler/pkg/eval"
	"github.com/samasno/little-compiler/pkg/lexer"
	"github.com/samasno/little-compiler/pkg/object"
)

// need to add error handling to lexer and parser

func Run() {
  println("Starting repl for little-compiler")
	scanner := bufio.NewScanner(os.Stdin)
  io.WriteString(os.Stdout, ">>")
  env := object.NewEnvironment()
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
				l := lexer.NewLexer(text)
        p := ast.New(l)
        prg := p.ParseProgram()
        o := eval.Eval(prg, env)
        io.WriteString(os.Stdout, o.Inspect())
        io.WriteString(os.Stdout, "\n>>")
				break inner
			}

		}
	}

	println("exiting repl")
}


