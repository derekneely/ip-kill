#!/bin//bash

# build linux
env GOOS=linux GOARCH=amd64 go build ip-kill.go
tar czvf ipkill