#!/hint/sh

# shellcheck source=../submodules/bash-zsh-compat-widgets/bindfunc.sh
. ~/.resh/bindfunc.sh
# shellcheck source=widgets.sh
. ~/.resh/widgets.sh

__resh_nop() {
    # does nothing
    true
}

__resh_bind_control_R() {
    bindfunc --revert '\C-r' __resh_widget_control_R_compat
    if [ "${__RESH_control_R_bind_enabled-0}" != 0 ]; then
        # Re-binding is a valid usecase but it shouldn't happen much
        # so this is a warning
        echo "Re-binding RESH SEARCH app to Ctrl+R ..."
    else
        # Only save original binding if resh binding was not enabled
        __RESH_bindfunc_revert_control_R_bind=$_bindfunc_revert
    fi
    __RESH_control_R_bind_enabled=1
    if [ -n "${BASH_VERSION-}" ]; then
        # fuck bash
        bind '"\C-r": "\u[31~\u[32~"'
        bind -x '"\u[31~": __resh_widget_control_R_compat'

        # execute
        # bind '"\u[32~": accept-line'

        # just paste
        # bind -x '"\u[32~": __resh_nop'
        true
    fi
    return 0
}

__resh_unbind_control_R() {
    if [ "${__RESH_control_R_bind_enabled-0}" != 1 ]; then
        echo "RESH SEARCH app Ctrl+R binding is already disabled!"
        return 1 
    fi
    if [ -z "${__RESH_bindfunc_revert_control_R_bind+x}" ]; then
        echo "Warn: Couldn't revert Ctrl+R binding because 'revert command' is empty."
    else
        eval "$__RESH_bindfunc_revert_control_R_bind"
    fi
    __RESH_control_R_bind_enabled=0
    return 0
}

# wrapper for resh-cli for calling resh directly
resh() {
    local buffer
    local git_remote; git_remote="$(git remote get-url origin 2>/dev/null)"
    buffer=$(resh-cli --sessionID "$__RESH_SESSION_ID" --pwd "$PWD" --gitOriginRemote "$git_remote" "$@")
    status_code=$?
    if [ $status_code = 111 ]; then
        # execute
        echo "$buffer" 
        eval "$buffer"
    elif [ $status_code = 0 ]; then
        # paste
        echo "$buffer" 
    elif [ $status_code = 130 ]; then
        true
    else
        printf "%s" "$buffer" >&2
    fi
}