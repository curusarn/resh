
PATH=$PATH:~/.resh/bin
# if [ -n "$ZSH_VERSION" ]; then
#     zmodload zsh/datetime
# fi

# shellcheck source=hooks.sh
. ~/.resh/hooks.sh
# shellcheck source=util.sh
. ~/.resh/util.sh
# shellcheck source=reshctl.sh
. ~/.resh/reshctl.sh

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

if [ -z "${__RESH_SESSION_ID+x}" ]; then
    export __RESH_SESSION_ID; __RESH_SESSION_ID=$(__resh_get_uuid)
    export __RESH_SESSION_PID="$$"
    # TODO add sesson time
    __resh_reset_variables
    __resh_session_init
fi

# do not add more hooks when shellrc is sourced again  
if [ -z "${__RESH_PREEXEC_PRECMD_HOOKS_ADDED+x}" ]; then
    preexec_functions+=(__resh_preexec)
    precmd_functions+=(__resh_precmd)
    __RESH_PREEXEC_PRECMD_HOOKS_ADDED=1
fi
