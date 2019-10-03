

# shellcheck source=../submodules/bash-zsh-compat-widgets/bindfunc.sh
. ~/.resh/bindfunc.sh
# shellcheck source=widgets.sh
. ~/.resh/widgets.sh

__resh_bind_arrows() {
    echo "bindfunc __resh_widget_arrow_up"
    echo "bindfunc __resh_widget_arrow_down"
    return 0
}

__resh_bind_control_R() {
    echo "bindfunc __resh_widget_control_R"
    return 0
}
__resh_unbind_arrows() {
    echo "\ bindfunc __resh_widget_arrow_up"
    echo "\ bindfunc __resh_widget_arrow_down"
    return 0
}

__resh_unbind_control_R() {
    echo "\ bindfunc __resh_widget_control_R"
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
