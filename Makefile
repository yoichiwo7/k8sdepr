NAME = k8sdepr
GOFLAGS = -ldflags "-s -w"

$(NAME): build

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	go build $(GOFLAGS) -o $(NAME)
