BINNAME=fountain-linux-amd64
RELEASE=-s -w
UPXBIN=/usr/local/bin/upx
GOBIN=/usr/local/bin/go
GOOS=$(shell uname -s | tr [A-Z] [a-z])
GOARGS=GOARCH=amd64 CGO_ENABLED=0
GOBUILD=$(GOARGS) $(GOBIN) build -ldflags="$(RELEASE)"

.PHONY: all
all: clean build
clean:
	rm -f $(BINNAME)
	@echo Clean all.
build:
	@echo Compile $(BINNAME) ...
	GOOS=$(GOOS) $(GOBUILD) -o $(BINNAME) .
	@echo Build success.
upx: build
	$(UPXBIN) $(BINNAME)
upxx: build
	$(UPXBIN) --ultra-brute $(BINNAME)
vend:
	GOOS=$(GOOS) $(GOBUILD) -mod=vendor -o $(BINNAME) .
