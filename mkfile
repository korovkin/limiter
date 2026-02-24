travis:
	go build

build:
	go build
	cd examples/simple && go build -o simple
	cd examples/tui && go build -o tui

run-simple:
	cd examples/simple && go run main.go

run-tui:
	cd examples/tui && go run main.go

test:
	go test -v .
