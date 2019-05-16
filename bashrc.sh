
PATH=$PATH:~/.resh/bin
export __RESH_RT_SESSION=$EPOCHREALTIME
export __RESH_RT_SESS_SINCE_BOOT=$(cat /proc/uptime | cut -d' ' -f1)
#resh-daemon & disown

preexec() {
    # core
    __RESH_COLLECT=1
    __RESH_CMDLINE="$1"

    # posix
    __RESH_COLS="$COLUMNS"
    __RESH_HOME="$HOME"
    __RESH_LANG="$LANG"
    __RESH_LC_ALL="$LC_ALL"
    # other LC ?
    __RESH_LINES="$LINES"
    __RESH_LOGIN="$LOGNAME"
    __RESH_PATH="$PATH"
    __RESH_PWD="$PWD"
    __RESH_SHELL="$SHELL"
    __RESH_TERM="$TERM"
    
    # non-posix
    __RESH_PID="$BASHPID" # current pid
    __RESH_SESSION_PID="$$" # pid of original shell 
    __RESH_WINDOWID="$WINDOWID" # session 
    __RESH_HOST="$HOSTNAME"
    __RESH_HOSTTYPE="$HOSTTYPE"
    __RESH_OSTYPE="$OSTYPE"
    __RESH_MACHTYPE="$MACHTYPE"

    # time
    __RESH_TZ_BEFORE=$(date +%:z)
    __RESH_RT_BEFORE="$EPOCHREALTIME"
}

precmd() {
    __RESH_EXIT_CODE=$?
    __RESH_RT_AFTER=$EPOCHREALTIME
    __RESH_SECS_UTC_AFTER=$(date +%s -u)
    __RESH_TZ_AFTER=$(date +%:z)
    if [ ! -z ${__RESH_COLLECT+x} ]; then
        resh-collect -cmdLine "$__RESH_CMDLINE" -exitCode "$__RESH_EXIT_CODE" \
                     -cols "$__RESH_COLS" \
                     -home "$__RESH_HOME" \
                     -lang "$__RESH_LANG" \
                     -lcAll "$__RESH_LC_ALL" \
                     -lines "$__RESH_LINES" \
                     -login "$__RESH_LOGIN" \
                     -path "$__RESH_PATH" \
                     -pwd "$__RESH_PWD" \
                     -shell "$__RESH_SHELL" \
                     -term "$__RESH_TERM" \
                     -pid "$__RESH_PID" -sessionPid "$__RESH_SESSION_PID" \
                     -windowId "$__RESH_WINDOWID" \
                     -host "$__RESH_HOST" \
                     -hosttype "$__RESH_HOSTTYPE" \
                     -ostype "$__RESH_OSTYPE" \
                     -machtype "$__RESH_MACHTYPE" \
                     -realtimeBefore "$__RESH_RT_BEFORE" \
                     -realtimeAfter "$__RESH_RT_AFTER" \
                     -realtimeSession "$__RESH_RT_SESSION" \
                     -realtimeSessSinceBoot "$__RESH_RT_SESS_SINCE_BOOT" \
                     -timezoneBefore "$__RESH_TZ_BEFORE" \
                     -timezoneAfter "$__RESH_TZ_AFTER"
    fi
    unset __RESH_COLLECT
}

