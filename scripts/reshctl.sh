
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

__resh_bind_control_R() {
    # TODO
    echo "bindfunc __resh_widget_control_R_compat"
    return 0
}

__resh_unbind_arrows() {
    if [ "${__RESH_arrow_keys_bind_enabled-0}" != 1 ]; then
        echo "Error: Can't disable arrow key bindings because they are not enabled!"
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
    # TODO
    echo "\ bindfunc __resh_widget_control_R_compat"
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

# wrapper for resh-cli
# meant to be launched on ctrl+R
resh() {
    if resh-cli --sessionID "$__RESH_SESSION_ID" --pwd "$PWD" > ~/.resh/cli_last_run_out.txt 2>&1; then
        # insert on cmdline
        cat ~/.resh/cli_last_run_out.txt
        eval "$(cat ~/.resh/cli_last_run_out.txt)"
        # TODO: get rid of eval
    else
        # print errors
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
    # 100)
    #     # enable all
    #     __resh_bind_all
    #     return 0
    #     ;;
    101)
        # enable arrow keys
        __resh_bind_arrows
        return 0
        ;;
    # disable
    # 110)
    #     # disable all
    #     __resh_unbind_all
    #     return 0
    #     ;;
    111)
        # disable arrow keys
        __resh_unbind_arrows
        return 0
        ;;
    200)
        # reload rc files
        . ~/.resh/shellrc
        return 0
        ;;
    201)
        # inspect session history 
        # reshctl debug inspect N
        resh-inspect --sessionID "$__RESH_SESSION_ID" --count "${3-10}"
        return 0
        ;;
    202)
        # show status 
        if [ "${__RESH_arrow_keys_bind_enabled-0}" != 0 ]; then
            echo " * this session: ENABLED"
        else
            echo " * this session: DISABLED"
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