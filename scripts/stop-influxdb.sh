#!/bin/bash

set -e

influxd_pid="$(cat "${HOME}"/.influxdbv2/pid)"

if [ -n "$influxd_pid" ]; then
    echo "Stopping InfluxDB..."
    while ps -p "$influxd_pid" >/dev/null; do
          kill -INT "$influxd_pid" >/dev/null
          sleep 1
    done
    echo "InfluxDB has been stopped."
fi