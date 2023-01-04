#! /usr/bin/env sh
# Run daemon in background - don't block
# Redirect stdin, stdout, and stderr to /dev/null - detach all I/O
resh-daemon </dev/null >/dev/null 2>/dev/null &

# After resh-daemon-start.sh exits the resh-daemon process loses its parent
# and it gets adopted by init

# NOTES:
# No disown - job control of this shell doesn't affect the parent shell 
# No nohup - SIGHUP signals won't be sent to orphaned resh-daemon (plus the daemon ignores them)
# No setsid - SIGINT signals won't be sent to orphaned resh-daemon (plus the daemon ignores them)