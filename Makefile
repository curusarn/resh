SHELL=/bin/bash

build: submodules resh-collect

install: build $(HOME)/.resh $(HOME)/.resh/bin
	cp submodules/bash-preexec/bash-preexec.sh ~/.bash-preexec.sh
	cp src/bashrc.sh ~/.resh/bashrc
	cp resh-collect ~/.resh/bin/
	grep '[[ -f ~/.resh/bashrc ]] && source ~/.resh/bashrc' ~/.bashrc ||\
		echo '[[ -f ~/.resh/bashrc ]] && source ~/.resh/bashrc' >> ~/.bashrc
	grep '[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh' ~/.bashrc ||\
		echo '[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh' >> ~/.bashrc

resh-collect: src/resh-collect.go
	go build -o resh-collect src/resh-collect.go

$(HOME)/.resh:
	mkdir $(HOME)/.resh

$(HOME)/.resh/bin:
	mkdir $(HOME)/.resh/bin

.PHONY: submodules build install

submodules: submodules/bash-preexec/bash-preexec.sh
	# this is always run and updates submodules
	git submodule update --recursive

submodules/%:
	git submodule init
