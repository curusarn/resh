
PATH=$PATH:~/.resh/bin

preexec() {
    __RESH_PWD="$PWD"
    __RESH_CMDLINE="$1"
}

precmd() {
    __RESH_EXIT_CODE=$?
    resh-collect $__RESH_EXIT_CODE "$PWD" "$__RESH_CMDLINE"
}

