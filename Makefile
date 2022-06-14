default: lint test

lint:
	golangci-lint run

test:
	go test ./... -gcflags=all=-l

clean:
	rm -rf dist

proto:
	protoc --go_out=. pb/*.proto
	protoc --go_out=. pb_test/*.proto
	protoc-go-inject-tag -input=./pb_test/*.pb.go

build:
	go build -o dist/gpa main.go
