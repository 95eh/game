#!/bin/sh
set -v on
#export CGO_ENABLED=0
#export GOOS=linux
export GOOS=windows
export GOARCH=amd64
#export GOARCH=arm64
#go build -o ../bin/game ../cmd/game.go
#go build -o ../bin/global ../cmd/global.go
#go build -o ../bin/test.exe ../cmd/test.go
#go build -o ../bin/test_server ../cmd/test_server.go
go build -o ../bin/test_client.exe ../cmd/test_client.go