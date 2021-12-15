default: testacc
TF_ACC=1

# Run acceptance tests
.PHONY: testacc
testacc:
	go test ./... -v $(TESTARGS) -timeout 120m
