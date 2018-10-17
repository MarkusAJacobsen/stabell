all: ensure_format lint transpile

ensure_format:
	gofmt -s -w .

lint:
	golint ./... --vendor

transpile:
	gopherjs build stabell.go -o ./app.js
