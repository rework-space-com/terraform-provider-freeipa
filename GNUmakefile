default: install

install: 
	go build -o ~/go/bin/terraform-provider-freeipa
	
build: 
	go build -o $(shell pwd)/terraform-provider-freeipa

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

doc:
	go generate ./...

fmt:
	go fmt ./...

deps:
	go mod tidy

.PHONY: install build testacc doc fmt deps