#!/usr/bin/env bash

set -e

# TIMESTAMP_REGEX='^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$'

LOG="$1"
METHOD="$2"

# Check if the user provided a path or a file as an argument
if [ $# -eq 0 ]; then
  echo ""
  echo "-----------------------------------------------------"
  echo "-------- Consul RPC Call Rate Log Parser ------------"
  echo "-----------------------------------------------------"
  echo "Requires:"
  echo "    - Valid consul monitor or consul log file with '.log' extension"
  echo "    - TRACE level capture enabled on agent's log or monitor"
  echo "      - agent cmd: -log-level=trace"
  echo "      - agent conf: log_level=\"trace\""
  echo "      - monitor cmd: consul monitor -log-level=trace"
  echo ""
  echo "Description:"
  echo "    Parses consul trace logs for all (default) or specified ([method]) rpc method calls and provides"
  echo "        => Rate-per-minute count of rpc call(s) sorted from highest to lowest"
  echo "        => Total log capture count of rpc call(s) sorted from highest to lowest"
  echo ""
  echo "Usage: "
  echo "    $0 <log> [method]"
  echo ""
  echo "Parameters:"
  echo "  <log>: parameter identifying either"
  echo "        => Directory containing multiple log files (.log extension)"
  echo "        => Single log file name (path to file)."
  echo ""
  echo "Options:"
  echo "  [method]: Specify an RPC method to filter results (e.g., 'Catalog.NodeServiceList') - Optional"

  exit 1
fi


logs=
if [ -d "$LOG" ]; then
  logs=( "$(ls "$LOG"/*.log)" )
elif [ -f "$LOG" ]; then
  logs=("$LOG")
else
  echo "invalid input: $LOG is neither a directory nor a log file."
  exit 1
fi



for log in "${logs[@]}"; do
  log_name="$(basename "$log")"
  echo "---------------- START $log_name ----------------"
  if [ -n "$METHOD" ]; then
      echo "retrieving rpc $METHOD counts/minute: $log"
      # use awk to grep log for rpc method name and show calls/minute
      printf "%-5s | %-30s | %-16s\n" "Count" "Method" "Minute-Interval"
      printf "%-5s-|-%-30s-|-%-16s\n" "-----" "------------------------------" "----------------"
      awk -v log_name="$log_name" -v method_regex="method=$METHOD" -v method="$METHOD" '
        $0 ~ method_regex {
          timestamp = $1 " " $2;
          gsub(/[\[\]]/, "", $3);               # Remove square brackets
          interval = substr(timestamp, 1, 16);  # Extract year, month, day, hour, and minute
          count[interval]++;
        }
        END {
          for (interval in count) {
            printf "%-5d | %-30s | %-16s\n", count[interval], method, interval;
          }
        }
      ' "$log" | sort -n -r
    printf '%s\n%s' "---------------------------------------------------------" "Total: "
    grep rpc_ < "$log" | grep --only-matching "method=$METHOD" | sed 's/method=//' | sort | uniq -c | sort -r -n
    echo ""
  else
      log_name="$(basename "$log")"
      echo "retrieving all rpc method counts/minute: $log"
      # Use awk to process the log file
      printf "%-5s | %-30s | %-16s\n" "Count" "Method" "Minute-Interval"
      printf "%-5s-|-%-30s-|-%-16s\n" "-----" "------------------------------" "----------------"
      awk -v log_name="$log_name" '
        BEGIN {
          method_regex = "rpc_server_call: method=[^ ]*";
        }
        {
          timestamp = $1 " " $2;
          gsub(/[\[\]]/, "", $3);               # Remove square brackets
          interval = substr(timestamp, 1, 16);  # Extract year, month, day, hour, and minute
        }
        $0 ~ method_regex {
          # Extract and count method names
          if (match($0, method_regex)) {
            method = substr($0, RSTART, RLENGTH);
            sub("rpc_server_call: method=", "", method);
            count[interval, method]++;
            methods[interval, method] = 1;  # Store the method name for this interval
          }
        }
        END {
          for (interval_method in methods) {
            split(interval_method, arr, SUBSEP);
            interval = arr[1];
            method = arr[2];
            printf "%-5d | %-30s | %-16s\n", count[interval, method], method, interval;
          }
        }
      ' "$log" | sort -n -r
    printf '%s\n%s\n' "---------------------------------------------------------" "Totals: "
    grep rpc_ < "$log" | grep --only-matching "method=[^ ]*" | sed 's/method=//' | sort | uniq -c | sort -r -n
    echo ""
  fi
  echo "---------------- END $log_name --------------------"
  echo ""
done
