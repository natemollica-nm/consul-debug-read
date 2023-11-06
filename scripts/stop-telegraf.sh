#!/bin/bash

set -e

telegraf_pids=()
# Use pgrep to find Telegraf PIDs and add them to the array
while IFS= read -r pid; do
  telegraf_pids+=("$pid")
done < <(pgrep telegraf)

if [ "${#telegraf_pids[@]}" -gt 1 ]; then
  # Gracefully terminate Telegraf processes
  for pid in "${telegraf_pids[@]}"; do
    echo "stopping telegraf process with PID $pid gracefully..."
    kill -INT "$pid"  # Send SIGINT (interrupt signal)
  done

  # Wait for processes to terminate gracefully
  sleep 5

  # Forcefully kill remaining Telegraf processes
  for pid in "${telegraf_pids[@]}"; do
    if kill -0 "$pid" 2>/dev/null; then
      echo "forcefully killing telegraf process with PID $pid"
      kill -9 "$pid"  # Send SIGKILL
    fi
  done
  echo "telegraf stopped"
elif [ "${#telegraf_pids[@]}" -eq 1 ] && [ -n "${telegraf_pids[0]}" ]; then
      echo "stopping telegraf process with PID ${telegraf_pids[0]} gracefully..."
      kill -INT "${telegraf_pids[0]}" >/dev/null
      sleep 5
      if kill -0 "${telegraf_pids[0]}" 2>/dev/null; then
        echo "forcefully killing telegraf process with PID ${telegraf_pids[0]}"
        kill -9 "${telegraf_pids[0]}"  # Send SIGKILL
      fi
      echo "telegraf stopped"
else
  echo "telegraf not running."
fi