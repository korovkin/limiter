travis: build test
	go build

build:
	go build
	go build -o simple.exe ./examples/simple/*.go
	go build -o tui.exe ./examples/tui/*.go

run-simple:
	cd examples/simple && go run main.go

run-tui:
	cd examples/tui && go run main.go

test:
	go test -v .

tidy:
	go mod vendor
	go mod tidy
