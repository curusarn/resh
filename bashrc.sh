
PATH=$PATH:~/.resh/bin
#resh-daemon & disown

preexec() {
    __RESH_COLLECT=1
    __RESH_PWD="$PWD"
    __RESH_CMDLINE="$1"
}

precmd() {
    __RESH_EXIT_CODE=$?
    if [ ! -z ${__RESH_COLLECT+x} ]; then
        resh-collect $__RESH_EXIT_CODE "$PWD" "$__RESH_CMDLINE"
    fi
    unset __RESH_COLLECT
}

