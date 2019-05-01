
build: submodules/bash-preexec/bash-preexec.sh
	@echo "build"

install: build
	@echo "install"

submodules/%:
	-git submodule init
	git submodule update
