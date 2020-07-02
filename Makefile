EXENAME=demo_server
ROOT_PATH=$(CURDIR)/
BIN=$(ROOT_PATH)bin

$(info ========================================)


all: build

format:
	gofmt -l -w -s ./

main:  
	go build -o $(BIN)/$(EXENAME) 
	chmod +x $(BIN)/$(EXENAME)

test:
	go test -gcflags='-N -l' -covermode=count -coverprofile=coverage.out -coverpkg ./... ./...
	@#workaround:https://github.com/golang/go/issues/22430
	@#sed -i "s/_${PWDSLASH}/./g" coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out -o coverage.txt
	@tail -n 1 coverage.txt | awk '{print $$1,$$3}'
	go test -c -o $(BIN)/$(EXENAME).test -covermode=count -coverpkg ./...

init:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.27.0

staticcheck:
	bin/golangci-lint run --no-config --issues-exit-code=0 --deadline=1m --disable-all --enable=govet --enable=structcheck --enable=lll --enable=staticcheck --enable=deadcode --enable=misspell --enable=nakedret  --enable=gocyclo --enable=varcheck --enable=structcheck --enable=errcheck --enable=ineffassign --enable=interfacer --tests=false ./...

build: format main test


clean:
	go clean
	rm -rf $(BIN)/$(EXENAME)
	rm -rf $(BIN)/$(EXENAME).test
	rm -f coverage.html
	rm -f coverage.out
	rm -f coverage.txt

.PHONY:main test clean format init staticcheck


