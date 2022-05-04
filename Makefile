SHELL=/bin/bash
LATEST_TAG=$(shell git describe --tags)
COMMIT=$(shell [ -z "$(git status --untracked-files=no --porcelain)" ] && git rev-parse --short=12 HEAD || echo "no_commit")
VERSION="${LATEST_TAG}-DEV"
GOFLAGS=-ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.development=true"


build: submodules bin/resh-session-init bin/resh-collect bin/resh-postcollect bin/resh-daemon\
  bin/resh-control bin/resh-config bin/resh-cli bin/resh-config-setup

install: build
	scripts/install.sh

test:
	go test -v ./...
	go vet ./...
	scripts/test.sh

rebuild:
	make clean
	make build

clean:
	rm -f bin/resh-*

uninstall:
	# Uninstalling ...
	-rm -rf ~/.resh/

go_files = $(shell find -name '*.go')
bin/resh-%: $(go_files)
	grep $@ .goreleaser.yml -q # all build targets need to be included in .goreleaser.yml
	go build ${GOFLAGS} -o $@ cmd/$*/*.go

.PHONY: submodules build install rebuild uninstall clean test


submodules: | submodules/bash-preexec/bash-preexec.sh submodules/bash-zsh-compat-widgets/bindfunc.sh
	@# sets submodule.recurse to true if unset
	@# sets status.submoduleSummary to true if unset
	@git config --get submodule.recurse >/dev/null || git config --global submodule.recurse true
	@#git config --get status.submoduleSummary >/dev/null || git config --global status.submoduleSummary true
	@#git config --get diff.submodule >/dev/null || git config --global diff.submodule log
	@# warns user if submodule.recurse is not set to true
	@[[ "true" == `git config --get submodule.recurse` ]] || echo "WARN: You should REALLY set 'git config --global submodule.recurse true'!"
	@#git config --global push.recurseSubmodules check

submodules/%:
	# Getting submodules ...
	git submodule sync --recursive 
	git submodule update --init --recursive

