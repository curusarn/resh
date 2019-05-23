#!/usr/bin/env bash

if ! go version 2>/dev/null; then
    echo "Please install Golang and rerun this script"
    exit 1 
fi

go_version=$(go version | cut -d' ' -f3)
go_version_major=$(echo ${go_version:2} | cut -d'.' -f1)
go_version_minor=$(echo ${go_version:2} | cut -d'.' -f2)

if [ "$go_version_major" -gt 1 ]; then
    # good to go - future proof ;)
    echo "Building & installing ..."
    make install
elif [ "$go_version_major" -eq 1 ] && [ "$go_version_minor" -ge 11 ]; then
    # good to go - we have go modules
    echo "Building & installing ..."
    make install
elif dep version >/dev/null; then
    if ! dep init; then
        echo "`dep init` failed - bootstraping GOPATH ..."
        export GOPATH=/tmp/gopath
        project_path=$GOPATH/src/github.com/curusarn/resh
        mkdir -p $project_path
        cp -rf . $project_path
        cd $project_path
        if ! dep init; then
            echo "Unexpected ERROR while running `dep init`!"
            echo "Try fixing it yourself - I'm sorry :("
            exit 3
        fi
        echo "Succesfuly bootstraped GOPATH"
    fi
    echo "Building & installing ..."
    make install
else
    echo "Please install dep and rerun this script"
    exit 2
fi
