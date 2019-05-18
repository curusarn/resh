SHELL=/bin/bash


build: submodules resh-collect resh-daemon


install: build | $(HOME)/.resh $(HOME)/.resh/bin $(HOME)/.config/resh
	cp submodules/bash-preexec/bash-preexec.sh ~/.bash-preexec.sh -f
	cp config.toml ~/.config/resh.toml -f
	cp bashrc.sh ~/.resh/bashrc -f
	cp resh-* ~/.resh/bin/ -f
	grep '[[ -f ~/.resh/bashrc ]] && source ~/.resh/bashrc' ~/.bashrc ||\
		echo '[[ -f ~/.resh/bashrc ]] && source ~/.resh/bashrc' >> ~/.bashrc
	grep '[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh' ~/.bashrc ||\
		echo '[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh' >> ~/.bashrc
	-pkill resh-daemon
	nohup resh-daemon &>/dev/null & disown

resh-daemon: daemon/resh-daemon.go common/resh-common.go
	go build -o $@ $<

resh-collect: collect/resh-collect.go common/resh-common.go
	go build -o $@ $<


$(HOME)/.resh $(HOME)/.resh/bin $(HOME)/.config/resh:
	mkdir -p $@

.PHONY: submodules build install


submodules: | submodules/bash-preexec/bash-preexec.sh
	@# sets submodule.recurse to true if unset
	@# sets status.submoduleSummary to true if unset
	@git config --get submodule.recurse >/dev/null || git config --global submodule.recurse true
	@git config --get status.submoduleSummary >/dev/null || git config --global status.submoduleSummary true
	@git config --get diff.submodule >/dev/null || git config --global diff.submodule log
	@# warns user if submodule.recurse is not set to true
	@[[ "true" == `git config --get submodule.recurse` ]] || echo "WARN: You should REALLY set 'git config --global submodule.recurse true'!"
	@#git config --global push.recurseSubmodules check

submodules/%:
	git submodule sync --recursive 
	git submodule update --init --recursive

