#!/bin/bash

set -e
# 4Gi: 4294967296
# 2Gi: 2147483648
# 1Gi: 1073741824

# Implement influxdb tsi1 disk based db vice in-mem (large metrics input)
INFLUXDB_DATA_MAX_INDEX_LOG_FILE_SIZE="1g" \
  INFLUXDB_DATA_SERIES_ID_SET_CACHE_SIZE=0 \
  INFLUXDB_DATA_INDEX_VERSION=tsi1 \
  INFLUXDB_DATA_MAX_INDEX_LOG_FILE_SIZE=1g \
  INFLUXDB_DATA_SERIES_ID_SET_CACHE_SIZE=0 \
  influxd run \
  --http-bind-address 127.0.0.1:8086 \
  --log-level debug \
  --bolt-path "${HOME}"/.influxdbv2/influxd.bolt \
  --engine-path "${HOME}"/.influxdbv2/engine \
  1>/tmp/influxdb.log &

# Track pid for cleanup (if necessary)
influxd_pid=$!

# Check if InfluxDB started successfully
if [ -n "$influxd_pid" ]; then
    echo "influxDB has been started with PID: $influxd_pid"
    echo "$influxd_pid" > "${HOME}"/.influxdbv2/pid
else
    echo "failed to start InfluxDB."
fi

# Add a signal handler to gracefully stop InfluxDB when the script is terminated
cleanup() {
    if [ -n "$influxd_pid" ]; then
        echo "stopping InfluxDB..."
        kill -TERM "$influxd_pid"
        wait "$influxd_pid"
        echo "influxDB has been stopped."
    fi
}
trap cleanup SIGINT SIGTERM
sleep 5
echo "influxdb init complete"