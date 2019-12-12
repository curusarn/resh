
# shellcheck source=../submodules/bash-zsh-compat-widgets/bindfunc.sh
. ~/.resh/bindfunc.sh
# shellcheck source=widgets.sh
. ~/.resh/widgets.sh

__resh_bind_arrows() {
    if [ "${__RESH_arrow_keys_bind_enabled-0}" != 0 ]; then
        echo "Error: RESH arrow key bindings are already enabled!"
        return 1 
    fi
    bindfunc --revert '\e[A' __resh_widget_arrow_up_compat
    __RESH_bindfunc_revert_arrow_up_bind=$_bindfunc_revert
    bindfunc --revert '\e[B' __resh_widget_arrow_down_compat
    __RESH_bindfunc_revert_arrow_down_bind=$_bindfunc_revert
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
        echo "RESH arrow up binding successfully disabled ✓"
    fi
    if [ -z "${__RESH_bindfunc_revert_arrow_down_bind+x}" ]; then
        echo "Warn: Couldn't revert arrow DOWN binding because 'revert command' is empty."
    else
        eval "$__RESH_bindfunc_revert_arrow_down_bind"
        echo "RESH arrow down binding successfully disabled ✓"
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

reshctl() {
    # run resh-control aka the real reshctl
    resh-control "$@"
    # modify current shell session based on exit status
    local _status=$?
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
    *)
        echo "reshctl() FATAL ERROR: unknown status" >&2
        return "$_status"
        ;;
    esac
}