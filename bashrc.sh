
PATH=$PATH:~/.resh/bin
#resh-daemon & disown

preexec() {
    __RESH_COLLECT=1
    __RESH_PWD="$PWD"
    __RESH_CMDLINE="$1"
    __RESH_RT_BEFORE=$EPOCHREALTIME
}

precmd() {
    __RESH_EXIT_CODE=$?
    __RESH_RT_AFTER=$EPOCHREALTIME
    if [ ! -z ${__RESH_COLLECT+x} ]; then
        resh-collect -cmd "$__RESH_CMDLINE" -status $__RESH_EXIT_CODE \
                     -pwd "$PWD" \
                     -realtimeBefore $__RESH_RT_BEFORE \
                     -realtimeAfter $__RESH_RT_AFTER
    fi
    unset __RESH_COLLECT
}

