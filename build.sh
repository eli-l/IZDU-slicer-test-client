#!/bin/bash
GOOS=darwin GOARCH=arm64 go build -o ./.build/client-macos-arm64 main.go
GOOS=linux GOARCH=amd64 go build -o ./.build/client-linux-amd64 main.go
