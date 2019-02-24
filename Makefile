.PHONY:build$

serverConfig=config.json
installPath=/usr/local/bin

deploy: build src/server $(serverConfig)
	mv src/server $(installPath)/server
	cp $(serverConfig) $(installPath)/$(serverConfig)

build: server.go
	go build -o src/server server.go
