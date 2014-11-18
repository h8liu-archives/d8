.PHONY: all fmt test

all:
	go install ./...

fmt:
	gofmt -s -w -l .

test:
	go test ./...

lc:
	wc -l `find . -name "*.go" | grep -v regmap`
