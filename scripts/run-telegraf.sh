#!/bin/bash

set -e

cd metrics/telegraf/

# Retrieve newly minted token for influxdb operations
influx_token=$(grep -E '^\s*token\s*=' "${HOME}/.influxdbv2/configs" | awk '{printf $3}' | tr -d '"')

echo "telegraf: starting telegraf with metrics/telegraf/telegraf.conf"
sed -i -e 's/token = \".*\"/token = "'"$influx_token"'"/' telegraf.conf

# Start telegraph with config file and update the token value for the new influxdb instance
telegraf --config telegraf.conf --once --debug 1>/tmp/telegraf.log &
telegraf_pid=$!

# Check if InfluxDB started successfully
if [ -n "$telegraf_pid" ]; then
    echo "Telegraf agent has been started with PID: $telegraf_pid"
    echo "$telegraf_pid" > "${HOME}"/.influxdbv2/telegraf_pid
else
    echo "Failed to start Telegraf agent."
fi

# Add a signal handler to gracefully stop InfluxDB when the script is terminated
cleanup() {
    if [ -n "$telegraf_pid" ]; then
        echo "Stopping Telegraf agent..."
        kill -TERM "$telegraf_pid"
        wait "$telegraf_pid"
        echo "Telegraf agent has been stopped."
    fi
}
trap cleanup SIGINT SIGTERM

echo "telegraf run script complete"