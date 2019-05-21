
PATH=$PATH:~/.resh/bin
# if [ -n "$ZSH_VERSION" ]; then
#     zmodload zsh/datetime
# fi

# export __RESH_RT_SESSION=$EPOCHREALTIME
export __RESH_RT_SESSION=$(date +%s)
#export __RESH_RT_SESS_SINCE_BOOT=$(cat /proc/uptime | cut -d' ' -f1)
export __RESH_RT_SESS_SINCE_BOOT=1
export __RESH_SESSION_ID=$(cat /proc/sys/kernel/random/uuid)

if [ $(uname) = "Darwin" ]; then
    export __RESH_OS_RELEASE_ID="macos"
    export __RESH_OS_RELEASE_VERSION_ID=$(sw_vers -productVersion 2>/dev/null)
    export __RESH_OS_RELEASE_NAME="macOS"
    export __RESH_OS_RELEASE_PRETTY_NAME="Mac OS X"
    export __RESH_WINDOWID="1" # session 
else
    export __RESH_OS_RELEASE_ID=$(source /etc/os-release; echo $ID)
    export __RESH_OS_RELEASE_VERSION_ID=$(source /etc/os-release; echo $VERSION_ID)
    export __RESH_OS_RELEASE_ID_LIKE=$(source /etc/os-release; echo $ID_LIKE)
    export __RESH_OS_RELEASE_NAME=$(source /etc/os-release; echo $NAME)
    export __RESH_OS_RELEASE_PRETTY_NAME=$(source /etc/os-release; echo $PRETTY_NAME)
fi

nohup resh-daemon &>/dev/null & disown

__resh_preexec() {
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
    __RESH_SHELL_ENV="$SHELL"
    __RESH_TERM="$TERM"
    
    # non-posix
    __RESH_SHELL_PID="$$" # pid (subshells don't affect it)
#    __RESH_WINDOWID="$WINDOWID" # session 
    __RESH_OSTYPE="$OSTYPE" 
    __RESH_MACHTYPE="$MACHTYPE"
    __RESH_SHLVL="$SHLVL"
    __RESH_GIT_CDUP="$(git rev-parse --show-cdup 2>/dev/null)"
    __RESH_GIT_CDUP_EXIT_CODE=$?
    __RESH_GIT_REMOTE="$(git remote get-url origin 2>/dev/null)"
    __RESH_GIT_REMOTE_EXIT_CODE=$?
    #__RESH_GIT_TOPLEVEL="$(git rev-parse --show-toplevel)"
    #__RESH_GIT_TOPLEVEL_EXIT_CODE=$?

    if [ -n "$ZSH_VERSION" ]; then
        # assume Zsh
        __RESH_PID="$$"
        __RESH_HOST="$HOST"
        __RESH_HOSTTYPE="$CPUTYPE"
        __RESH_SHELL="zsh"
    elif [ -n "$BASH_VERSION" ]; then
        # assume Bash
        __RESH_PID="$BASHPID" # current pid
        __RESH_HOST="$HOSTNAME"
        __RESH_HOSTTYPE="$HOSTTYPE"
        __RESH_SHELL="bash"
    else
        # asume something else
        echo "resh PANIC unrecognized shell"
    fi
    # time
    #__RESH_TZ_BEFORE=$(date +%:z)
    __RESH_TZ_BEFORE=2:00
    # __RESH_RT_BEFORE="$EPOCHREALTIME"
    __RESH_RT_BEFORE="$(date +%s)"

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
    # __RESH_RT_AFTER=$EPOCHREALTIME
    __RESH_RT_AFTER="$(date +%s.0)"
    #__RESH_TZ_AFTER=$(date +%:z)
    __RESH_TZ_AFTER=2:00
    __RESH_PWD_AFTER="$PWD"
    if [ ! -z ${__RESH_COLLECT+x} ]; then
        resh-collect -cmdLine "$__RESH_CMDLINE" -exitCode "$__RESH_EXIT_CODE" \
                     -shell "$__RESH_SHELL" \
                     -cols "$__RESH_COLS" \
                     -home "$__RESH_HOME" \
                     -lang "$__RESH_LANG" \
                     -lcAll "$__RESH_LC_ALL" \
                     -lines "$__RESH_LINES" \
                     -login "$__RESH_LOGIN" \
                     -path "$__RESH_PATH" \
                     -pwd "$__RESH_PWD" \
                     -pwdAfter "$__RESH_PWD_AFTER" \
                     -shellEnv "$__RESH_SHELL_ENV" \
                     -term "$__RESH_TERM" \
                     -pid "$__RESH_PID" -shellPid "$__RESH_SHELL_PID" \
                     -windowId "$__RESH_WINDOWID" \
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
                     -osReleasePrettyName "$__RESH_OS_RELEASE_PRETTY_NAME"
    fi
    unset __RESH_COLLECT
}

preexec_functions+=(__resh_preexec)
precmd_functions+=(__resh_precmd)
