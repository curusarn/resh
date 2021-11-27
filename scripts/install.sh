#!/usr/bin/env bash

set -euo pipefail

# Setting this to `1` skips prompts and always uses the deault option
SKIP_ASK_PROMPTS=0

# Helper for "ask [Y/n]" and "ask [y/N]"
ask() {
    yn="$1"
    shift
    case "$yn" in
        0) options="[Y/n]" ;;
        1) options="[y/N]" ;;
        *)
            echo "FATAL: ask() got invalid argument."
            exit 2
        ;;
    esac
    if [ "$SKIP_ASK_PROMPTS" = "0" ]; then
        echo
        # We are using echo -e to allow multiline messages
        echo -en "$@"
        printf " %s\n" "$options"
        read reply
    else
        reply=""
    fi
    case "$reply" in
        y*|Y*) return 0 ;;
        n*|N*) return 1 ;;
        *)     return "$yn" ;;
    esac
}
ask_Yn() {
    ask 0 "$@"
}
ask_yN() {
    ask 1 "$@"
}

echo
echo "Checking your system ..."

# /usr/bin/zsh -> zsh
login_shell=$(echo "$SHELL" | rev | cut -d'/' -f1 | rev)

if [ "$login_shell" != bash ] && [ "$login_shell" != zsh ]; then
    echo "ERROR: Unsupported/unknown login shell: $login_shell"
    exit 1
fi
echo " * Login shell: $login_shell - OK"


# check like if we are not running bash
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
        echo "   > Updating to bash 4.3+ is strongly RECOMMENDED!"
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
            echo "   > Updating to Zsh 5.0+ is strongly RECOMMENDED!"
        else
            echo "   > Your zsh version is old"
            echo "   > Zsh is not your login shell so it should not be an issue."
            echo "   > Updating to zsh 5.0+ is recommended."
        fi
    else
        echo " * Zsh version: $zsh_version - OK"
    fi
fi


if [ "$(uname)" = Darwin ]; then
    if gnohup --version >/dev/null 2>&1; then
        echo " * Nohup installed: OK"
    else
        echo " * Nohup installed: NOT INSTALLED!"
        echo "   > You don't have nohup"
        echo "   > Please install GNU coreutils"
        echo
        echo "   $ brew install coreutils"
        echo
        exit 1
    fi
else
    if setsid --version >/dev/null 2>&1; then
        echo " * Setsid installed: OK"
    else
        echo " * Setsid installed: NOT INSTALLED!"
        echo "   > You don't have setsid"
        echo "   > Please install unix-util"
        echo
        exit 1
    fi
fi

echo
echo "This installations has two modes:"
echo " * Automatic - no question asked - RECOMMENDED"
echo " * Guided - prompts to give you more control - useful if you have heavily customized shell configuration"
if ask_Yn ">>> Would you like to use the AUTOMATIC install mode?"; then
    SKIP_ASK_PROMPTS=1
    echo "Using automatic install mode ..."
else
    echo "Using guided install mode ..."
fi

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

cp -f scripts/shellrc.sh ~/.resh/shellrc
cp -f scripts/reshctl.sh scripts/widgets.sh scripts/hooks.sh scripts/util.sh ~/.resh/
cp -f scripts/rawinstall.sh ~/.resh/

update_config() {
    version=$1
    key=$2
    value=$3
    # TODO: create bin/semver-lt
    if bin/semver-lt "${__RESH_VERSION:-0.0.0}" "$1" && [ "$(bin/resh-config -key $key)" != "$value" ] ; then
        echo " * config option $key was updated to $value"
        # TODO: enable resh-config value setting
        # resh-config -key "$key" -value "$value"
    fi
}

# Do not overwrite config if it exists
if [ ! -f ~/.config/resh.toml ]; then
    echo "Copying config file ..."
    cp -f conf/config.toml ~/.config/resh.toml
# else 
    # echo "Merging config files ..."
    # NOTE: This is where we will merge configs when we make changes to the upstream config
    # HINT: check which version are we updating FROM and make changes to config based on that 
fi

echo "Generating shell completions ..."
bin/resh-control completion bash > ~/.resh/bash_completion.d/_reshctl
bin/resh-control completion zsh > ~/.resh/zsh_completion.d/_reshctl

