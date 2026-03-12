# Makefile for Taxi-app

.PHONY: lint format

lint:
	golangci-lint run ./backend/...

format:
	gofmt -w ./backend 