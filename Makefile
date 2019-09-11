SHELL=/bin/bash
VERSION=$(shell cat version)
REVISION=$(shell [ -z "$(git status --untracked-files=no --porcelain)" ] && git rev-parse --short=12 HEAD || echo "no_revision")
GOFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Revision=${REVISION}"

autoinstall: 
	./install_helper.sh

sanitize:
	#
	#
	# I'm going to create a sanitized version of your resh history.
	# Everything is done locally - your history won't leave this machine.
	# The way this works is that any sensitive information in your history is going to be replaced with its SHA1 hash.
	# There is also going to be a second version with hashes trimed to 12 characters for readability
	#
	#
	# > full hashes: ~/resh_history_sanitized.json
	# > 12 char hashes: ~/resh_history_sanitized_trim12.json
	#
	#
	# Encountered any issues? Got questions? -> Hit me up at https://github.com/curusarn/resh/issues
	#
	#
	# Running history sanitization ...
	resh-sanitize-history -trim-hashes 0 --output ~/resh_history_sanitized.json
	resh-sanitize-history -trim-hashes 12 --output ~/resh_history_sanitized_trim12.json
	# 
	# 
	# SUCCESS - ALL DONE!
	#
	# 
	# PLEASE HAVE A LOOK AT THE RESULT USING THESE COMMANDS:
	#
	# > pretty print JSON:
	@echo 'cat ~/resh_history_sanitized_trim12.json | jq'
	#
	# > only show executed commands, don't show metadata:
	@echo "cat ~/resh_history_sanitized_trim12.json | jq '.[\"cmdLine\"]'"
	#
	#
	#


build: submodules resh-collect resh-daemon resh-sanitize-history resh-evaluate

rebuild:
	make clean
	make build

clean:
	rm resh-*

install: build submodules/bash-preexec/bash-preexec.sh shellrc.sh config.toml uuid.sh | $(HOME)/.resh $(HOME)/.resh/bin $(HOME)/.config
	# Copying files to resh directory ...
	cp -f submodules/bash-preexec/bash-preexec.sh ~/.bash-preexec.sh
	cp -f config.toml ~/.config/resh.toml
	cp -f shellrc.sh ~/.resh/shellrc
	cp -f uuid.sh ~/.resh/bin/resh-uuid
	cp -f resh-* ~/.resh/bin/
	cp -f evaluate/resh-evaluate-plot.py ~/.resh/bin/
	cp -fr sanitizer_data ~/.resh/
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
	#
	#
	##########################################################
	#                                                        #
	#    SUCCESS - thank you for trying out this project!    #
	#                                                        #
	##########################################################
	#
	#
	# WHAT'S NEXT
	# Please RESTART ALL OPEN TERMINAL WINDOWS (or reload your rc files)
	# Your resh history is located in `~/.resh_history.json`
	# You can look at it using e.g. `tail -f ~/.resh_history.json | jq`
	#
	#
	# ISSUES
	# If anything looks broken create an issue: https://github.com/curusarn/resh/issues
	# You can uninstall this at any time by running `rm -rf ~/.resh/`
	# You won't lose any collected history by removing `~/.resh` directory
	#
	#
	# Please give me some contact info using this form: https://forms.gle/227SoyJ5c2iteKt98
	#
	#
	#

uninstall:
	# Uninstalling ...
	-rm -rf ~/.resh/

resh-daemon: daemon/resh-daemon.go common/resh-common.go version
	go build ${GOFLAGS} -o $@ $<

resh-collect: collect/resh-collect.go common/resh-common.go version
	go build ${GOFLAGS} -o $@ $<

resh-sanitize-history: sanitize-history/resh-sanitize-history.go common/resh-common.go version
	go build ${GOFLAGS} -o $@ $<

resh-evaluate: evaluate/resh-evaluate.go evaluate/strategy-*.go common/resh-common.go version
	go build ${GOFLAGS} -o $@ $< evaluate/strategy-*.go 

$(HOME)/.resh $(HOME)/.resh/bin $(HOME)/.config:
	# Creating dirs ...
	mkdir -p $@

$(HOME)/.resh/resh-uuid:
	# Generating random uuid for this device ...
	cat /proc/sys/kernel/random/uuid > $@ 2>/dev/null || ./uuid.sh 

.PHONY: submodules build install rebuild uninstall clean autoinstall


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

