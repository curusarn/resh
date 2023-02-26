#!/usr/bin/env sh
set -eu

q=0
[ "${1-}" != "-q" ] || q=1

xdg_pid() {
    local path="${XDG_DATA_HOME-}"/resh/daemon.pid
    [ -n "${XDG_DATA_HOME-}" ] && [ -f "$path" ] || return 1
    cat "$path"
    rm "$path"
}
default_pid() {
    local path=~/.local/share/resh/daemon.pid
    [ -f "$path" ] || return 1
    cat "$path"
    rm "$path"
}

kill_by_pid() {
    [ -n "$1" ] || return 1
    [ "$q" = "1" ] || printf "Stopping RESH daemon ... (PID: %s)\n" "$1"
    kill "$1"
}
kill_by_name() {
    [ "$q" = "1" ] || printf "Stopping RESH daemon ...\n"
    killall -q resh-daemon
}
failed_to_kill() {
    [ "$q" = "1" ] || echo "Failed to kill the RESH daemon - it probably isn't running"
    return 1
}

kill_by_pid "$(xdg_pid)" || kill_by_pid "$(default_pid)" || kill_by_name || failed_to_kill
