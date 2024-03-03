package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/samasno/little-compiler/pkg/lexer"
	"github.com/samasno/little-compiler/pkg/ast"

)

// need to add error handling to lexer and parser

func Run() {
  println("Starting repl for little-compiler")
	scanner := bufio.NewScanner(os.Stdin)
outer:
	for {
		prompt()
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
        println(prg.String())
				break inner
			}

		}
	}
	println("exiting repl")
}

func prompt() {
	fmt.Printf(">>")
}
