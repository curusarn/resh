#!/usr/bin/env sh

set -euo pipefail

# TODO: There is a lot of hardcoded stuff here (paths mostly)
# TODO: Split this into installation and setup because we want to suport package manager installation eventually
# TODO: "installation" should stay here and be simple, "setup" should be moved behind "reshctl setup"

echo
echo "Checking your system ..."
printf '\e[31;1m' # red color on

cleanup() {
    printf '\e[0m' # reset
    exit
}
trap cleanup EXIT INT TERM

# /usr/bin/zsh -> zsh
login_shell=$(echo "$SHELL" | rev | cut -d'/' -f1 | rev)

if [ "$login_shell" != bash ] && [ "$login_shell" != zsh ]; then
    echo "* UNSUPPORTED login shell: $login_shell"
    echo " -> RESH supports zsh and bash"
    echo
    if [ -z "${RESH_INSTALL_IGNORE_LOGIN_SHELL-}" ]; then
        echo 'EXITING!'
        echo ' -> You can skip this check with `export RESH_INSTALL_IGNORE_LOGIN_SHELL=1`'
        exit 1
    fi
fi

# TODO: Explicitly ask users if they want to enable RESH in shells
#       Only offer shells with supported versions
#       E.g. Enable RESH in: Zsh (your login shell), Bash, Both shells

bash_version=$(bash -c 'echo ${BASH_VERSION}')
bash_version_major=$(bash -c 'echo ${BASH_VERSINFO[0]}')
bash_version_minor=$(bash -c 'echo ${BASH_VERSINFO[1]}')
bash_ok=1
if [ "$bash_version_major" -le 3 ]; then
    bash_ok=0
elif [ "$bash_version_major" -eq 4 ] && [ "$bash_version_minor" -lt 3 ]; then
    bash_ok=0
fi

if [ "$bash_ok" != 1 ]; then
    echo "* UNSUPPORTED bash version: $bash_version"
    echo " -> Update to bash 4.3+ if you want to use RESH in bash"
    echo
fi

zsh_ok=1
if ! zsh --version >/dev/null 2>&1; then
    echo "* Zsh not installed"
    zsh_ok=0
else
    zsh_version=$(zsh -c 'echo ${ZSH_VERSION}')
    zsh_version_major=$(echo "$zsh_version" | cut -d'.' -f1)
    if [ "$zsh_version_major" -lt 5 ]; then
        echo "* UNSUPPORTED zsh version: $zsh_version"
        echo " -> Updatie to zsh 5.0+ if you want to use RESH in zsh"
        echo
        zsh_ok=0
    fi
fi

if [ "$bash_ok" != 1 ] && [ "$zsh_ok" != 1 ]; then
    echo "* You have no shell that is supported by RESH!"
    echo " -> Please install/update zsh or bash and run this installation again"
    echo
    if [ -z "${RESH_INSTALL_IGNORE_NO_SHELL-}" ]; then
        echo 'EXITING!'
        echo ' -> You can prevent this check by setting `export RESH_INSTALL_IGNORE_NO_SHELL=1`'
        echo
        exit 1
    fi
fi

printf '\e[0m' # reset
# echo "Continue with installation? (Any key to CONTINUE / Ctrl+C to ABORT)"
# # shellcheck disable=2034
# read -r x

if [ -z "${__RESH_VERSION-}" ]; then
    # First installation
    # Stop the daemon anyway just to be sure
    # But don't output anything
    ./scripts/resh-daemon-stop.sh -q
else
    ./scripts/resh-daemon-stop.sh
fi

echo "Installing ..."

# Crete dirs first to get rid of edge-cases
# If we fail we don't roll back - directories are harmless
mkdir_if_not_exists() {
    if [ ! -d "$1" ]; then
        mkdir "$1"
    fi
}

mkdir_if_not_exists ~/.resh
mkdir_if_not_exists ~/.resh/bin
mkdir_if_not_exists ~/.config

# Run setup and update tasks

./bin/resh-install-utils setup-device
# migrate-all updates format of config and history
# migrate-all restores original config and history on error
# There is no need to roll back anything else because we haven't replaced
#   anything in the previous installation.
./bin/resh-install-utils migrate-all


# Copy files

# echo "Copying files ..."
cp -f submodules/bash-preexec/bash-preexec.sh ~/.bash-preexec.sh
cp -f submodules/bash-zsh-compat-widgets/bindfunc.sh ~/.resh/bindfunc.sh

cp -f scripts/shellrc.sh ~/.resh/shellrc
cp -f scripts/resh-daemon-start.sh ~/.resh/bin/resh-daemon-start
cp -f scripts/resh-daemon-stop.sh ~/.resh/bin/resh-daemon-stop
cp -f scripts/resh-daemon-restart.sh ~/.resh/bin/resh-daemon-restart
cp -f scripts/hooks.sh ~/.resh/
cp -f scripts/rawinstall.sh ~/.resh/

# echo "Copying more files ..."
# Copy all go executables. We don't really need to omit install-utils from the bin directory
cp -f bin/resh-* ~/.resh/bin/
# Rename reshctl
mv ~/.resh/bin/resh-control ~/.resh/bin/reshctl


echo "Handling shell files ..."
# Only add shell directives into bash if it passed version checks
if [ "$bash_ok" = 1 ]; then
    if [ ! -f ~/.bashrc ]; then
        touch ~/.bashrc
    fi
    # Adding resh shellrc to .bashrc ...
    grep -q '[[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc' ~/.bashrc ||\
        echo -e '\n[[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc # this line was added by RESH' >> ~/.bashrc
    # Adding bash-preexec to .bashrc ...
    grep -q '[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh' ~/.bashrc ||\
        echo -e '\n[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh # this line was added by RESH' >> ~/.bashrc
fi

# Only add shell directives into zsh if it passed version checks
if [ "$zsh_ok" = 1 ]; then
    # Adding resh shellrc to .zshrc ...
    if [ -f ~/.zshrc ]; then
        grep -q '[ -f ~/.resh/shellrc ] && source ~/.resh/shellrc' ~/.zshrc ||\
            echo -e '\n[ -f ~/.resh/shellrc ] && source ~/.resh/shellrc # this line was added by RESH' >> ~/.zshrc
    fi
fi

~/.resh/bin/resh-daemon-start

# bright green
high='\e[1m'
reset='\e[0m'

printf '
Installation finished successfully.


QUICK START
\e[32;1m    Press CTRL+R to launch RESH SEARCH    \e[0m
'
if [ -z "${__RESH_VERSION-}" ]; then
    printf 'You will need to restart your terminal first!\n'
fi
printf '
    Full-text search your shell history.
    Relevant results are displayed first based on current directory, git repo, and exit status.

    RESH will locally record and save shell history with context (directory, time, exit status, ...)
    Start using RESH right away because bash and zsh history are also searched.

    Update RESH by running: reshctl update
    Thank you for using RESH!
'

# Show banner if RESH is not loaded in the terminal
if [ -z "${__RESH_VERSION-}" ]; then printf '
##############################################################
#                                                            #
#    Finish the installation by RESTARTING this terminal!    #
#                                                            #
##############################################################
'
fi
