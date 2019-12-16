#!/usr/bin/env bash

set -euo pipefail

echo
echo "Please report any issues you encounter to: https://github.com/curusarn/resh/issues"
echo

echo "Looking for the latest release ..."
json=$(curl --silent "https://api.github.com/repos/curusarn/resh/releases/latest")

# not very robust but we don't want any dependencies to parse to JSON
tag=$(echo "$json" | grep '"tag_name":' | cut -d':' -f2 | tr -d ',' | cut -d'"' -f2)

if [ ${#tag} -lt 2 ]; then
    echo "ERROR: Couldn't determine the latest release! (extracted git tag is too short \"${tag}\")"
    exit 1
fi
if [ "${tag:0:1}" != v ]; then
    echo "ERROR: Couldn't determine the latest release! (extracted git tag doesn't start with 'v' \"${tag}\")"
    exit 1
fi
version="${tag:1}"
# TODO: check if version is numeric

echo " * Latest version: $version (git tag: $tag)"
echo

if [ "$(uname)" = "Darwin" ]; then
    OS=darwin
elif [ "$(uname)" = "Linux" ]; then
    OS=linux
else
    OS=unknown
fi

if [ "$(uname -m)" = "x86_64" ]; then
    ARCH=amd64
elif [ "$(uname -m)" = "i386" ]; then
    ARCH=386
else
    ARCH=unknown
fi


if [ "$OS" = unknown ] || [ "$ARCH" = unknown ]; then
    echo "Couldn't detect your OS and architecture - exiting!"
    echo "Expected Linux or macOS on x86_64 or i386"
    exit 1
fi

dl_base="https://github.com/curusarn/resh/releases/download/${tag}"

fname_checksums="resh_${version}_checksums.txt"
dl_checksums="$dl_base/$fname_checksums"

fname_binaries="resh_${version}_${OS}_${ARCH}.tar.gz"
dl_binaries="$dl_base/$fname_binaries"


tmpdir="$(mktemp -d /tmp/resh-rawinstall-XXX)"
# echo
# echo "Changing to $tmpdir ..."
cd "$tmpdir"

echo "Downloading files ..."

curl_opt="--location --remote-name --progress-bar"

echo " * $fname_checksums"
# shellcheck disable=2086
COLUMNS=80 curl $curl_opt "$dl_checksums"

echo " * $fname_binaries"
# shellcheck disable=2086
COLUMNS=80 curl $curl_opt "$dl_binaries"

# TODO: check if we downloaded anything
# Github serves you a "Not found" page so the curl doesn't error out

echo
echo "Checking integrity ..."

# macOS doesn't have sha256sum
if [ "$OS" = darwin ]; then
    function sha256sum() { shasum -a 256 "$@" ; } && export -f sha256sum
fi

if [ "$(sha256sum "$fname_binaries")" != "$(grep "$fname_binaries" "$fname_checksums")" ]; then
    echo "ERROR: integrity check failed - exiting!"
    exit 1
fi
echo " * OK"

echo
echo "Extracting downloaded files ..."
tar -xzf "$fname_binaries"
echo " * OK"

if ! scripts/install.sh; then
    if [ $? != 130 ]; then
        echo
        echo "INSTALLATION FAILED!"
        echo "I'm sorry for the inconvenience."
        echo
        echo "Please create an issue: https://github.com/curusarn/resh/issues"
    fi
    echo
    echo "You can rerun the installation by executing: (this will skip downloading)"
    echo
    echo "cd $PWD && scripts/install.sh"
    exit 1
fi