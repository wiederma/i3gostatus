#!/bin/sh

if [ "$1" = "-docker" ]; then
    docker build -t i3gostatus-build .
    docker create --name i3gostatus-build-cont i3gostatus-build
    docker cp "i3gostatus-build-cont:/var/build/i3gostatus" .
    docker rm i3gostatus-build-cont
fi

go build -o i3gostatus cmd/main.go
