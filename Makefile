build:
	go build .

local-install:
	go install ./

install:
	go install github.com/jorgejr568/cloudflare-cli@latest

test:
	go test -v ./...

install-tools:
	go install go.uber.org/mock/mockgen@latest

generate:
	go generate ./...