#!/usr/bin/env bash

export GOPATH=/tmp/gopath
mkdir $GOPATH/bin
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
