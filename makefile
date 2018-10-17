all: ensure_format lint transpile

no_dev_check: transpile

ensure_format:
	gofmt -s -w .

lint:
	golint ./... --vendor

transpile:
	gopherjs build stabell.go -o ./app.js
