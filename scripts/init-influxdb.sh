#!/bin/bash

set -e

influxd run --http-bind-address 127.0.0.1:8086 1>/tmp/influxdb.log &
influxd_pid=$!

# Check if InfluxDB started successfully
if [ -n "$influxd_pid" ]; then
    echo "InfluxDB has been started with PID: $influxd_pid"
    echo "$influxd_pid" > "${HOME}"/.influxdbv2/pid
else
    echo "Failed to start InfluxDB."
fi

# Add a signal handler to gracefully stop InfluxDB when the script is terminated
cleanup() {
    if [ -n "$influxd_pid" ]; then
        echo "Stopping InfluxDB..."
        kill -TERM "$influxd_pid"
        wait "$influxd_pid"
        echo "InfluxDB has been stopped."
    fi
}
trap cleanup SIGINT SIGTERM
sleep 5
echo "influxdb init complete"