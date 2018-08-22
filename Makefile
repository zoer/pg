generate:
	go generate ./...
dev:
	go get github.com/golang/lint/golint
	go get honnef.co/go/tools/cmd/megacheck
	go get github.com/golang/dep/cmd/dep
	dep ensure

lint:
	@go vet -v $(go list ./... | grep -v /vendor/)
	@golint $(go list ./... | grep -v /vendor/)
	@megacheck $(go list ./... | grep -v /vendor/)

test:
	@go test -v -parallel 2 ./...
