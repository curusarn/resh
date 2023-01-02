#!/hint/sh

__resh_reset_variables() {
    __RESH_RECORD_ID=$(resh-generate-uuid)
}

__resh_preexec() {
    # core
    __RESH_COLLECT=1
    __RESH_CMDLINE="$1" # not local to preserve it for postcollect (useful as sanity check)
    __resh_collect --cmdLine "$__RESH_CMDLINE"
}

# used for collect and collect --recall
__resh_collect() {
    # posix
    local __RESH_PWD="$PWD"
    
    # non-posix
    local __RESH_SHLVL="$SHLVL"
    local __RESH_GIT_REMOTE; __RESH_GIT_REMOTE="$(git remote get-url origin 2>/dev/null)"

    __RESH_RT_BEFORE=$(resh-get-epochtime)

    if [ "$__RESH_VERSION" != "$(resh-collect -version)" ]; then
        # shellcheck source=shellrc.sh
        source ~/.resh/shellrc 
        if [ "$__RESH_VERSION" != "$(resh-collect -version)" ]; then
            echo "RESH WARNING: You probably just updated RESH - PLEASE RESTART OR RELOAD THIS TERMINAL SESSION (resh version: $(resh-collect -version); resh version of this terminal session: ${__RESH_VERSION})"
        else
            echo "RESH INFO: New RESH shellrc script was loaded - if you encounter any issues please restart this terminal session."
        fi
    elif [ "$__RESH_REVISION" != "$(resh-collect -revision)" ]; then
        # shellcheck source=shellrc.sh
        source ~/.resh/shellrc
        if [ "$__RESH_REVISION" != "$(resh-collect -revision)" ]; then
            echo "RESH WARNING: You probably just updated RESH - PLEASE RESTART OR RELOAD THIS TERMINAL SESSION (resh revision: $(resh-collect -revision); resh revision of this terminal session: ${__RESH_REVISION})"
        fi
    fi
    if [ "$__RESH_VERSION" = "$(resh-collect -version)" ] && [ "$__RESH_REVISION" = "$(resh-collect -revision)" ]; then
        resh-collect -requireVersion "$__RESH_VERSION" \
                    -requireRevision "$__RESH_REVISION" \
                    -shell "$__RESH_SHELL" \
                    -sessionID "$__RESH_SESSION_ID" \
                    -recordID "$__RESH_RECORD_ID" \
                    -home "$__RESH_HOME" \
                    -pwd "$__RESH_PWD" \
                    -sessionPID "$__RESH_SESSION_PID" \
                    -shlvl "$__RESH_SHLVL" \
                    -gitRemote "$__RESH_GIT_REMOTE" \
                    -time "$__RESH_RT_BEFORE" \
                    "$@"
        return $?
    fi
    return 1
}

__resh_precmd() {
    local __RESH_EXIT_CODE=$?
    local __RESH_RT_AFTER
    local __RESH_SHLVL="$SHLVL"
    __RESH_RT_AFTER=$(resh-get-epochtime)
    if [ -n "${__RESH_COLLECT}" ]; then
        if [ "$__RESH_VERSION" != "$(resh-postcollect -version)" ]; then
            # shellcheck source=shellrc.sh
            source ~/.resh/shellrc 
            if [ "$__RESH_VERSION" != "$(resh-postcollect -version)" ]; then
                echo "RESH WARNING: You probably just updated RESH - PLEASE RESTART OR RELOAD THIS TERMINAL SESSION (resh version: $(resh-collect -version); resh version of this terminal session: ${__RESH_VERSION})"
            else
                echo "RESH INFO: New RESH shellrc script was loaded - if you encounter any issues please restart this terminal session."
            fi
        elif [ "$__RESH_REVISION" != "$(resh-postcollect -revision)" ]; then
            # shellcheck source=shellrc.sh
            source ~/.resh/shellrc 
            if [ "$__RESH_REVISION" != "$(resh-postcollect -revision)" ]; then
                echo "RESH WARNING: You probably just updated RESH - PLEASE RESTART OR RELOAD THIS TERMINAL SESSION (resh revision: $(resh-collect -revision); resh revision of this terminal session: ${__RESH_REVISION})"
            fi
        fi
        if [ "$__RESH_VERSION" = "$(resh-postcollect -version)" ] && [ "$__RESH_REVISION" = "$(resh-postcollect -revision)" ]; then
            resh-postcollect -requireVersion "$__RESH_VERSION" \
                        -requireRevision "$__RESH_REVISION" \
                        -timeBefore "$__RESH_RT_BEFORE" \
                        -exitCode "$__RESH_EXIT_CODE" \
                        -sessionID "$__RESH_SESSION_ID" \
                        -recordID "$__RESH_RECORD_ID" \
                        -shlvl "$__RESH_SHLVL" \
                        -timeAfter "$__RESH_RT_AFTER"
        fi
        __resh_reset_variables
    fi
    unset __RESH_COLLECT
}
