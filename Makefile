.PHONY: build run test

GOCACHE ?= $(CURDIR)/backend/.gocache

build:
	mkdir -p bin backend/.gocache
	cd backend && GOCACHE=$(GOCACHE) go build -o ../bin/customer-order-app ./cmd/server

run:
	mkdir -p backend/.gocache
	cd backend && GOCACHE=$(GOCACHE) go run ./cmd/server

test:
	mkdir -p backend/.gocache
	cd backend && GOCACHE=$(GOCACHE) go test ./...
