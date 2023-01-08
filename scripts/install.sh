#!/usr/bin/env bash

set -euo pipefail

# TODO: There is a lot of hardcoded stuff here (paths mostly)
# TODO: Split this into installation and setup because we want to suport package manager installation eventually
# TODO: "installation" should stay here and be simple, "setup" should be moved behind "reshctl setup"

echo
echo "Checking your system ..."

# /usr/bin/zsh -> zsh
login_shell=$(echo "$SHELL" | rev | cut -d'/' -f1 | rev)

if [ "$login_shell" != bash ] && [ "$login_shell" != zsh ]; then
    echo "ERROR: Unsupported/unknown login shell: $login_shell"
    exit 1
fi
echo " * Login shell: $login_shell - OK"

# TODO: Explicitly ask users if they want to enable RESH in shells
#       Only offer shells with supported versions
#       E.g. Enable RESH in: Zsh (your login shell), Bash, Both shells

# check like we are not running bash
bash_version=$(bash -c 'echo ${BASH_VERSION}')
bash_version_major=$(bash -c 'echo ${BASH_VERSINFO[0]}')
bash_version_minor=$(bash -c 'echo ${BASH_VERSINFO[1]}')
bash_too_old=""
if [ "$bash_version_major" -le 3 ]; then 
    bash_too_old=true
elif [ "$bash_version_major" -eq 4 ] && [ "$bash_version_minor" -lt 3 ]; then 
    bash_too_old=true
fi

if [ "$bash_too_old" = true ]; then
    echo " * Bash version: $bash_version - WARNING!"
    if [ "$login_shell" = bash ]; then
        echo "   > Your bash version is old."
        echo "   > Bash is also your login shell."
        echo "   > Updating to bash 4.3+ is STRONGLY RECOMMENDED!"
    else
        echo "   > Your bash version is old"
        echo "   > Bash is not your login shell so it should not be an issue."
        echo "   > Updating to bash 4.3+ is recommended."
    fi
else
    echo " * Bash version: $bash_version - OK"
fi


if ! zsh --version >/dev/null 2>&1; then
    echo " * Zsh version: ? - not installed!"
else
    zsh_version=$(zsh -c 'echo ${ZSH_VERSION}')
    zsh_version_major=$(echo "$zsh_version" | cut -d'.' -f1)
    if [ "$zsh_version_major" -lt 5 ]; then 
        echo " * Zsh version: $zsh_version - UNSUPPORTED!"
        if [ "$login_shell" = zsh ]; then
            echo "   > Your zsh version is old."
            echo "   > Zsh is also your login shell."
            echo "   > Updating to Zsh 5.0+ is STRONGLY RECOMMENDED!"
        else
            echo "   > Your zsh version is old"
            echo "   > Zsh is not your login shell so it should not be an issue."
            echo "   > Updating to zsh 5.0+ is recommended."
        fi
    else
        echo " * Zsh version: $zsh_version - OK"
    fi
fi


if [ "$(uname)" = Darwin ] && gnohup --version >/dev/null 2>&1; then
    echo " * Nohup installed: OK"
elif nohup --version >/dev/null 2>&1; then
    echo " * Nohup installed: OK"
else
    echo " * Nohup installed: NOT INSTALLED!"
    echo "   > You don't have nohup"
    echo "   > Please install GNU coreutils"
    echo
    if [ "$(uname)" = Darwin ]; then
        echo "   $ brew install coreutils"
        echo
    fi
    exit 1
fi

# echo 
# echo "Continue with installation? (Any key to CONTINUE / Ctrl+C to ABORT)"
# # shellcheck disable=2034
# read -r x

echo "Creating directories ..."

mkdir_if_not_exists() {
    if [ ! -d "$1" ]; then
        mkdir "$1" 
    fi
}

mkdir_if_not_exists ~/.resh
mkdir_if_not_exists ~/.resh/bin
mkdir_if_not_exists ~/.config

echo "Copying files ..."
cp -f submodules/bash-preexec/bash-preexec.sh ~/.bash-preexec.sh
cp -f submodules/bash-zsh-compat-widgets/bindfunc.sh ~/.resh/bindfunc.sh

cp -f scripts/shellrc.sh ~/.resh/shellrc
cp -f scripts/resh-daemon-start.sh ~/.resh/bin/resh-daemon-start
cp -f scripts/reshctl.sh scripts/widgets.sh scripts/hooks.sh ~/.resh/
cp -f scripts/rawinstall.sh ~/.resh/

