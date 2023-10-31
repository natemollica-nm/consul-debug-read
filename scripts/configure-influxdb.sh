#!/bin/bash

set -e

# Setup initial default influxdb username, password, org, and generate operator token
echo "influxdb_configure: running influx setup"
influx setup \
   --org hashicorp \
   --bucket consul-debug-metrics \
   --username consul \
   --password hashicorp \
   --force

sleep 2
# Retrieve newly minted token for influxdb operations
influx_token=$(grep -E '^\s*token\s*=' "${HOME}/.influxdbv2/configs" | awk '{printf $3}')

echo "influxdb_token: $influx_token"
echo "influxdb_configure complete"
