.PHONY: build install doc fmt lint run

export GOPATH
export GO15VENDOREXPERIMENT=1

default: build

install: build
	sudo cp ./bin/bumblebee-ui /usr/bin/

build: vet vendors
	go build -v -o ./bin/bumblebee-ui ./src/main.go

vendors:
	glide install

doc:
	godoc -http=:6060 -index

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	go fmt ./src/...

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	golint ./src

run: build
	./bin/main_app

# http://godoc.org/code.google.com/p/go.tools/cmd/vet
# go get code.google.com/p/go.tools/cmd/vet
vet:
	go vet ./src/...