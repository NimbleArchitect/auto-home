
export GO111MODULE=on

.PHONY: test
test:
	go test ./... -coverprofile cover.out

.PHONY: bin
bin: fmt vet test
	go build -o bin/server .

CLIENTSUBDIRS := $(wildcard clients/*/.)

clients: $(CLIENTSUBDIRS)
$(CLIENTSUBDIRS):
	$(MAKE) -C $@ bin

.PHONY: clients $(CLIENTSUBDIRS)

PLUGINSUBDIRS := $(wildcard plugins/*/.)

plugins: $(PLUGINSUBDIRS)
$(PLUGINSUBDIRS):
	$(MAKE) -C $@ bin

.PHONY: plugins $(PLUGINSUBDIRS)


.PHONY: all
all: bin clients plugins


.PHONY: fmt
fmt:
	go fmt .

.PHONY: vet
vet:
	go vet .

.PHONY: server-deps
server-deps:
	go get github.com/lucas-clemente/quic-go@v0.29.0
	go get github.com/dop251/goja@v0.0.0-20220915101355-d79e1b125a30

# .PHONY: setup
# setup:
# 	make -C setup