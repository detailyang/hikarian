GOPATH := ${GOPATH}:$(shell pwd)
.PHONY: clean test

all:
		@GOPATH=$(GOPATH) go install tunnel
		@GOPATH=$(GOPATH) go install socks5
		@GOPATH=$(GOPATH) go install icmp

tunnel:
		@GOPATH=$(GOPATH) go install tunnel

socks5:
		@GOPATH=$(GOPATH) go install socks5

icmp:
		@GOPATH=$(GOPATH) go install icmp

clean:
		@rm -fr bin pkg

test:
		@GOPATH=$(GOPATH) go test tunnel
