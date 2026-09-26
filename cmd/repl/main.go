package main

import (
	"bufio"
	"fmt"
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
		fmt.Println("<<", line)
	}
}
