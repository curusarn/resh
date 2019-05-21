SHELL=/bin/bash


build: submodules resh-collect resh-daemon


install: build | $(HOME)/.resh $(HOME)/.resh/bin $(HOME)/.config $(HOME)/.resh/resh-uuid
	cp -f submodules/bash-preexec/bash-preexec.sh ~/.bash-preexec.sh
	cp -f config.toml ~/.config/resh.toml
	cp -f shellrc.sh ~/.resh/shellrc
	cp -f resh-* ~/.resh/bin/
	[ ! -f ~/resh-history.json ] || mv ~/resh-history.json ~/.resh/history.json 
	grep '[[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc' ~/.bashrc ||\
		echo '[[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc' >> ~/.bashrc
	grep '[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh' ~/.bashrc ||\
		echo '[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh' >> ~/.bashrc
	grep '[ -f ~/.resh/shellrc ] && source ~/.resh/shellrc' ~/.zshrc ||\
		echo '[ -f ~/.resh/shellrc ] && source ~/.resh/shellrc' >> ~/.zshrc
	[ ! -f ~/.resh/resh.pid ] || kill -SIGTERM $$(cat ~/.resh/resh.pid)
	nohup resh-daemon &>/dev/null & disown

uninstall:
	-mv ~/.resh/history.json ~/resh-history.json
	-rm -rf ~/.resh

resh-daemon: daemon/resh-daemon.go common/resh-common.go
	go build -o $@ $<

resh-collect: collect/resh-collect.go common/resh-common.go
	go build -o $@ $<


$(HOME)/.resh $(HOME)/.resh/bin $(HOME)/.config:
	mkdir -p $@

$(HOME)/.resh/resh-uuid:
	-cat /proc/sys/kernel/random/uuid > $@

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

