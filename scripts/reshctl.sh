
# shellcheck source=../submodules/bash-zsh-compat-widgets/bindfunc.sh
. ~/.resh/bindfunc.sh
# shellcheck source=widgets.sh
. ~/.resh/widgets.sh

__resh_bind_arrows() {
    bindfunc '\e[A' __resh_widget_arrow_up_compat
    bindfunc '\e[B' __resh_widget_arrow_down_compat
    return 0
}

__resh_bind_control_R() {
    echo "bindfunc __resh_widget_control_R_compat"
    return 0
}
__resh_unbind_arrows() {
    echo "\ bindfunc __resh_widget_arrow_up_compat"
    echo "\ bindfunc __resh_widget_arrow_down_compat"
    return 0
}

__resh_unbind_control_R() {
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
    100)
        # enable all
        __resh_bind_all
        return 0
        ;;
    # disable
    110)
        # disable all
        __resh_unbind_all
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