#!/usr/bin/env sh

failed_to_kill() {
    [ "${1-}" != "-q" ] && echo "Failed to kill the RESH daemon - it probably isn't running"
}

xdg_pid() {
    local path="${XDG_DATA_HOME-}"/resh/daemon.pid
    [ -n "${XDG_DATA_HOME-}" ] && [ -f "$path" ] || return 1
    cat "$path"
}
default_pid() {
    local path=~/.local/share/resh/daemon.pid
    [ -f "$path" ] || return 1
    cat "$path"
}
legacy_pid() {
    local path=~/.resh/resh.pid
    [ -f "$path" ] || return 1
    cat "$path"
}
pid=$(xdg_pid || default_pid || legacy_pid)

if [ -n "$pid" ]; then
    [ "${1-}" != "-q" ] && printf "Stopping RESH daemon ... (PID: %s)\n" "$pid"
    kill "$pid" || failed_to_kill
else
    [ "${1-}" != "-q" ] && printf "Stopping RESH daemon ...\n"
    killall -q resh-daemon || failed_to_kill
fi

