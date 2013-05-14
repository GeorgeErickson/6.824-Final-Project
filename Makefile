setup:
	@export GOPATH=${CURDIR}
	@mkdir -p /usr/local/go
	@export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
	@go get github.com/garyburd/go-websocket github.com/garyburd/redigo github.com/hoisie/web

watch:
	@cd src/server/client/; brunch w

run:
	@export REDIS_TCP=127.0.0.1:6379
	@go run src/server/main.go localhost:8000

run_external:
	@export REDIS_TCP=pub-redis-11830.us-east-1-4.1.ec2.garantiadata.com:11830
	@export REDIS_AUTH=davidgeorge
	@go run src/server/main.go localhost:8000

.PHONY: setup, run, run_external, watch