
# shellcheck source=hooks.sh
. ~/.resh/hooks.sh

__resh_helper_arrow_pre() {
    # set prefix
    __RESH_PREFIX=${BUFFER:0:CURSOR}
    # cursor not at the end of the line => end "NO_PREFIX_MODE"
    [ "$CURSOR" -ne "${#BUFFER}" ] && __RESH_HIST_NO_PREFIX_MODE=0
    # if user made any edits from last recall action => restart histno AND deactivate "NO_PREFIX_MODE"
    [ "$BUFFER" != "$__RESH_HIST_PREV_LINE" ] && __RESH_HISTNO=0 && __RESH_HIST_NO_PREFIX_MODE=0
    # "NO_PREFIX_MODE" => set prefix to empty string
    [ "$__RESH_HIST_NO_PREFIX_MODE" -eq 1 ] && __RESH_PREFIX=""
    # append curent recall action
    __RESH_HIST_RECALL_ACTIONS="$__RESH_HIST_RECALL_ACTIONS;arrow_up:$__RESH_PREFIX"
    # histno == 0 => save current line
    [ $__RESH_HISTNO -eq 0 ] && __RESH_HISTNO_ZERO_LINE=$BUFFER
}
__resh_helper_arrow_post() {
    # cursor at the beginning of the line => activate "NO_PREFIX_MODE"
    [ "$CURSOR" -eq 0 ] && __RESH_HIST_NO_PREFIX_MODE=1
    # "NO_PREFIX_MODE" => move cursor to the end of the line
    [ "$__RESH_HIST_NO_PREFIX_MODE" -eq 1 ] && CURSOR=${#BUFFER}
    # save current line so we can spot user edits next time
    __RESH_HIST_PREV_LINE=$BUFFER
}

__resh_widget_arrow_up() {
    # run helper function
    __resh_helper_arrow_pre
    # increment histno
    (( __RESH_HISTNO++ ))
    # back at histno == 0 => restore original line
    if [ $__RESH_HISTNO -eq 0 ]; then
        BUFFER=$__RESH_HISTNO_ZERO_LINE
    else
        # run recall
        local NEW_BUFFER
        NEW_BUFFER="$(__resh_collect --recall --histno "$__RESH_HISTNO" --prefix-search "$__RESH_PREFIX" 2> ~/.resh/arrow_up_last_run_out.txt)"
        # IF new buffer in non-empty THEN use the new buffer ELSE revert histno change
        # shellcheck disable=SC2015
        [ ${#NEW_BUFFER} -gt 0 ] && BUFFER=$NEW_BUFFER || (( __RESH_HISTNO-- ))
    fi
    # run post helper
    __resh_helper_arrow_post
}
__resh_widget_arrow_down() {
    # run helper function
    __resh_helper_arrow_pre
    # increment histno
    (( __RESH_HISTNO-- ))
    # prevent HISTNO from getting negative (for now)
    [ $__RESH_HISTNO -lt 0 ] && __RESH_HISTNO=0
    # back at histno == 0 => restore original line
    if [ $__RESH_HISTNO -eq 0 ]; then
        BUFFER=$__RESH_HISTNO_ZERO_LINE
    else
        # run recall
        local NEW_BUFFER
        NEW_BUFFER="$(__resh_collect --recall --histno "$__RESH_HISTNO" --prefix-search "$__RESH_PREFIX" 2> ~/.resh/arrow_down_last_run_out.txt)"
        # IF new buffer in non-empty THEN use the new buffer ELSE revert histno change
        # shellcheck disable=SC2015
        [ ${#NEW_BUFFER} -gt 0 ] && BUFFER=$NEW_BUFFER || (( __RESH_HISTNO++ ))
    fi
    __resh_helper_arrow_post
}
__resh_widget_control_R() {
    local __RESH_LBUFFER=${BUFFER:0:CURSOR}
    __RESH_HIST_RECALL_ACTIONS="$__RESH_HIST_RECALL_ACTIONS;control_R:$__RESH_LBUFFER"
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
