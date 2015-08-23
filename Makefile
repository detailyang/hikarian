GOPATH := $(shell pwd)
.PHONY: clean test

tunnel:
		@GOPATH=$(GOPATH) go install tunnel

socks5:
		@GOPATH=$(GOPATH) go install socks5

clean:
		@rm -fr bin pkg

test:
		@GOPATH=$(GOPATH) go test tunnel
