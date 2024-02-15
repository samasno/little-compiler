package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Run() {
	scanner := bufio.NewScanner(os.Stdin)
	prompt()
loop:
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.Replace(text, "\n", "", -1)
		switch text {
		case "quit":
			println("received quit command")
			break loop
		default:
			println(text)
		}
		prompt()
	}
	println("exiting repl")
}

func prompt() {
	fmt.Printf(">>")
}
