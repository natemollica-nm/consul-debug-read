#!/usr/bin/env bash

set -e

CONSUL_SERVER_MONITORING='https://raw.githubusercontent.com/sstarcher/grafana-dashboards/master/influxdb/consul.json'
DASHBOARDS_PATH=/opt/homebrew/var/lib/grafana/dashboards
API_KEY="$(CAT scripts/files/grafana-api-key)"



if ! test -d "${DASHBOARDS_PATH}"; then
  mkdir -p "${DASHBOARDS_PATH}"
fi

curl -sSFl https://raw.githubusercontent.com/sstarcher/grafana-dashboards/master/influxdb/consul.json > ${DASHBOARDS_PATH}/consul.json

status_code=$(curl -X POST -s \
   -H "Authorization: Bearer $API_KEY" \
   -H "Content-Type: application/json" \
   -d "@${DASHBOARDS_PATH}/consul.json" \
   http://consul:admin@localhost:3000/api/dashboards/db -w "%{http_code}")
