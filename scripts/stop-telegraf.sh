#!/bin/bash

set -e

telegraf_pid="$(cat "${HOME}"/.influxdbv2/telegraf_pid)"

if [ -n "$telegraf_pid" ]; then
    count=1
    printf "Stopping telegraf..."
    while ps -p "$telegraf_pid" >/dev/null; do
          kill -INT "$telegraf_pid" >/dev/null
          sleep 1
          count+=1
          printf '.'
          if [[ $count == 30 ]]; then
            count=0
            printf "\n force killing telegraf (30s timeout reached)"
            sudo kill -9 "$telegraf_pid" >/dev/null
            break
          fi
    done
    printf "\n telegraf has been stopped."
fi