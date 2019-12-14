# util.sh - resh utility functions 
__resh_get_uuid() {
    cat /proc/sys/kernel/random/uuid 2>/dev/null || resh-uuid
}

__resh_get_pid() {
    if [ -n "$ZSH_VERSION" ]; then
        # assume Zsh
        local __RESH_PID="$$" # current pid
    elif [ -n "$BASH_VERSION" ]; then
        # assume Bash
        local __RESH_PID="$BASHPID" # current pid
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
    elif [ -n "$ZSH_VERSION" ]; then
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
    if [ -n "$ZSH_VERSION" ]; then
        setopt LOCAL_OPTIONS NO_NOTIFY NO_MONITOR
    fi
    nohup resh-daemon &>~/.resh/daemon_last_run_out.txt & disown
}

__resh_bash_completion_init() {
    local bash_completion_dir=~/.resh/bash_completion.d
    # source user completion directory definitions
    # taken from /usr/share/bash-completion/bash_completion
    if [[ -d $bash_completion_dir && -r $bash_completion_dir && \
        -x $bash_completion_dir ]]; then
        for i in $(LC_ALL=C command ls "$bash_completion_dir"); do
            i=$bash_completion_dir/$i
            # shellcheck disable=SC2154
            # shellcheck source=/dev/null
            [[ ${i##*/} != @($_backup_glob|Makefile*|$_blacklist_glob) \
                && -f $i && -r $i ]] && . "$i"
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
    source ~/.resh/zsh_completion.d/_reshctl && compdef _reshctl reshctl

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
                    &>~/.resh/session_init_last_run_out.txt || echo "resh-session-init ERROR: $(head -n 1 ~/.resh/session_init_last_run_out.txt)"
        fi
}
