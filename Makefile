deploy:
	@export GOPATH=${CURDIR}
	@mkdir -p /usr/local/go
	@export GOROOT=/usr/local/go
	@export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

.PHONY: deploy