#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/electron-download-linux-amd64 electron-download.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./bin/electron-download-linux-arm64 electron-download.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/electron-download-windows-amd64.exe electron-download.go
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o ./bin/electron-download-windows-arm64.exe electron-download.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go  build -o ./bin/electron-download-mac-amd64 electron-download.go
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go  build -o ./bin/electron-download-mac-arm64 electron-download.go