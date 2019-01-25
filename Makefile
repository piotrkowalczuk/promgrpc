all: setup get test

setup:
	ln -s -f ../../.githooks/pre-commit.sh .git/hooks/pre-commit

test: lint
	GO111MODULE=on go test -race -cover -coverprofile=cover-v2.out -count=5
	cd ./v3/ && GO111MODULE=on go test -race -cover -coverprofile=cover-v3.out -count=5
	cd ./v4/ && GO111MODULE=on go test -race -cover -coverprofile=cover-v4.out -count=5

lint:
	GO111MODULE=on gofmt -s -l .
	GO111MODULE=on goimports -l .
	GO111MODULE=on go vet .
	cd ./v3/ && GO111MODULE=on gofmt -s -l .
	cd ./v3/ && GO111MODULE=on goimports -l .
	cd ./v3/ && GO111MODULE=on go vet .
	cd ./v4/ && GO111MODULE=on gofmt -s -l .
	cd ./v4/ && GO111MODULE=on goimports -l .
	cd ./v4/ && GO111MODULE=on go vet .

get:
	go get -u -t golang.org/x/tools/cmd/goimports/...
	go get -u github.com/golang/lint/golint
#	go get -u honnef.co/go/tools/...

gen:
	protoc --go_out plugins=grpc:${GOPATH}/src v4/pb/private/test/test.proto
	goimports -w ./v4/pb/