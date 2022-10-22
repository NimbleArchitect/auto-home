
export GO111MODULE=on

.PHONY: test
test:
	go test . -coverprofile cover.out

.PHONY: bin
bin: fmt vet
	go build -o ../../bin/plugins/$(notdir $(CURDIR))/$(notdir $(CURDIR)) .


.PHONY: fmt
fmt:
	go fmt .

.PHONY: vet
vet:
	go vet .
