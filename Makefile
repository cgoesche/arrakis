.DEFAULT_GOAL := build
INSTALL_DIR   ::= /usr/local/bin

.PHONY:fmt vet build
fmt: 
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o ./build/arrakis

install: build
	sudo cp ./build/arrakis $(INSTALL_DIR)/

clean: 
	go clean
	rm ./build/*
