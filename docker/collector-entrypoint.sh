#!/bin/sh

set -eu

if [ -z "${IKNOW_USER_ID:-}" ]; then
    echo "IKNOW_USER_ID is required" >&2
    exit 1
fi

mkdir -p "${IKNOW_DATA_DIR:-/data}/images"

trap 'exit 0' INT TERM

# Run cron.sh at minute 59 of every hour. The container uses UTC by default.
while :; do
    now=$(date +%s)
    next_run=$((now / 3600 * 3600 + 59 * 60))
    if [ "$next_run" -le "$now" ]; then
        next_run=$((next_run + 3600))
    fi

    sleep $((next_run - now)) &
    wait $!

    if ! /app/cron.sh; then
        echo "cron.sh failed; retrying at the next scheduled time" >&2
    fi
done
