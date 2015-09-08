@echo off
env GOARCH=amd64 GOOS=windows go build -o connection_test_amd64.exe
env GOARCH=386 GOOS=windows go build -o connection_test_i386.exe
env GOARCH=amd64 GOOS=darwin go build -o connection_test_amd64darwin
env GOARCH=386 GOOS=darwin go build -o connection_test_i386darwin