# Copy all executables. We don't really need to omit install-utils from the bin directory
echo "Copying more files ..."
cp -f bin/resh-* ~/.resh/bin/
# rename reshctl
mv ~/.resh/bin/resh-control ~/.resh/bin/reshctl

# Shutting down resh daemon ...
echo "Shutting down resh daemon ..."
pid_file="${XDG_DATA_HOME-~/.local/share}/resh/daemon.pid"
if [ ! -f "$pid_file" ]; then
    # old pid file location
    pid_file=~/.resh/resh.pid
fi

failed_to_kill() {
    echo "ERROR: Failed to kill the resh-daemon - maybe it wasn't running?"
}

if [ -f "$pid_file" ]; then
    kill -SIGKILL "$pid_file" || failed_to_kill
    rm "$pid_file"
else
    pkill -SIGKILL "resh-daemon" || failed_to_kill
fi

echo "Creating/updating config file ..."
./bin/resh-install-utils migrate-config

echo "Checking/setting up device files ..."
./bin/resh-install-utils setup-device

echo "Updating format of history file ..."
./bin/resh-install-utils migrate-history

echo "Finishing up ..."
# Adding resh shellrc to .bashrc ...
if [ ! -f ~/.bashrc ]; then
    touch ~/.bashrc
fi
grep -q '[[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc' ~/.bashrc ||\
	echo -e '\n[[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc # this line was added by RESH (Rich Enhanced Shell History)' >> ~/.bashrc
# Adding bash-preexec to .bashrc ...
grep -q '[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh' ~/.bashrc ||\
	echo -e '\n[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh # this line was added by RESH (Rich Enhanced Shell History)' >> ~/.bashrc
# Adding resh shellrc to .zshrc ...
if [ -f ~/.zshrc ]; then
    grep -q '[ -f ~/.resh/shellrc ] && source ~/.resh/shellrc' ~/.zshrc ||\
        echo -e '\n[ -f ~/.resh/shellrc ] && source ~/.resh/shellrc # this line was added by RESH (Rich Enhanced Shell History)' >> ~/.zshrc
fi

# Deleting zsh completion cache - for future use
# [ ! -e ~/.zcompdump ] || rm ~/.zcompdump

echo "Launching resh daemon ..."
~/.resh/bin/resh-daemon-start

info="---- Scroll down using arrow keys ----
#####################################
#      ____  _____ ____  _   _      #
#     |  _ \| ____/ ___|| | | |     #
#     | |_) |  _| \___ \| |_| |     #
#     |  _ <| |___ ___) |  _  |     #
#     |_| \_\_____|____/|_| |_|     #
#    Rich Enhanced Shell History    #
#####################################
"

info="$info
RESH SEARCH APPLICATION = Redesigned reverse search that actually works

    >>> Launch RESH SEARCH app by pressing CTRL+R <<<
    (you will need to restart your terminal first)
     
    Search your history by commands. 
    Device, directories, git remote, and exit status is used to display relevant results first.

    At first, the search application will use the standard shell history without context. 
    All history recorded from now on will have context which will be used by the RESH SEARCH app.

CHECK FOR UPDATES
    To check for (and install) updates use reshctl command:
     $ reshctl update

HISTORY
    Your resh history will be recorded to '${XDG_DATA_HOME-~/.local/share}/resh/history.reshjson'
    Look at it using e.g. following command (you might need to install jq)
     $ cat ${XDG_DATA_HOME-~/.local/share}/resh/history.reshjson | sed 's/^v[^{]*{/{/' | jq .

ISSUES & FEEDBACK
    Please report issues to: https://github.com/curusarn/resh/issues
    Feedback and suggestions are very welcome!
"
if [ -z "${__RESH_VERSION:-}" ]; then info="$info
##############################################################
#                                                            #
#    Finish the installation by RESTARTING this terminal!    #
#                                                            #
##############################################################"
fi
info="$info
---- Close this by pressing Q ----" 

printf "%s\n" "$info" | ${PAGER:-less}

echo
echo "All done!"
echo "Thank you for using RESH"
echo "Report issues here: https://github.com/curusarn/resh/issues"
echo "Ctrl+R launches the RESH SEARCH app"

if [ -z "${__RESH_VERSION:-}" ]; then printf "
##############################################################
#                                                            #
#    Finish the installation by RESTARTING this terminal!    #
#                                                            #
##############################################################\n"
fi