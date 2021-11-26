
# shellcheck source=hooks.sh
. ~/.resh/hooks.sh

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
    BUFFER=$(resh-cli --sessionID "$__RESH_SESSION_ID" --host "$__RESH_HOST" --pwd "$PWD" --gitOriginRemote "$git_remote" --query "$BUFFER")
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
    else
        echo "$BUFFER" >| ~/.resh/cli_last_run_out.txt
        echo "# RESH SEARCH APP failed - sorry for the inconvinience (error output was saved to ~/.resh/cli_last_run_out.txt)" 
        BUFFER="$PREVBUFFER"
    fi
    CURSOR=${#BUFFER}
    # recorded to history
    __RESH_HIST_PREV_LINE=${BUFFER}
}

__resh_widget_control_R_compat() {
   __bindfunc_compat_wrapper __resh_widget_control_R
}
