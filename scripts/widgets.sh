
# shellcheck source=hooks.sh
. ~/.resh/hooks.sh

__resh_widget_arrow_up() {
    (( __RESH_HISTNO++ ))
    BUFFER="$(__resh_collect --recall --histno "$__RESH_HISTNO" 2> ~/.resh/arrow_up_last_run_out.txt)"
}
__resh_widget_arrow_down() {
    (( __RESH_HISTNO-- ))
    BUFFER="$(__resh_collect --recall --histno "$__RESH_HISTNO" 2> ~/.resh/arrow_down_last_run_out.txt)"
}
__resh_widget_control_R() {
    # resh-collect --hstr
    hstr
}

__resh_widget_arrow_up_compat() {
   __bindfunc_compat_wrapper __resh_widget_arrow_up
}

__resh_widget_arrow_down_compat() {
   __bindfunc_compat_wrapper __resh_widget_arrow_down
}

__resh_widget_control_R_compat() {
   __bindfunc_compat_wrapper __resh_widget_control_R
}
