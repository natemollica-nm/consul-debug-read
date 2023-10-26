#!/bin/bash

set -e

pid_file="${HOME}"/.influxdbv2/telegraf_pid
telegraf_pid=

if test -f "$pid_file"; then
  telegraf_pid="$(cat "$pid_file")"
else
  exit 0
fi

if [ -n "$telegraf_pid" ]; then
    count=1
    printf "stopping telegraf..."
    while ps -p "$telegraf_pid" >/dev/null; do
          kill -INT "$telegraf_pid" >/dev/null
          sleep 1
          count+=1
          printf '.'
          if [[ $count == 30 ]]; then
            count=0
            printf "\nforce killing telegraf (30s timeout reached)"
            sudo kill -9 "$telegraf_pid" >/dev/null
            break
          fi
    done
    printf "\ntelegraf has been stopped.\n"
fi