echo "Copying more files ..."
cp -f scripts/uuid.sh ~/.resh/bin/resh-uuid
cp -f bin/* ~/.resh/bin/
cp -f scripts/resh-evaluate-plot.py ~/.resh/bin/
cp -fr data/sanitizer ~/.resh/sanitizer_data

# backward compatibility: We have a new location for resh history file 
[ ! -f ~/.resh/history.json ] || mv ~/.resh/history.json ~/.resh_history.json 

echo "Adding RESH to shell rc files ..."

setup_bashrc() {
    # Creating .bashrc ...
    if [ ! -f ~/.bashrc ]; then
        ask_Yn \
           "It looks like there is no '~/.bashrc'." \
           "RESH must be sourced when your shell starts otherwise it won't work." \
           ">>> Create '~/.bashrc'?"
        if [ "$?" = "0" ]; then
            echo "Creating '~/.bashrc' ..."
            touch ~/.bashrc
        else
            return 1
        fi
    fi
    # Adding resh shellrc to .bashrc ...
    if ! grep -q '[[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc' ~/.bashrc; then
        ask_Yn \
           "It looks like there is no '[[ -f ~/.resh/shellrc ]] && source ~/.resh/shelrc' directive in your '~/.bashrc'." \
           "\nRESH must be sourced when your shell starts otherwise it won't work." \
           "\n>>> Add source directive to '~/.bashrc'?"
        if [ "$?" = "0" ]; then
            echo "Adding RESH source directive to '~/.bashrc' ..."
            echo -e '\n[[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc # this line was added by RESH (Rich Enchanced Shell History)' >> ~/.bashrc
        else
            return 1
        fi
    fi
    # Adding bash-preexec to .bashrc ...
    if ! grep -q '[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh' ~/.bashrc; then
        ask_Yn \
           "It looks like there is no '[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh' directive in your '~/.bashrc'." \
           "\nBash-preexec must be sourced when your shell starts otherwise RESH won't work." \
           "\n>>> Add bash-preexec source directive to '~/.bashrc'?"
        if [ "$?" = "0" ]; then
            echo "Adding bash-preexec source directive to '~/.bashrc' ..."
            echo -e '\n[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh # this line was added by RESH (Rich Enchanced Shell History)' >> ~/.bashrc
        else
            return 1
        fi
    fi
}

setup_zshrc() {
    # Creating .zshrc ...
    if [ ! -f ~/.zshrc ]; then
        echo "There is no '~/.zshrc' - skipping zsh setup. (This is fine if you don't use zsh.)"
        return 0
    fi
    # Adding resh shellrc to .zshrc ...
    if ! grep -q '[[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc' ~/.zshrc; then
        ask_Yn \
           "It looks like there is no '[[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc' directive in your '~/.zshrc'." \
           "\nRESH must be sourced when your shell starts otherwise it won't work." \
           "\n>>> Add source directive to '~/.zshrc'?"
        if [ "$?" = "0" ]; then
            echo "Adding RESH source directive to '~/.zshrc' ..."
            echo -e '\n[[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc # this line was added by RESH (Rich Enchanced Shell History)' >> ~/.zshrc
        else
            return 1
        fi
    fi
}

bash_setup_info() {
    echo
    echo "WARNING: Shell config setup didn't complete for BASH! (You probably answered 'no' somewhere.)"
    echo
    echo "It is likely that you will need to modify your bash startup scripts to make RESH work in bash."
    echo "Consider rerunning the installation if you do not want to modify your shell configs yourself."
    echo
    echo "Instructions for manual config setup for BASH:"
    echo " 1) Add following lines to the end of your bash startup file:"
    echo
    echo "    [[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc # enable RESH"
    echo "    [[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh # enable bash-preexec (required by RESH)"
    echo
    echo " 2) Make sure that you have the right config - '~/.bashrc' is usually the right one."
    echo " 3) Make sure that you added the lines to the *end* of the file."
    echo " 4) Make sure that you added the lines in the correct order - bash-preexec needs to come last."
    echo
    echo "Press any key to continue ..."
    read -n 1 x
    echo
    echo
}
zsh_setup_info() {
    echo
    echo "WARNING: Shell config setup didn't complete for ZSH! (You probably answered 'no' somewhere.)"
    echo
    echo "It is likely that you will need to modify your zsh startup scripts to make RESH work in zsh."
    echo "Consider rerunning the installation if you do not want to modify your shell configs yourself."
    echo
    echo "Instructions for manual config setup for ZSH:"
    echo " 1) Add following line to the end of your zsh startup file:"
    echo
    echo "    [[ -f ~/.resh/shellrc ]] && source ~/.resh/shellrc # enable RESH"
    echo
    echo " 2) Make sure that you have the right config - '~/.zshrc' is usually the right one."
    echo " 3) Make sure that you added the line to the *end* of the file."
    echo
    echo "Press any key to continue ..."
    read -n 1 x
    echo
    echo
}

if ! setup_bashrc; then
    bash_setup_info
fi

if ! setup_zshrc; then
    zsh_setup_info
fi

echo "Finishing up ..."

# Deleting zsh completion cache - for future use
# [ ! -e ~/.zcompdump ] || rm ~/.zcompdump

# Final touch
touch ~/.resh_history.json

# Generating resh-uuid ...
[ -e ~/.resh/resh-uuid ] \
	|| cat /proc/sys/kernel/random/uuid > ~/.resh/resh-uuid 2>/dev/null \
	|| scripts/uuid.sh > ~/.resh/resh-uuid 2>/dev/null 

# Source utils to get __resh_run_daemon function
# shellcheck source=util.sh
. ~/.resh/util.sh

# Restarting resh daemon ...
if [ -f ~/.resh/resh.pid ]; then
    kill -SIGTERM "$(cat ~/.resh/resh.pid)" || true
    rm ~/.resh/resh.pid
else
    pkill -SIGTERM "resh-daemon" || true
fi
__resh_run_daemon


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
    (you will need to restart your teminal first)
     
    Search your history by commands. 
    Host, directories, git remote, and exit status is used to display relevant results first.

    At first, the search application will use the standard shell history without context. 
    All history recorded from now on will have context which will be used by the RESH SEARCH app.

    Enable/disable Ctrl+R binding using reshctl command:
     $ reshctl enable ctrl_r_binding
     $ reshctl disable ctrl_r_binding

CHECK FOR UPDATES
    To check for (and install) updates use reshctl command:
     $ reshctl update

HISTORY
    Your resh history will be recorded to '~/.resh_history.json'
    Look at it using e.g. following command (you might need to install jq)
     $ tail -f ~/.resh_history.json | jq

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


echo "$info" | ${PAGER:-less}

echo
echo "All done!"
echo "Thank you for using RESH"
echo "Issues go here: https://github.com/curusarn/resh/issues"
echo "Ctrl+R launches the RESH SEARCH app"
# echo "Do not forget to restart your terminal"
if [ -z "${__RESH_VERSION:-}" ]; then echo "
##############################################################
#                                                            #
#    Finish the installation by RESTARTING this terminal!    #
#                                                            #
##############################################################"
fi