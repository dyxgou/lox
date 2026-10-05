package scanner

import (
	"log/slog"
	"testing"

	"github.com/dyxgou/lox/pkg/token"
)

func TestExample(t *testing.T) {
	src := []byte("cos(x) + 1i * sin(x)")

	var file token.File
	var scanner Scanner

	file.Init("euler", 1, len(src))
	scanner.Init(&file, src, nil)

	end := 0
	for scanner.ch != eof {
		pos, tok, lit := scanner.Scan()
		end = int(pos)

		slog.Info("euler", "pos", pos, "tok", tok, "lit", lit, "fPos", file.Position(pos))
	}

	if end != len(src) {
		t.Fatalf("ScanEuler source len expected=%d. got=%d", len(src), end)
	}
}
