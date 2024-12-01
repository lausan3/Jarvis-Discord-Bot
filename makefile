all: build run

build:
	go build -o bot main.go

run:
	./bot