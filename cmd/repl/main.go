package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/dyxgou/lox/pkg/scanner"
	"github.com/dyxgou/lox/pkg/token"
)

const prompt = ">> "

func main() {
	scn := bufio.NewScanner(os.Stdin)
	var file token.File
	var s scanner.Scanner

outter:
	for {
		fmt.Print(prompt)

		if !scn.Scan() {
			return
		}

		buf := scn.Bytes()

		file.Init("repl", 1, len(buf))
		s.Init(&file, buf, nil)

		for pos, tok, lit := s.Scan(); tok != token.EOF; pos, tok, lit = s.Scan() {
			if lit == "exit" {
				break outter
			}

			fmt.Printf("%q: %s: %d\n", tok.String(), lit, pos)
		}
	}

	if err := scn.Err(); err != nil {
		log.Fatalf("Error scanning lox: %v", err)
	}
}
