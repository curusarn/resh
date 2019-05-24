SHELL=/bin/bash


autoinstall: 
	./install_helper.sh


build: submodules resh-collect resh-daemon


install: build | $(HOME)/.resh $(HOME)/.resh/bin $(HOME)/.config $(HOME)/.resh/resh-uuid
	# Copying files to resh directory ...
	cp -f submodules/bash-preexec/bash-preexec.sh ~/.bash-preexec.sh
	cp -f config.toml ~/.config/resh.toml
	cp -f shellrc.sh ~/.resh/shellrc
	cp -f uuid.sh ~/.resh/bin/resh-uuid
	cp -f resh-* ~/.resh/bin/
	# backward compatibility: We have a new location for resh history file 
	[ ! -f ~/.resh/history.json ] || mv ~/.resh/history.json ~/.resh_history.json 
	# Adding resh shellrc to .bashrc ...
	grep '[[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc' ~/.bashrc ||\
		echo '[[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc' >> ~/.bashrc
	# Adding bash-preexec to .bashrc ...
	grep '[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh' ~/.bashrc ||\
		echo '[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh' >> ~/.bashrc
	# Adding resh shellrc to .zshrc ...
	grep '[ -f ~/.resh/shellrc ] && source ~/.resh/shellrc' ~/.zshrc ||\
		echo '[ -f ~/.resh/shellrc ] && source ~/.resh/shellrc' >> ~/.zshrc
	# Restarting resh daemon ...
	[ ! -f ~/.resh/resh.pid ] || kill -SIGTERM $$(cat ~/.resh/resh.pid)
	nohup resh-daemon &>/dev/null & disown
	# Final touch
	touch ~/.resh_history.json
	#
	##########################################################
	#                                                        #
	#    SUCCESS - thank you for trying out this project!    #
	#                                                        #
	##########################################################
	#
	# WHAT'S NEXT
	# Please close all open terminal windows (or reload your rc files)
	# Your resh history is located in `~/.resh_history.json`
	# You can look at it using e.g. `tail -f ~/.resh_history.json | jq`
	#
	# ISSUES
	# If anything looks broken create an issue: https://github.com/curusarn/resh/issues
	# You can uninstall this at any time by running `rm -rf ~/.resh/`
	# You won't lose any collected history by removing `~/.resh` directory
	#

uninstall:
	# Uninstalling ...
	-rm -rf ~/.resh/

resh-daemon: daemon/resh-daemon.go common/resh-common.go
	go build -o $@ $<

resh-collect: collect/resh-collect.go common/resh-common.go
	go build -o $@ $<


$(HOME)/.resh $(HOME)/.resh/bin $(HOME)/.config:
	# Creating dirs ...
	mkdir -p $@

$(HOME)/.resh/resh-uuid:
	# Generating random uuid for this device ...
	cat /proc/sys/kernel/random/uuid > $@ 2>/dev/null || ./uuid.sh 

.PHONY: submodules build install


submodules: | submodules/bash-preexec/bash-preexec.sh
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

