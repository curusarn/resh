#!/hint/sh

# util.sh - resh utility functions 

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
