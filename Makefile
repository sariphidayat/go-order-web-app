.PHONY: build build-frontend build-backend run-backend run-frontend

GOCACHE ?= $(CURDIR)/backend/.gocache

build: build-frontend build-backend

build-frontend:
	cd frontend && npm run build

build-backend:
	mkdir -p bin backend/.gocache
	cd backend && GOCACHE=$(GOCACHE) go build -o ../bin/customer-order-app ./cmd/server

run-backend:
	mkdir -p backend/.gocache
	cd backend && GOCACHE=$(GOCACHE) go run ./cmd/server

run-frontend:
	cd frontend && npm run dev
