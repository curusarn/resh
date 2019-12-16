#!/usr/bin/env bash

set -euo pipefail

die() {
    if [ $# -eq 1 ]; then
        echo "$1"
    elif [ $# -eq 0 ]; then
        echo "ERROR: SomEtHiNg WeNt wRonG - ExItiNg!"
        echo "THIS IS NOT SUPPOSED TO HAPPEN!"
    fi
    echo
    echo "Please report any issues you encounter to: https://github.com/curusarn/resh/issues"
    echo
    echo "You can rerun this installation by executing: (you will skip downloading the files)"
    echo "cd $PWD && scripts/install.sh"
    exit 1
}

echo "Checking your system ..."

# /usr/bin/zsh -> zsh
login_shell=$(echo "$SHELL" | rev | cut -d'/' -f1 | rev)

if [ "$login_shell" != bash ] && [ "$login_shell" != zsh ]; then
    die "ERROR: Unsupported/unknown login shell: $login_shell"
fi
echo "Login shell: $login_shell - OK"


# check like we are not running bash
bash_version=$(bash -c 'echo ${BASH_VERSION}')
bash_version_major=$(bash -c 'echo ${BASH_VERSINFO[0]}')
bash_version_minor=$(bash -c 'echo ${BASH_VERSINFO[1]}')
bash_too_old=""
if [ "$bash_version_major" -lt 3 ]; then 
    bash_too_old=true
elif [ "$bash_version_major" -eq 4 ] && [ "$bash_version_minor" -lt 3 ]; then 
    bash_too_old=true
fi

if [ "$bash_too_old" = true ]; then
    echo "Bash version: $bash_version - UNSUPPORTED!"
    if [ "$login_shell" = bash ]; then
        echo " > Your bash version is old."
        echo " > Bash is also your login shell."
        echo " > Updating to bash 4.3+ is strongly RECOMMENDED!"
    else
        echo " > Your bash version is old"
        echo " > Bash is not your login shell so it should not be an issue."
        echo " > Updating to bash 4.3+ is recommended."
    fi
else
    echo "Bash version: $bash_version - OK"
fi


if ! zsh --version &>/dev/null; then
    echo "Zsh version: ? - not installed!"
else
    zsh_version=$(zsh -c 'echo ${ZSH_VERSION}')
    zsh_version_major=$(echo "$zsh_version" | cut -d'.' -f1)
    if [ "$zsh_version_major" -lt 5 ]; then 
        echo "Zsh version: $zsh_version - UNSUPPORTED!"
        if [ "$login_shell" = zsh ]; then
            echo " > Your zsh version is old."
            echo " > Zsh is also your login shell."
            echo " > Updating to Zsh 5.0+ is strongly RECOMMENDED!"
        else
            echo " > Your zsh version is old"
            echo " > Zsh is not your login shell so it should not be an issue."
            echo " > Updating to zsh 5.0+ is recommended."
        fi
    else
        echo "Zsh version: $zsh_version - OK"
    fi
fi


echo 
echo "Continue with installation? (Any key to CONTINUE / Ctrl+C to ABORT)"
echo
echo "Creating directories ..."

mkdir_if_not_exists() {
    if [ ! -d "$1" ]; then
        mkdir "$1" 
    fi
}

mkdir_if_not_exists ~/.resh
mkdir_if_not_exists ~/.resh/bin
mkdir_if_not_exists ~/.resh/bash_completion.d
mkdir_if_not_exists ~/.resh/zsh_completion.d
mkdir_if_not_exists ~/.config

echo "Copying files ..."
cp -f submodules/bash-preexec/bash-preexec.sh ~/.bash-preexec.sh
cp -f submodules/bash-zsh-compat-widgets/bindfunc.sh ~/.resh/bindfunc.sh

cp -f conf/config.toml ~/.config/resh.toml

cp -f scripts/shellrc.sh ~/.resh/shellrc
cp -f scripts/reshctl.sh scripts/widgets.sh scripts/hooks.sh scripts/util.sh ~/.resh/

echo "Generating completions for reshctl ..."
bin/resh-control completion bash > ~/.resh/bash_completion.d/_reshctl
bin/resh-control completion zsh > ~/.resh/zsh_completion.d/_reshctl

echo "Copying more files ..."
cp -f scripts/uuid.sh ~/.resh/bin/resh-uuid
cp -f bin/* ~/.resh/bin/
cp -f scripts/resh-evaluate-plot.py ~/.resh/bin/
cp -fr data/sanitizer ~/.resh/sanitizer_data

# backward compatibility: We have a new location for resh history file 
[ ! -f ~/.resh/history.json ] || mv ~/.resh/history.json ~/.resh_history.json 

echo "Finishing up ..."
# Adding resh shellrc to .bashrc ...
grep -q '[[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc' ~/.bashrc ||\
	echo -e '\n[[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc' >> ~/.bashrc
# Adding bash-preexec to .bashrc ...
grep -q '[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh' ~/.bashrc ||\
	echo -e '\n[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh' >> ~/.bashrc
# Adding resh shellrc to .zshrc ...
grep -q '[ -f ~/.resh/shellrc ] && source ~/.resh/shellrc' ~/.zshrc ||\
	echo -e '\n[ -f ~/.resh/shellrc ] && source ~/.resh/shellrc' >> ~/.zshrc

# Deleting zsh completion cache - for future use
# [ ! -e ~/.zcompdump ] || rm ~/.zcompdump

# Final touch
touch ~/.resh_history.json

# Restarting resh daemon ...
if [ -f ~/.resh/resh.pid ]; then
    kill -SIGTERM "$(cat ~/.resh/resh.pid)"
    rm ~/.resh/resh.pid
else
    pkill -SIGTERM "resh-daemon" || true
fi
nohup resh-daemon &>/dev/null & disown

# Generating resh-uuid ...
[ -e "$(HOME)/.resh/resh-uuid" ] \
	|| cat /proc/sys/kernel/random/uuid > "$(HOME)/.resh/resh-uuid" 2>/dev/null \
	|| scripts/uuid.sh > "$(HOME)/.resh/resh-uuid" 2>/dev/null 

echo "\ 

##########################################################
#                                                        #
#    SUCCESS - thank you for trying out this project!    #
#                                                        #
##########################################################

 WARNING 
 It's recommended to RESTART ALL OPEN TERMINAL WINDOWS (or reload your rc files)

 HISTORY
 Your resh history will be recorded to '~/.resh_history.json'
 You can look at it using e.g. 'tail -f ~/.resh_history.json | jq' (you might need to install jq)

 SANITIZATION
 History can be sanitized by running '... to be included'
 This will create sanitized version of your history

 GRAPHS
 You can get some graphs of your history by running '... to be included'

 ISSUES
 Please report issues to: https://github.com/curusarn/resh/issues

 UNINSTALL
 You can uninstall this at any time by running 'rm -rf ~/.resh/'
 You won't lose any collected history by removing '~/.resh/' directory

 Please give me some contact info using this form: https://forms.gle/227SoyJ5c2iteKt98

"