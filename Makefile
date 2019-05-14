SHELL=/bin/bash

build: submodules resh-collect resh-daemon

install: build $(HOME)/.resh $(HOME)/.resh/bin
	cp submodules/bash-preexec/bash-preexec.sh ~/.bash-preexec.sh -f
	cp bashrc.sh ~/.resh/bashrc -f
	cp resh-* ~/.resh/bin/ -f
	grep '[[ -f ~/.resh/bashrc ]] && source ~/.resh/bashrc' ~/.bashrc ||\
		echo '[[ -f ~/.resh/bashrc ]] && source ~/.resh/bashrc' >> ~/.bashrc
	grep '[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh' ~/.bashrc ||\
		echo '[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh' >> ~/.bashrc
#	-pkill resh-daemon
#	resh-daemon &


resh-daemon: daemon/resh-daemon.go common/resh-common.go
	go build -o $@ $<

resh-collect: collect/resh-collect.go common/resh-common.go
	go build -o $@ $<


$(HOME)/.resh:
	mkdir $(HOME)/.resh

$(HOME)/.resh/bin:
	mkdir $(HOME)/.resh/bin

.PHONY: submodules build install

submodules:
	# update (and intialize) submodules (recursively)
	git submodule update --init --recursive

submodules_to_latest_commit:
	git submodule foreach --recursive git pull origin master
