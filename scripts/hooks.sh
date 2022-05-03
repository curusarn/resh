#!/hint/sh

__resh_reset_variables() {
    __RESH_RECORD_ID=$(__resh_get_uuid)
}

__resh_preexec() {
    # core
    __RESH_COLLECT=1
    __RESH_CMDLINE="$1" # not local to preserve it for postcollect (useful as sanity check)
    local fpath_last_run="$__RESH_XDG_CACHE_HOME/collect_last_run_out.txt"
    __resh_collect --cmdLine "$__RESH_CMDLINE" \
        >| "$fpath_last_run" 2>&1 || echo "resh-collect ERROR: $(head -n 1 $fpath_last_run)"
}

# used for collect and collect --recall
__resh_collect() {
    # posix
    local __RESH_COLS="$COLUMNS"
    local __RESH_LANG="$LANG"
    local __RESH_LC_ALL="$LC_ALL"
    local __RESH_LINES="$LINES"
    local __RESH_PWD="$PWD"
    
    # non-posix
    local __RESH_SHLVL="$SHLVL"
    local __RESH_GIT_CDUP; __RESH_GIT_CDUP="$(git rev-parse --show-cdup 2>/dev/null)"
    local __RESH_GIT_CDUP_EXIT_CODE=$?
    local __RESH_GIT_REMOTE; __RESH_GIT_REMOTE="$(git remote get-url origin 2>/dev/null)"
    local __RESH_GIT_REMOTE_EXIT_CODE=$?

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
    # time
    local __RESH_TZ_BEFORE; __RESH_TZ_BEFORE=$(date +%z)
    # __RESH_RT_BEFORE="$EPOCHREALTIME"
    __RESH_RT_BEFORE=$(__resh_get_epochrealtime)

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
                    -uname "$__RESH_UNAME" \
                    -sessionId "$__RESH_SESSION_ID" \
                    -recordId "$__RESH_RECORD_ID" \
                    -cols "$__RESH_COLS" \
                    -home "$__RESH_HOME" \
                    -lang "$__RESH_LANG" \
                    -lcAll "$__RESH_LC_ALL" \
                    -lines "$__RESH_LINES" \
                    -login "$__RESH_LOGIN" \
                    -pwd "$__RESH_PWD" \
                    -shellEnv "$__RESH_SHELL_ENV" \
                    -term "$__RESH_TERM" \
                    -pid "$__RESH_PID" \
                    -sessionPid "$__RESH_SESSION_PID" \
                    -host "$__RESH_HOST" \
                    -hosttype "$__RESH_HOSTTYPE" \
                    -ostype "$__RESH_OSTYPE" \
                    -machtype "$__RESH_MACHTYPE" \
                    -shlvl "$__RESH_SHLVL" \
                    -gitCdup "$__RESH_GIT_CDUP" \
                    -gitCdupExitCode "$__RESH_GIT_CDUP_EXIT_CODE" \
                    -gitRemote "$__RESH_GIT_REMOTE" \
                    -gitRemoteExitCode "$__RESH_GIT_REMOTE_EXIT_CODE" \
                    -realtimeBefore "$__RESH_RT_BEFORE" \
                    -realtimeSession "$__RESH_RT_SESSION" \
                    -realtimeSessSinceBoot "$__RESH_RT_SESS_SINCE_BOOT" \
                    -timezoneBefore "$__RESH_TZ_BEFORE" \
                    -osReleaseId "$__RESH_OS_RELEASE_ID" \
                    -osReleaseVersionId "$__RESH_OS_RELEASE_VERSION_ID" \
                    -osReleaseIdLike "$__RESH_OS_RELEASE_ID_LIKE" \
                    -osReleaseName "$__RESH_OS_RELEASE_NAME" \
                    -osReleasePrettyName "$__RESH_OS_RELEASE_PRETTY_NAME" \
                    "$@"
            return $?
        fi
        return 1
}

__resh_precmd() {
    local __RESH_EXIT_CODE=$?
    local __RESH_RT_AFTER
    local __RESH_TZ_AFTER
    local __RESH_PWD_AFTER
    local __RESH_GIT_CDUP_AFTER
    local __RESH_GIT_CDUP_EXIT_CODE_AFTER
    local __RESH_GIT_REMOTE_AFTER
    local __RESH_GIT_REMOTE_EXIT_CODE_AFTER
    local __RESH_SHLVL="$SHLVL"
    __RESH_RT_AFTER=$(__resh_get_epochrealtime)
    __RESH_TZ_AFTER=$(date +%z)
    __RESH_PWD_AFTER="$PWD"
    __RESH_GIT_CDUP_AFTER="$(git rev-parse --show-cdup 2>/dev/null)"
    __RESH_GIT_CDUP_EXIT_CODE_AFTER=$?
    __RESH_GIT_REMOTE_AFTER="$(git remote get-url origin 2>/dev/null)"
    __RESH_GIT_REMOTE_EXIT_CODE_AFTER=$?
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
            local fpath_last_run="$__RESH_XDG_CACHE_HOME/postcollect_last_run_out.txt"
            resh-postcollect -requireVersion "$__RESH_VERSION" \
                        -requireRevision "$__RESH_REVISION" \
                        -cmdLine "$__RESH_CMDLINE" \
                        -realtimeBefore "$__RESH_RT_BEFORE" \
                        -exitCode "$__RESH_EXIT_CODE" \
                        -sessionId "$__RESH_SESSION_ID" \
                        -recordId "$__RESH_RECORD_ID" \
                        -shell "$__RESH_SHELL" \
                        -shlvl "$__RESH_SHLVL" \
                        -pwdAfter "$__RESH_PWD_AFTER" \
                        -gitCdupAfter "$__RESH_GIT_CDUP_AFTER" \
                        -gitCdupExitCodeAfter "$__RESH_GIT_CDUP_EXIT_CODE_AFTER" \
                        -gitRemoteAfter "$__RESH_GIT_REMOTE_AFTER" \
                        -gitRemoteExitCodeAfter "$__RESH_GIT_REMOTE_EXIT_CODE_AFTER" \
                        -realtimeAfter "$__RESH_RT_AFTER" \
                        -timezoneAfter "$__RESH_TZ_AFTER" \
                        >| "$fpath_last_run" 2>&1 || echo "resh-postcollect ERROR: $(head -n 1 $fpath_last_run)"
        fi
        __resh_reset_variables
    fi
    unset __RESH_COLLECT
}
