.PHONY:build$

deploy: build src/server config.json
	mv src/server /usr/local/bin/server
	cp config.json /usr/local/bin/config.json

build: server.go
	go build -o src/server server.go
