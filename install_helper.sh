#!/usr/bin/env bash

if ! go version &>/dev/null; then
    echo
    echo "==========================================================================="
    echo
    echo "Please INSTALL GOLANG and run this again"
    echo
    if [ "$(uname)" = "Darwin" ]; then
        echo 'You can probably use `brew install go`.'
        echo
        echo "==========================================================================="
        echo
        exit 1
    elif [ "$(uname)" = "Linux" ]; then
        . /etc/os-release
        if [ "${ID}" = "ubuntu" ]; then
            echo 'You can probably use `sudo snap install go --classic` (gets latest golang - RECOMMENDED)'
            echo 'OR `sudo apt install go` (this might give you old golang)' 
            echo
            echo "==========================================================================="
            echo
            exit 1
        elif [ "${ID_LIKE}" = "debian" ]; then
            echo 'You can probably use `sudo apt install go`' 
            echo
            echo "==========================================================================="
            echo
            exit 1
        fi
    fi
    echo "It's recomended to use your favourite package manager."
    echo
    echo "==========================================================================="
    echo
    exit 1 
fi

go_version=$(go version | cut -d' ' -f3)
go_version_major=$(echo "${go_version:2}" | cut -d'.' -f1)
go_version_minor=$(echo "${go_version:2}" | cut -d'.' -f2)

if [ "$go_version_major" -gt 1 ]; then
    # good to go - future proof ;)
    echo "Building & installing ..."
    make install
elif [ "$go_version_major" -eq 1 ] && [ "$go_version_minor" -ge 11 ]; then
    # good to go - we have go modules
    echo "Building & installing ..."
    make install
else
    echo
    echo "==========================================================================="
    echo "Your Golang version is older than 1.11 - we can't use go modules for build!"
    echo "I will try to build the project using dep. (I will let you review each step.)"
    echo "Continue? (Any key to continue / Ctrl+C to cancel)" 
    read x

    take_care_of_gopath=0
    if [ -z "${GOPATH+x}" ]; then
        echo
        echo "==========================================================================="
        echo "GOPATH env variable is unset!"
        echo "I will take care of GOPATH. (I will create tmp GOPATH.)"
        echo "Continue? (Any key to continue / Ctrl+C to cancel)" 
        read x

        GOPATH=$(mktemp -d /tmp/gopath-XXX) \
            && mkdir "$GOPATH/bin" \
            && echo "Created tmp GOPATH: $GOPATH"
        export GOPATH
        take_care_of_gopath=1
    fi

    echo "GOPATH=$GOPATH"
    PATH=$GOPATH/bin:$PATH

    if ! dep version &>/dev/null; then
        echo
        echo "==========================================================================="
        echo "It appears that you don't have dep installed!"
        echo "I will install dep. (I will install it from GitHub.)"
        echo "Continue? (Any key to continue / Ctrl+C to cancel)" 
        read x
        
        curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
        [ $? -eq 0 ] && echo "Installed dep."
    fi

    project_path=$GOPATH/src/github.com/curusarn/resh
    mkdir -p "$project_path" &>/dev/null
    if [ "$project_path" != "$PWD" ]; then
        if [ "$take_care_of_gopath" -eq 0 ]; then
            echo
            echo "==========================================================================="
            echo "It seems that current directory is not in the GOPATH!"
            echo "I will copy the project to appropriate GOPATH directory."
            echo "Continue? (Any key to continue / Ctrl+C to cancel)" 
            read x
        fi
        cp -rf ./* .git* "$project_path" && echo "Copied files to $project_path"
        cd "$project_path"
    fi

    echo "Running \`dep ensure\` ..."
    if ! dep ensure; then
        echo "Unexpected ERROR while running \`dep ensure\`!"
        exit 2
    fi
    echo
    echo "==========================================================================="
    echo "Building & installing ..."
    make install
fi
