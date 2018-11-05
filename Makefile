.PHONY: build$

build: server.go
	go build -o server server.go

deploy: build:
	mv server /usr/local/bin/server
