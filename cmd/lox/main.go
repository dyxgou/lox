package main

import (
	"bufio"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

const loxExtention = ".lox"

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("lox exec arguments expected path. got=%d args", len(os.Args))
	}

	path := os.Args[1]
	stat, err := os.Stat(path)

	if err != nil {
		log.Fatalf("invalid executable path: %s", err)
	}

	if ext := filepath.Ext(stat.Name()); ext != loxExtention {
		log.Fatalf("lox extension expected=%q. got=%q", loxExtention, ext)
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open lox file: %s", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		slog.Info("lox line", "l", line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}
}
