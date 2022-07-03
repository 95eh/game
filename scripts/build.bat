#!/bin/sh
set -v on
export CGO_ENABLED=0
export GOOS=windows
export GOARCH=amd64
go build -o ../output/game.exe ../cmd/game.go