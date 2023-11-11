#!/bin/bash

set -e

WORKING_DIR=metrics/telegraf

# Retrieve newly minted auth token for influxdb operations
influx_token=$(grep -E '^\s*token\s*=' "${HOME}/.influxdbv2/configs" | awk '{printf $3}' | tr -d '"')

# update telegraf.conf with token and hostname from bundle capture
echo "telegraf: starting telegraf with metrics/telegraf/telegraf.conf"
sed -i -e 's/token = \".*\"/token = "'"$influx_token"'"/' "${WORKING_DIR}"/telegraf.conf

# Start telegraph with config file and update the token value for the new influxdb instance
telegraf --config "${WORKING_DIR}"/telegraf.conf --once --debug >/tmp/telegraf.log 2>&1 &
telegraf_pid=$!

# Check if InfluxDB started successfully
if [ -n "$telegraf_pid" ]; then
    echo "telegraf agent has been started with PID: $telegraf_pid"
    echo "$telegraf_pid" > "${HOME}"/.influxdbv2/telegraf_pid
else
    echo "failed to start Telegraf agent."
fi

# Add a signal handler to gracefully stop InfluxDB when the script is terminated
cleanup() {
    if [ -n "$telegraf_pid" ]; then
        echo "stopping Telegraf agent..."
        kill -TERM "$telegraf_pid"
        wait "$telegraf_pid"
        echo "telegraf agent has been stopped."
    fi
}
trap cleanup SIGINT SIGTERM

echo "telegraf: init complete and metrics ingestion started"
echo "    ==> observe progress at '/tmp/telegraf.log'"
echo "    ==> visit http://localhost:8086 (influxdb ui) to explore metrics"
echo "        ==> un: consul | pw: hashicorp"