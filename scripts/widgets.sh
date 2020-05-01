
# shellcheck source=hooks.sh
. ~/.resh/hooks.sh

__resh_helper_arrow_pre() {
    # this is a very bad workaround
    # force bash-preexec to run repeatedly because otherwise premature run of bash-preexec overshadows the next poper run
    # I honestly think that it's impossible to make widgets work in bash without hacks like this
    # shellcheck disable=2034
    __bp_preexec_interactive_mode="on"
    # set recall strategy
    __RESH_HIST_RECALL_STRATEGY="bash_recent - history-search-{backward,forward}"
    # set prefix
    __RESH_PREFIX=${BUFFER:0:$CURSOR}
    # cursor not at the end of the line => end "NO_PREFIX_MODE"
    [ "$CURSOR" -ne "${#BUFFER}" ] && __RESH_HIST_NO_PREFIX_MODE=0
    # if user moved the cursor or made edits (to prefix) from last recall action
    # => restart histno AND deactivate "NO_PREFIX_MODE" AND clear end of recall list histno
    [ "$__RESH_PREFIX" != "$__RESH_HIST_PREV_PREFIX" ] && __RESH_HISTNO=0 && __RESH_HIST_NO_PREFIX_MODE=0 && __RESH_HISTNO_MAX=""
    # "NO_PREFIX_MODE" => set prefix to empty string
    [ "$__RESH_HIST_NO_PREFIX_MODE" -eq 1 ] && __RESH_PREFIX=""
    # histno == 0 => save current line
    [ "$__RESH_HISTNO" -eq 0 ] && __RESH_HISTNO_ZERO_LINE=$BUFFER
}
__resh_helper_arrow_post() {
    # cursor at the beginning of the line => activate "NO_PREFIX_MODE"
    [ "$CURSOR" -eq 0 ] && __RESH_HIST_NO_PREFIX_MODE=1
    # "NO_PREFIX_MODE" => move cursor to the end of the line
    [ "$__RESH_HIST_NO_PREFIX_MODE" -eq 1 ] && CURSOR=${#BUFFER}
    # save current prefix so we can spot when user moves cursor or edits (the prefix part of) the line
    __RESH_HIST_PREV_PREFIX=${BUFFER:0:$CURSOR}
    # recorded to history
    __RESH_HIST_PREV_LINE=${BUFFER}
}

__resh_widget_arrow_up() {
    # run helper function
    __resh_helper_arrow_pre
    # append curent recall action
    __RESH_HIST_RECALL_ACTIONS="$__RESH_HIST_RECALL_ACTIONS|||arrow_up:$__RESH_PREFIX"
    # increment histno
    __RESH_HISTNO=$((__RESH_HISTNO+1))
    if [ "${#__RESH_HISTNO_MAX}" -gt 0 ] && [ "${__RESH_HISTNO}" -gt "${__RESH_HISTNO_MAX}" ]; then
        # end of the recall list -> don't recall, do nothing
        # fix histno
        __RESH_HISTNO=$((__RESH_HISTNO-1))
    elif [ "$__RESH_HISTNO" -eq 0 ]; then
        # back at histno == 0 => restore original line
        BUFFER=$__RESH_HISTNO_ZERO_LINE
    else
        # run recall
        local NEW_BUFFER
        local status_code
        NEW_BUFFER="$(__resh_collect --recall --prefix-search "$__RESH_PREFIX" 2>| ~/.resh/arrow_up_last_run_out.txt)"
        status_code=$?
        # revert histno change on error
        # shellcheck disable=SC2015
        if [ "${status_code}" -eq 0 ]; then
            BUFFER=$NEW_BUFFER
        else
            __RESH_HISTNO=$((__RESH_HISTNO-1))
            __RESH_HISTNO_MAX=$__RESH_HISTNO
        fi
    fi
    # run post helper
    __resh_helper_arrow_post
}
__resh_widget_arrow_down() {
    # run helper function
    __resh_helper_arrow_pre
    # append curent recall action
    __RESH_HIST_RECALL_ACTIONS="$__RESH_HIST_RECALL_ACTIONS|||arrow_down:$__RESH_PREFIX"
    # increment histno
    __RESH_HISTNO=$((__RESH_HISTNO-1))
    # prevent HISTNO from getting negative (for now)
    [ "$__RESH_HISTNO" -lt 0 ] && __RESH_HISTNO=0
    # back at histno == 0 => restore original line
    if [ "$__RESH_HISTNO" -eq 0 ]; then
        BUFFER=$__RESH_HISTNO_ZERO_LINE
    else
        # run recall
        local NEW_BUFFER
        NEW_BUFFER="$(__resh_collect --recall --prefix-search "$__RESH_PREFIX" 2>| ~/.resh/arrow_down_last_run_out.txt)"
        # IF new buffer in non-empty THEN use the new buffer ELSE revert histno change
        # shellcheck disable=SC2015
        [ "${#NEW_BUFFER}" -gt 0 ] && BUFFER=$NEW_BUFFER || (( __RESH_HISTNO++ ))
    fi
    __resh_helper_arrow_post
}
__resh_widget_control_R() {
    # this is a very bad workaround
    # force bash-preexec to run repeatedly because otherwise premature run of bash-preexec overshadows the next poper run
    # I honestly think that it's impossible to make widgets work in bash without hacks like this
    # shellcheck disable=2034
    __bp_preexec_interactive_mode="on"

    # local __RESH_PREFIX=${BUFFER:0:CURSOR}
    # __RESH_HIST_RECALL_ACTIONS="$__RESH_HIST_RECALL_ACTIONS;control_R:$__RESH_PREFIX"
    local PREVBUFFER=$BUFFER
    __RESH_HIST_RECALL_ACTIONS="$__RESH_HIST_RECALL_ACTIONS|||control_R:$BUFFER"

    local status_code
    local git_remote; git_remote="$(git remote get-url origin 2>/dev/null)"
    BUFFER=$(resh-cli --sessionID "$__RESH_SESSION_ID" --host "$HOST" --pwd "$PWD" --gitOriginRemote "$git_remote" --query "$BUFFER")
    status_code=$?
    if [ $status_code = 111 ]; then
        # execute
        if [ -n "${ZSH_VERSION-}" ]; then
            # zsh
            zle accept-line
        elif [ -n "${BASH_VERSION-}" ]; then
            # bash
            # set chained keyseq to accept-line
            bind '"\u[32~": accept-line'
        fi
    elif [ $status_code = 0 ]; then
        if [ -n "${BASH_VERSION-}" ]; then
            # bash
            # set chained keyseq to nothing
            bind -x '"\u[32~": __resh_nop'
        fi
    elif [ $status_code = 130 ]; then
        BUFFER="$PREVBUFFER"
    else
        echo "$BUFFER" >| ~/.resh/cli_last_run_out.txt
        echo "# RESH cli failed - sorry for the inconvinience (error output was saved to ~/.resh/cli_last_run_out.txt)" 
        BUFFER="$PREVBUFFER"
    fi
    CURSOR=${#BUFFER}
    # recorded to history
    __RESH_HIST_PREV_LINE=${BUFFER}
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
