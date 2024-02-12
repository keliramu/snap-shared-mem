#!/bin/bash

rm ./client ./server
rm ./*.snap

echo "Build server"
go build -o server cmd/server/main.go

echo "Build client"
go build -o client cmd/client/main.go

echo "Build snap"
time snapcraft 2>&1 | tee build-snap.log


echo "Done."
