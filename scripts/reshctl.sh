
# source bindutil - contains functions to bind and unbind RESH widgets
# shellcheck source=bindutil.sh
. ~/.resh/bindutil.sh

reshctl() {
    # run resh-control aka the real reshctl
    resh-control "$@"
    # modify current shell session based on exit status
    local status=$?
    case "$status" in
    0|1)
        # success | fail
        return "$status"
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
    *)
        echo "reshctl() FATAL ERROR: unknown status" >&2
        return "$status"
        ;;
    esac
}