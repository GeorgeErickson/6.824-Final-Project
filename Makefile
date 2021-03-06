

setup:
	rm -rf src/github.com
	mkdir -p /usr/local/go
	export GOPATH=${CURDIR}; \
	go get github.com/garyburd/go-websocket github.com/garyburd/redigo github.com/hoisie/web

watch:
	@cd src/server/client/; brunch w

run:
	export REDIS_TCP=127.0.0.1:6379; \
	export GOPATH=${CURDIR}; \
	go run src/server/main.go localhost:8000

run_external:
	export REDIS_TCP=pub-redis-11830.us-east-1-4.1.ec2.garantiadata.com:11830; \
	export REDIS_AUTH=davidgeorge; \
	export GOPATH=${CURDIR}; \
	go run src/server/main.go localhost:8000

pdev: watch run

dev:
	${MAKE} -j2 pdev


.PHONY: setup, run, run_external, watch, dev