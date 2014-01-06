.PHONY: all fmt test testv tags doc vet lc

all: build

build:
	@ GOPATH=`pwd` go install ./src/...

fmt: 
	@ GOPATH=`pwd` gofmt -s -w -l src

vet: 
	@ GOPATH=`pwd` go vet ./src/...

testv:
	@ GOPATH=`pwd` go test -v ./src/...

test:
	@ GOPATH=`pwd` go test ./src/...

clean:
	@ rm -rf pkg bin

fix:
	@ GOPATH=`pwd` go fix ./src/...

tags:
	@ gotags `find src -name "*.go"` > tags

doc:
	@ GOPATH=`pwd` godoc -http=:8000

lc:
	@ wc -l `find src -name "*.go" | grep -v regmap.go`
