#!/bin/sh

git pull

echo go build .
go build .

echo running executable
./blogDownloadServer
