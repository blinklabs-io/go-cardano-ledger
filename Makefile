.PHONY: test go-mod-tidy

test:
	go test -v ./...

# This will update the go.mod file based on imports in the code
go-mod-tidy:
	go mod tidy
