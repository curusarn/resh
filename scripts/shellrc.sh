#!/hint/sh

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

if [ -n "${ZSH_VERSION-}" ]; then
    # shellcheck disable=SC1009
    __RESH_SHELL="zsh"
    __RESH_HOST="$HOST"
    __RESH_HOSTTYPE="$CPUTYPE"
    __resh_zsh_completion_init
elif [ -n "${BASH_VERSION-}" ]; then
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

# shellcheck disable=2155
export __RESH_VERSION=$(resh-collect -version)
# shellcheck disable=2155
export __RESH_REVISION=$(resh-collect -revision)

# FIXME: this does not exist anymore
# __resh_set_xdg_home_paths

__resh_run_daemon

[ "$(resh-config --key BindControlR)" = true ] && __resh_bind_control_R

# block for anything we only want to do once per session
# NOTE: nested shells are still the same session
if [ -z "${__RESH_SESSION_ID+x}" ]; then
    export __RESH_SESSION_ID; __RESH_SESSION_ID=$(__resh_get_uuid)
    export __RESH_SESSION_PID="$$"
    # TODO add sesson time
    __resh_reset_variables
    __resh_session_init
fi

# block for anything we only want to do once per shell
if [ -z "${__RESH_INIT_DONE+x}" ]; then
    preexec_functions+=(__resh_preexec)
    precmd_functions+=(__resh_precmd)

    __resh_reset_variables

    __RESH_INIT_DONE=1
fi