
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

__resh_bind_all() {
    __resh_bind_control_R
}

__resh_unbind_all() {
    __resh_unbind_control_R
}

# wrapper for resh-cli for calling resh directly
resh() {
    local buffer
    local git_remote; git_remote="$(git remote get-url origin 2>/dev/null)"
    buffer=$(resh-cli --sessionID "$__RESH_SESSION_ID" --host "$__RESH_HOST" --pwd "$PWD" --gitOriginRemote "$git_remote" "$@")
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
        local fpath_last_run="$__RESH_XDG_CACHE_HOME/cli_last_run_out.txt"
        echo "$buffer" >| "$fpath_last_run"
        echo "resh-cli failed - check '$fpath_last_run' and '~/.resh/cli.log'"
    fi
}

reshctl() {
    # export current shell because resh-control needs to know
    export __RESH_ctl_shell=$__RESH_SHELL
    # run resh-control aka the real reshctl
    resh-control "$@"

    # modify current shell session based on exit status
    local _status=$?
    # echo $_status
    # unexport current shell  
    unset __RESH_ctl_shell
    case "$_status" in
    0|1)
        # success | fail
        return "$_status"
        ;;
    # enable
    # 30)
    #     # enable all
    #     __resh_bind_all
    #     return 0
    #     ;;
    32)
        # enable control R
        __resh_bind_control_R
        return 0
        ;;
    # disable
    # 40)
    #     # disable all
    #     __resh_unbind_all
    #     return 0
    #     ;;
    42)
        # disable control R
        __resh_unbind_control_R
        return 0
        ;;
    50)
        # reload rc files
        . ~/.resh/shellrc
        return 0
        ;;
    51)
        # inspect session history 
        # reshctl debug inspect N
        resh-inspect --sessionID "$__RESH_SESSION_ID" --count "${3-10}"
        return 0
        ;;
    52)
        # show status 
        echo
		echo 'Control R binding ...'
        if [ "$(resh-config --key BindControlR)" = true ]; then
			echo ' * future sessions: ENABLED'
		else
			echo ' * future sessions: DISABLED'
        fi
        if [ "${__RESH_control_R_bind_enabled-0}" != 0 ]; then
            echo ' * this session: ENABLED'
        else
            echo ' * this session: DISABLED'
        fi
        return 0
        ;;
    *)
        echo "reshctl() FATAL ERROR: unknown status ($_status)" >&2
        echo "Possibly caused by version mismatch between installed resh and resh in this session." >&2
        echo "Please REPORT this issue here: https://github.com/curusarn/resh/issues" >&2
        echo "Please RESTART your terminal window." >&2
        return "$_status"
        ;;
    esac
}
