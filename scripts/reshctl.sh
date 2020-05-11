
# shellcheck source=../submodules/bash-zsh-compat-widgets/bindfunc.sh
. ~/.resh/bindfunc.sh
# shellcheck source=widgets.sh
. ~/.resh/widgets.sh

__resh_bind_arrows() {
    if [ "${__RESH_arrow_keys_bind_enabled-0}" != 0 ]; then
        echo "RESH arrow key bindings are already enabled!"
        return 1 
    fi
    bindfunc --revert '\eOA' __resh_widget_arrow_up_compat
    __RESH_bindfunc_revert_arrow_up_bind=$_bindfunc_revert
    bindfunc --revert '\e[A' __resh_widget_arrow_up_compat
    __RESH_bindfunc_revert_arrow_up_bind_vim=$_bindfunc_revert
    bindfunc --vim-cmd --revert 'k' __resh_widget_arrow_up_compat
    __RESH_bindfunc_revert_k_bind_vim=$_bindfunc_revert
    bindfunc --revert '\eOB' __resh_widget_arrow_down_compat
    __RESH_bindfunc_revert_arrow_down_bind=$_bindfunc_revert
    bindfunc --revert '\e[B' __resh_widget_arrow_down_compat
    __RESH_bindfunc_revert_arrow_down_bind_vim=$_bindfunc_revert
    bindfunc --vim-cmd --revert 'j' __resh_widget_arrow_down_compat
    __RESH_bindfunc_revert_j_bind_vim=$_bindfunc_revert
    __RESH_arrow_keys_bind_enabled=1
    return 0
}

__resh_nop() {
    # does nothing
    true
}

__resh_bind_control_R() {
    if [ "${__RESH_control_R_bind_enabled-0}" != 0 ]; then
        echo "RESH SEARCH app Ctrl+R binding is already enabled!"
        return 1 
    fi
    bindfunc --revert '\C-r' __resh_widget_control_R_compat
    __RESH_bindfunc_revert_control_R_bind=$_bindfunc_revert
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

__resh_unbind_arrows() {
    if [ "${__RESH_arrow_keys_bind_enabled-0}" != 1 ]; then
        echo "RESH arrow key bindings are already disabled!"
        return 1 
    fi

    if [ -z "${__RESH_bindfunc_revert_arrow_up_bind+x}" ]; then
        echo "Warn: Couldn't revert arrow UP binding because 'revert command' is empty."
    else
        eval "$__RESH_bindfunc_revert_arrow_up_bind"
        [ -z "${__RESH_bindfunc_revert_arrow_up_bind_vim+x}" ] || eval "$__RESH_bindfunc_revert_arrow_up_bind_vim"
        [ -z "${__RESH_bindfunc_revert_k_bind_vim+x}" ] || eval "$__RESH_bindfunc_revert_k_bind_vim"
        echo "RESH arrow up binding successfully disabled"
        __RESH_arrow_keys_bind_enabled=0
    fi

    if [ -z "${__RESH_bindfunc_revert_arrow_down_bind+x}" ]; then
        echo "Warn: Couldn't revert arrow DOWN binding because 'revert command' is empty."
    else
        eval "$__RESH_bindfunc_revert_arrow_down_bind"
        [ -z "${__RESH_bindfunc_revert_arrow_down_bind_vim+x}" ] || eval "$__RESH_bindfunc_revert_arrow_down_bind_vim"
        [ -z "${__RESH_bindfunc_revert_j_bind_vim+x}" ] || eval "$__RESH_bindfunc_revert_j_bind_vim"
        echo "RESH arrow down binding successfully disabled"
        __RESH_arrow_keys_bind_enabled=0
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
    __resh_bind_arrows
    __resh_bind_control_R
}

__resh_unbind_all() {
    __resh_unbind_arrows
    __resh_unbind_control_R
}

# wrapper for resh-cli for calling resh directly
resh() {
    local buffer
    local git_remote; git_remote="$(git remote get-url origin 2>/dev/null)"
    buffer=$(resh-cli --sessionID "$__RESH_SESSION_ID" --host "$__RESH_HOST" --pwd "$PWD" --gitOriginRemote "$git_remote")
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
        echo "$buffer" >| ~/.resh/cli_last_run_out.txt
        echo "resh-cli ERROR:"
        cat ~/.resh/cli_last_run_out.txt
    fi
}

reshctl() {
    # local log=~/.resh/reshctl.log
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
    31)
        # enable arrow keys
        __resh_bind_arrows
        return 0
        ;;
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
    41)
        # disable arrow keys
        __resh_unbind_arrows
        return 0
        ;;
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
        if [ "${__RESH_arrow_keys_bind_enabled-0}" != 0 ]; then
            echo ' * this session: ENABLED'
        else
            echo ' * this session: DISABLED'
        fi
        echo
		echo 'Control R binding ...'
        if [ "$(resh-config --key BindControlR)" = true ]; then
			echo ' * future sessions: ENABLED (experimental)'
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