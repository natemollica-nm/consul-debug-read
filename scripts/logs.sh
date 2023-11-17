#!/usr/bin/env bash

set -e

# TIMESTAMP_REGEX='^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$'
LOG_PATH=bundles/consul-debug-2023-11-16T10-12-21-0500
LOG_NAME="$1"
METHOD="$2"

logs=
if [ -z "$LOG_NAME" ]; then
  logs=( "$(ls "$LOG_PATH"/*.log)" )
else
  logs=("$LOG_PATH"/"$LOG_NAME")
fi



for log in ${logs[*]}; do
  log_name="$(basename "$log")"
  echo "---------------- START: $log_name --------------------"
  if ! [ -z "$METHOD" ]; then
      echo "retrieving rpc $METHOD counts/minute"
      # use awk to grep log for rpc method name and show calls/minute
      awk -v log_name="$log_name" -v method="method=$METHOD" '
      $0 ~ method {
          timestamp = $1 " " $2;
          gsub(/[\[\]]/, "", $3);               # Remove square brackets
          interval = substr(timestamp, 1, 16);  # Extract year, month, day, hour, and minute
          count[interval]++;
      }

      END {
          for (interval in count) {
              print "Log File: " log_name ", minute-interval: " interval ", " method ", Count: " count[interval];
          }
      }' "$log"
  else
      echo "retrieving all rpc method counts/minute"
      # Use awk to process the log file
      awk -v log_name="$log_name" '
        BEGIN {
          method_regex = "method=[^ ]+";
          printf "%-5s | %-30s | %-16s\n", "Count", "Method", "Minute-Interval";
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
            sub("method=", "", method);
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
      ' "$log" | tail -n +2 | sort -n -r
  fi
  echo ""
  echo "----- $log_name sorting all rpc method calls by counts -----"
  grep rpc_ < "$log" | grep --only-matching "method=[^ ]*" | sort | uniq -c | sort -r -n
  echo "-------------------------- END ------------------------------"
  echo ""
done