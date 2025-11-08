default: testacc

# Run acceptance tests
.PHONY: testacc docs
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

docs:
	go generate ./...
