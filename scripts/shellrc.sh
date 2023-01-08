#!/hint/sh

PATH=$PATH:~/.resh/bin

# shellcheck source=hooks.sh
. ~/.resh/hooks.sh

if [ -n "${ZSH_VERSION-}" ]; then
    # shellcheck disable=SC1009
    __RESH_SHELL="zsh"
elif [ -n "${BASH_VERSION-}" ]; then
    __RESH_SHELL="bash"
else
    echo "RESH PANIC: unrecognized shell - please report this to https://github.com/curusarn/resh/issues"
fi

# shellcheck disable=2155
export __RESH_VERSION=$(resh-collect -version)

resh-daemon-start

[ "$(resh-config --key BindControlR)" = true ] && __resh_bind_control_R

# block for anything we only want to do once per session
# NOTE: nested shells are still the same session
#       i.e. $__RESH_SESSION_ID will be set in nested shells
if [ -z "${__RESH_SESSION_ID+x}" ]; then
    # shellcheck disable=2155
    export __RESH_SESSION_ID=$(resh-generate-uuid)

    __resh_session_init
fi

# block for anything we only want to do once per shell
# NOTE: nested shells are new shells
#       i.e. $__RESH_INIT_DONE will NOT be set in nested shells
if [ -z "${__RESH_INIT_DONE+x}" ]; then
    preexec_functions+=(__resh_preexec)
    precmd_functions+=(__resh_precmd)

    __RESH_INIT_DONE=1
fi