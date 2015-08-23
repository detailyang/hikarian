GOPATH := $(shell pwd)
.PHONY: clean test

all:
		@GOPATH=$(GOPATH) go install tunnel
		@GOPATH=$(GOPATH) go install socks5
	
tunnel:
		@GOPATH=$(GOPATH) go install tunnel

socks5:
		@GOPATH=$(GOPATH) go install socks5

clean:
		@rm -fr bin pkg

test:
		@GOPATH=$(GOPATH) go test tunnel
