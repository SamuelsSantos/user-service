.PHONY: all
all: build
FORCE: ;

BIN_DIR = $(PWD)/bin

.PHONY: build

clean:
	rm -rf bin/*

dependencies:
	go mod download

build: dependencies build-api linux-binaries

build-api: 
	go build -o ./bin/users-api api/main.go

linux-binaries:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o $(BIN_DIR)/users-api-linux api/main.go

	

ci: dependencies test	

build-mocks:
	@go get github.com/golang/mock/gomock
	@go install github.com/golang/mock/mockgen
	@~/go/bin/mockgen -source=domain/entity/user/interface.go -destination=domain/entity/user/mock/user.go -package=mock


test:
	mkdir -p ./coverage
	go test -tags testing ./... -v -coverprofile=./coverage/tests.out
	go tool cover -html=./coverage/tests.out -o ./coverage/coverage-report.html


fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done