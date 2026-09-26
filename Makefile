run: build_repl
	@ ./bin/repl

exec: build_lox
	@ ./bin/lox $(FILE)

build_repl:
	@ go build -o ./bin/repl ./cmd/repl/main.go

build_lox:
	@ go build -o ./bin/lox ./cmd/lox/main.go


