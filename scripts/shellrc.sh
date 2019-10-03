
PATH=$PATH:~/.resh/bin
# if [ -n "$ZSH_VERSION" ]; then
#     zmodload zsh/datetime
# fi

# shellcheck source=reshctl.sh
. ~/.resh/reshctl.sh

__resh_get_uuid() {
    cat /proc/sys/kernel/random/uuid 2>/dev/null || resh-uuid
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
    nohup resh-daemon &>/dev/null & disown
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
    # shellcheck disable=SC2206
    fpath=(~/.resh/zsh_completion.d $fpath)
}

__RESH_MACOS=0
__RESH_LINUX=0
__RESH_UNAME=$(uname)

if [ "$__RESH_UNAME" = "Darwin" ]; then
    __RESH_MACOS=1
elif [ "$__RESH_UNAME" = "Linux" ]; then
    __RESH_LINUX=1
else
    echo "resh PANIC unrecognized OS"
fi

if [ -n "$ZSH_VERSION" ]; then
    # shellcheck disable=SC1009
    __RESH_SHELL="zsh"
    __RESH_HOST="$HOST"
    __RESH_HOSTTYPE="$CPUTYPE"
    __resh_zsh_completion_init
elif [ -n "$BASH_VERSION" ]; then
    __RESH_SHELL="bash"
    __RESH_HOST="$HOSTNAME"
    __RESH_HOSTTYPE="$HOSTTYPE"
    __resh_bash_completion_init
else
    echo "resh PANIC unrecognized shell"
fi

if [ -z "${__RESH_SESSION_ID+x}" ]; then
    export __RESH_SESSION_ID; __RESH_SESSION_ID=$(__resh_get_uuid)
    export __RESH_SESSION_PID="$$"
    # TODO add sesson time
fi

# posix
__RESH_HOME="$HOME"
__RESH_LOGIN="$LOGNAME"
__RESH_SHELL_ENV="$SHELL"
__RESH_TERM="$TERM"

# non-posix
__RESH_RT_SESSION=$(__resh_get_epochrealtime)
__RESH_OSTYPE="$OSTYPE" 
__RESH_MACHTYPE="$MACHTYPE"

if [ $__RESH_LINUX -eq 1 ]; then
    __RESH_OS_RELEASE_ID=$(. /etc/os-release; echo "$ID")
    __RESH_OS_RELEASE_VERSION_ID=$(. /etc/os-release; echo "$VERSION_ID")
    __RESH_OS_RELEASE_ID_LIKE=$(. /etc/os-release; echo "$ID_LIKE")
    __RESH_OS_RELEASE_NAME=$(. /etc/os-release; echo "$NAME")
    __RESH_OS_RELEASE_PRETTY_NAME=$(. /etc/os-release; echo "$PRETTY_NAME")
    __RESH_RT_SESS_SINCE_BOOT=$(cut -d' ' -f1 /proc/uptime)
elif [ $__RESH_MACOS -eq 1 ]; then
    __RESH_OS_RELEASE_ID="macos"
    __RESH_OS_RELEASE_VERSION_ID=$(sw_vers -productVersion 2>/dev/null)
    __RESH_OS_RELEASE_NAME="macOS"
    __RESH_OS_RELEASE_PRETTY_NAME="Mac OS X"
    __RESH_RT_SESS_SINCE_BOOT=$(sysctl -n kern.boottime | awk '{print $4}' | sed 's/,//g')
fi

__RESH_VERSION=$(resh-collect -version)
__RESH_REVISION=$(resh-collect -revision)

__resh_run_daemon

__resh_preexec() {
    # core
    __RESH_COLLECT=1
    __RESH_CMDLINE="$1"

    # posix
    __RESH_COLS="$COLUMNS"
    __RESH_LANG="$LANG"
    __RESH_LC_ALL="$LC_ALL"
    # other LC ?
    __RESH_LINES="$LINES"
    # __RESH_PATH="$PATH"
    __RESH_PWD="$PWD"
    
    # non-posix
    __RESH_SHLVL="$SHLVL"
    __RESH_GIT_CDUP="$(git rev-parse --show-cdup 2>/dev/null)"
    __RESH_GIT_CDUP_EXIT_CODE=$?
    __RESH_GIT_REMOTE="$(git remote get-url origin 2>/dev/null)"
    __RESH_GIT_REMOTE_EXIT_CODE=$?
    #__RESH_GIT_TOPLEVEL="$(git rev-parse --show-toplevel)"
    #__RESH_GIT_TOPLEVEL_EXIT_CODE=$?

    if [ -n "$ZSH_VERSION" ]; then
        # assume Zsh
        __RESH_PID="$$" # current pid
    elif [ -n "$BASH_VERSION" ]; then
        # assume Bash
        __RESH_PID="$BASHPID" # current pid
    fi
    # time
    __RESH_TZ_BEFORE=$(date +%z)
    # __RESH_RT_BEFORE="$EPOCHREALTIME"
    __RESH_RT_BEFORE=$(__resh_get_epochrealtime)

    # TODO: we should evaluate symlinks in preexec
    #       -> maybe create resh-precollect that could handle most of preexec
    #           maybe even move resh-collect here and send data to daemon and 
    #           send rest of the data ($?, timeAfter) to daemon in precmd 
    #           daemon will combine the data and save the record
    #           and save the unfinnished record even if it never finishes
    #           detect if the command died with the parent ps and save it then

}

__resh_precmd() {
    __RESH_EXIT_CODE=$?
    __RESH_RT_AFTER=$(__resh_get_epochrealtime)
    __RESH_TZ_AFTER=$(date +%z)
    __RESH_PWD_AFTER="$PWD"
    if [ -n "${__RESH_COLLECT}" ]; then
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
                        -cmdLine "$__RESH_CMDLINE" \
                        -exitCode "$__RESH_EXIT_CODE" \
                        -shell "$__RESH_SHELL" \
                        -uname "$__RESH_UNAME" \
                        -sessionId "$__RESH_SESSION_ID" \
                        -cols "$__RESH_COLS" \
                        -home "$__RESH_HOME" \
                        -lang "$__RESH_LANG" \
                        -lcAll "$__RESH_LC_ALL" \
                        -lines "$__RESH_LINES" \
                        -login "$__RESH_LOGIN" \
                        -pwd "$__RESH_PWD" \
                        -pwdAfter "$__RESH_PWD_AFTER" \
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
                        -realtimeAfter "$__RESH_RT_AFTER" \
                        -realtimeSession "$__RESH_RT_SESSION" \
                        -realtimeSessSinceBoot "$__RESH_RT_SESS_SINCE_BOOT" \
                        -timezoneBefore "$__RESH_TZ_BEFORE" \
                        -timezoneAfter "$__RESH_TZ_AFTER" \
                        -osReleaseId "$__RESH_OS_RELEASE_ID" \
                        -osReleaseVersionId "$__RESH_OS_RELEASE_VERSION_ID" \
                        -osReleaseIdLike "$__RESH_OS_RELEASE_ID_LIKE" \
                        -osReleaseName "$__RESH_OS_RELEASE_NAME" \
                        -osReleasePrettyName "$__RESH_OS_RELEASE_PRETTY_NAME" \
                        &>~/.resh/client_last_run_out.txt || echo "resh ERROR: $(head -n 1 ~/.resh/client_last_run_out.txt)"
                        # -path "$__RESH_PATH" \
        fi
    fi
    unset __RESH_COLLECT
}

preexec_functions+=(__resh_preexec)
precmd_functions+=(__resh_precmd)
