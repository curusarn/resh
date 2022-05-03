#!/hint/sh

# util.sh - resh utility functions 
__resh_get_uuid() {
    cat /proc/sys/kernel/random/uuid 2>/dev/null || resh-uuid
}

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

__resh_run_daemon() {
    if [ -n "${ZSH_VERSION-}" ]; then
        setopt LOCAL_OPTIONS NO_NOTIFY NO_MONITOR
    fi
    local fpath_last_run="$__RESH_XDG_CACHE_HOME/daemon_last_run_out.txt"
    if [ "$(uname)" = Darwin ]; then
        # hotfix
        gnohup resh-daemon >| "$fpath_last_run" 2>&1 & disown
    else
        # TODO: switch to nohup for consistency once you confirm that daemon is
        #       not getting killed anymore on macOS
        # nohup resh-daemon >| "$fpath_last_run" 2>&1 & disown
        setsid resh-daemon >| "$fpath_last_run" 2>&1 & disown
    fi
}

__resh_bash_completion_init() {
    # primitive check to find out if bash_completions are installed
    # skip completion init if they are not
    _get_comp_words_by_ref >/dev/null 2>/dev/null
    [[ $? == 127 ]] && return
    local bash_completion_dir=~/.resh/bash_completion.d
    if [[ -d $bash_completion_dir && -r $bash_completion_dir && \
        -x $bash_completion_dir ]]; then
        for i in $(LC_ALL=C command ls "$bash_completion_dir"); do
            i=$bash_completion_dir/$i
            # shellcheck disable=SC2154
            # shellcheck source=/dev/null
            [[ -f "$i" && -r "$i" ]] && . "$i"
        done
    fi
}

__resh_zsh_completion_init() {
    # NOTE: this is hacky - each completion needs to be added individually 
    # TODO: fix later
    # fpath=(~/.resh/zsh_completion.d $fpath)
    # we should be using fpath but that doesn't work well with oh-my-zsh
    #   so we are just adding it manually 
    # shellcheck disable=1090
    if typeset -f compdef >/dev/null 2>&1; then
        source ~/.resh/zsh_completion.d/_reshctl && compdef _reshctl reshctl
    else
        # fallback I guess
        fpath=(~/.resh/zsh_completion.d $fpath)
        __RESH_zsh_no_compdef=1
    fi

    # TODO: test and use this
    # NOTE: this is not how globbing works
    # for f in ~/.resh/zsh_completion.d/_*; do
    #   source ~/.resh/zsh_completion.d/_$f && compdef _$f $f
    # done
}

__resh_session_init() {
    # posix
    local __RESH_COLS="$COLUMNS"
    local __RESH_LANG="$LANG"
    local __RESH_LC_ALL="$LC_ALL"
    # other LC ?
    local __RESH_LINES="$LINES"
    local __RESH_PWD="$PWD"
    
    # non-posix
    local __RESH_SHLVL="$SHLVL"

    # pid
    local __RESH_PID; __RESH_PID=$(__resh_get_pid)

    # time
    local __RESH_TZ_BEFORE; __RESH_TZ_BEFORE=$(date +%z)
    local __RESH_RT_BEFORE; __RESH_RT_BEFORE=$(__resh_get_epochrealtime)

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
        local fpath_last_run="$__RESH_XDG_CACHE_HOME/session_init_last_run_out.txt"
        resh-session-init -requireVersion "$__RESH_VERSION" \
                    -requireRevision "$__RESH_REVISION" \
                    -shell "$__RESH_SHELL" \
                    -uname "$__RESH_UNAME" \
                    -sessionId "$__RESH_SESSION_ID" \
                    -cols "$__RESH_COLS" \
                    -home "$__RESH_HOME" \
                    -lang "$__RESH_LANG" \
                    -lcAll "$__RESH_LC_ALL" \
                    -lines "$__RESH_LINES" \
                    -login "$__RESH_LOGIN" \
                    -shellEnv "$__RESH_SHELL_ENV" \
                    -term "$__RESH_TERM" \
                    -pid "$__RESH_PID" \
                    -sessionPid "$__RESH_SESSION_PID" \
                    -host "$__RESH_HOST" \
                    -hosttype "$__RESH_HOSTTYPE" \
                    -ostype "$__RESH_OSTYPE" \
                    -machtype "$__RESH_MACHTYPE" \
                    -shlvl "$__RESH_SHLVL" \
                    -realtimeBefore "$__RESH_RT_BEFORE" \
                    -realtimeSession "$__RESH_RT_SESSION" \
                    -realtimeSessSinceBoot "$__RESH_RT_SESS_SINCE_BOOT" \
                    -timezoneBefore "$__RESH_TZ_BEFORE" \
                    -osReleaseId "$__RESH_OS_RELEASE_ID" \
                    -osReleaseVersionId "$__RESH_OS_RELEASE_VERSION_ID" \
                    -osReleaseIdLike "$__RESH_OS_RELEASE_ID_LIKE" \
                    -osReleaseName "$__RESH_OS_RELEASE_NAME" \
                    -osReleasePrettyName "$__RESH_OS_RELEASE_PRETTY_NAME" \
                    >| "$fpath_last_run" 2>&1 || echo "resh-session-init ERROR: $(head -n 1 $fpath_last_run)"
        fi
}

__resh_set_xdg_home_paths() {
    if [ -z "${XDG_CACHE_HOME-}" ]; then
        __RESH_XDG_CACHE_HOME="$HOME/.cache/resh"
    else
        __RESH_XDG_CACHE_HOME="$XDG_CACHE_HOME/resh"
    fi
    mkdir -p "$__RESH_XDG_CACHE_HOME" >/dev/null 2>/dev/null
    export __RESH_XDG_CACHE_HOME


    if [ -z "${XDG_DATA_HOME-}" ]; then
        __RESH_XDG_DATA_HOME="$HOME/.local/share/resh"
    else
        __RESH_XDG_DATA_HOME="$XDG_DATA_HOME/resh"
    fi
    mkdir -p "$__RESH_XDG_DATA_HOME" >/dev/null 2>/dev/null
    export __RESH_XDG_DATA_HOME
}
