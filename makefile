.PHONY: all fmt test testv tags doc vet

all: build

build:
	@ GOPATH=`pwd` go build -v ./src/...
	@ GOPATH=`pwd` go install ./src/...

fmt: 
	@ GOPATH=`pwd` go fmt ./src/...

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
	@ gotags `find src -name *.go` > tags

doc:
	@ GOPATH=`pwd` godoc -http=:8000
