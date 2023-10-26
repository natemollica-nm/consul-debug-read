#!/bin/bash

set -e

pid_file="${HOME}"/.influxdbv2/pid
influxd_pid=

if test -f "$pid_file"; then
  influxd_pid="$(cat "$pid_file")"
else
  exit 0
fi

if [ -n "$influxd_pid" ]; then
    echo "stopping influxDB..."
    while ps -p "$influxd_pid" >/dev/null; do
          kill -INT "$influxd_pid" >/dev/null
          sleep 1
    done
    echo "influxDB has been stopped."
fi