SHELL := /bin/bash

.PHONY: tools generate dev build package lint test

tools:
	go install github.com/wailsapp/wails/v3/cmd/wails@latest
	go install entgo.io/ent/cmd/ent@latest
	@echo "Tools installed."

generate:
	@echo "Generating ent code..."
	cd internal/store && ent generate ./ent/schema
	@echo "Done."

dev:
	wails dev

build:
	wails build

package:
	@echo "Packaging handled by Wails build outputs"

lint:
	golangci-lint run

test:
	go test ./...
