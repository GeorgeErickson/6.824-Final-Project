setup:
	@export GOPATH=${CURDIR}
	@mkdir -p /usr/local/go
	@export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
	@go get github.com/garyburd/go-websocket github.com/garyburd/redigo github.com/hoisie/web

.PHONY: setup