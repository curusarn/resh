#!/usr/bin/env bash

set -euo pipefail

while ! git --version &>/dev/null; do
    echo "Please install git."
    echo "Check again? (Any key to check again / Ctrl+C to exit)" 
    # shellcheck disable=2162 disable=2034
    read x
    echo
done

resh_git_dir=~/.resh_git
if [ ! -d "$resh_git_dir" ]; then 
    git clone https://github.com/curusarn/resh.git "$resh_git_dir"
    echo "Cloned https://github.com/curusarn/resh.git to $resh_git_dir"
    echo
fi

echo "Pulling the latest version of RESH ..."
cd "$resh_git_dir"
git checkout master
git pull
echo "Successfully pulled the latest version!"

make autoinstall

