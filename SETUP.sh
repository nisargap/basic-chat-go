#!/bin/sh

BUILD_PATH=$GOPATH/src/github.com/nisargap/basic-chat-go
cd $BUILD_PATH
go get -u github.com/kardianos/govendor
$GOBIN/govendor sync
echo "Brought in govender and installed dependencies via $GOBIN/govendor sync"


