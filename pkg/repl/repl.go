package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/samasno/little-compiler/pkg/lexer"
)

func Run() {
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
				ts := l.Tokenize()
				for _, t := range ts {
					fmt.Printf("%v\n", t)
				}
				break inner
			}

		}
	}
	println("exiting repl")
}

func prompt() {
	fmt.Printf(">>")
}
