#!/usr/bin/env bash

set -euo pipefail

tmpdir="$(mktemp -d /tmp/resh-XXX)"
cd "$tmpdir"
git clone https://github.com/curusarn/resh.git
cd resh
make autoinstall
