#!/hint/sh

# util.sh - resh utility functions 

__resh_get_pid() {
    if [ -n "${ZSH_VERSION-}" ]; then
        # assume Zsh
        local __RESH_PID="$$" # current pid
    elif [ -n "${BASH_VERSION-}" ]; then
        # assume Bash
        if [ "${BASH_VERSINFO[0]}" -ge "4" ]; then
            # $BASHPID is only available in bash4+
            # $$ is fairly similar so it should not be an issue
            local __RESH_PID="$BASHPID" # current pid
        else
            local __RESH_PID="$$" # current pid
        fi
    fi
    echo "$__RESH_PID"
}

__resh_get_epochrealtime() {
    if date +%s.%N | grep -vq 'N'; then
        # GNU date
        date +%s.%N
    elif gdate --version >/dev/null && gdate +%s.%N | grep -vq 'N'; then
        # GNU date take 2
        gdate +%s.%N
    elif [ -n "${ZSH_VERSION-}" ]; then
        # zsh fallback using $EPOCHREALTIME
        if [ -z "${__RESH_ZSH_LOADED_DATETIME+x}" ]; then
            zmodload zsh/datetime
            __RESH_ZSH_LOADED_DATETIME=1
        fi
        echo "$EPOCHREALTIME"
    else
        # dumb date 
        # XXX: we lost precison beyond seconds
        date +%s
        if [ -z "${__RESH_DATE_WARN+x}" ]; then
            echo "resh WARN: can't get precise time - consider installing GNU date!"
            __RESH_DATE_WARN=1
        fi
    fi
}

# FIXME: figure out if stdout/stderr should be discarded
__resh_run_daemon() {
    if [ -n "${ZSH_VERSION-}" ]; then
        setopt LOCAL_OPTIONS NO_NOTIFY NO_MONITOR
    fi
    if [ "$(uname)" = Darwin ]; then
        # hotfix
        gnohup resh-daemon >/dev/null 2>/dev/null & disown
    else
        # TODO: switch to nohup for consistency once you confirm that daemon is
        #       not getting killed anymore on macOS
        nohup resh-daemon >/dev/null 2>/dev/null & disown
        #setsid resh-daemon 2>&1 & disown
    fi
}

__resh_session_init() {
    if [ "$__RESH_VERSION" != "$(resh-session-init -version)" ]; then
        # shellcheck source=shellrc.sh
        source ~/.resh/shellrc 
        if [ "$__RESH_VERSION" != "$(resh-session-init -version)" ]; then
            echo "RESH WARNING: You probably just updated RESH - PLEASE RESTART OR RELOAD THIS TERMINAL SESSION (resh version: $(resh-session-init -version); resh version of this terminal session: ${__RESH_VERSION})"
        else
            echo "RESH INFO: New RESH shellrc script was loaded - if you encounter any issues please restart this terminal session."
        fi
    elif [ "$__RESH_REVISION" != "$(resh-session-init -revision)" ]; then
        # shellcheck source=shellrc.sh
        source ~/.resh/shellrc 
        if [ "$__RESH_REVISION" != "$(resh-session-init -revision)" ]; then
            echo "RESH WARNING: You probably just updated RESH - PLEASE RESTART OR RELOAD THIS TERMINAL SESSION (resh revision: $(resh-session-init -revision); resh revision of this terminal session: ${__RESH_REVISION})"
        fi
    fi
    if [ "$__RESH_VERSION" = "$(resh-session-init -version)" ] && [ "$__RESH_REVISION" = "$(resh-session-init -revision)" ]; then
        resh-session-init -requireVersion "$__RESH_VERSION" \
                    -requireRevision "$__RESH_REVISION" \
                    -sessionId "$__RESH_SESSION_ID" \
                    -sessionPid "$__RESH_SESSION_PID"
    fi
}
