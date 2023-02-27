.PHONY: default
default: clean build

.PHONY: clean
clean:
	- rm plugin
	- rm plugin.tar.gz

.PHONY: build
build:
	- GOOS=linux GOARCH=amd64 go build -o plugin plugin.go
	- tar -czvf plugin.tar.gz plugin plugin.json