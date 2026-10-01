package main

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
)

const prompt = ">> "

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(prompt)

		if !scanner.Scan() {
			return
		}

		line := scanner.Text()

		if line == "" || line == "exit" {
			slog.Info("exiting from repl")
			break
		}

		fmt.Println("<<", line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error scanning lox: %v", err)
	}
}
