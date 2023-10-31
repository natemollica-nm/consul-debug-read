#!/usr/bin/env bash

set -e

HOMEPATH=/opt/homebrew/opt/grafana/share/grafana
CONFIG_PATH=/opt/homebrew/etc/grafana/grafana.ini

# Add a signal handler to gracefully stop grafana when the script is terminated
cleanup() {
    echo "grafana: stopping grafana server"
    brew services stop grafana
}
trap cleanup SIGINT SIGTERM

function wait_timeout() {
  local timeout
  # Set a timeout (in seconds)
  timeout="$1"

  # Prompt message
  echo "Press any key within ${timeout} seconds or wait for timeout..."
  # Read user input with a timeout
  # shellcheck disable=SC2162
  # shellcheck disable=SC2034
  if read -t "$timeout" -n 1 -s key; then
      echo "continuing..."
  else
      echo "Timeout reached, continuing..."
  fi
}

function reset_admin_password() {
  local un pw
  un="$1"
  pw="$2"
  grafana cli --homepath "${HOMEPATH}" --config "${CONFIG_PATH}" "${un}" reset-admin-password "${pw}"
  sleep 8
}

brew services stop grafana
echo "grafana: creating grafana share, homepath, data, and conf directories"
for dir in \
  /opt/homebrew/opt/grafana/share \
  "${HOMEPATH}" \
  "${HOMEPATH}"/data \
  "${HOMEPATH}"/conf; do
  if [[ ! -d $dir ]]; then
    mkdir -p -m 0755 $dir
  fi
done
echo "updating defaults.ini and grafana.ini with consul/hashicorp un/pw"
cp -f scripts/files/grafana.ini "${HOMEPATH}"/conf/defaults.ini
cp -f scripts/files/grafana.ini "${CONFIG_PATH}"

# Start Grafana service
echo "starting grafana service (logs at - /opt/homebrew/var/log/grafana/grafana.log)"
brew services start grafana

echo "waiting for grafana to initialize"
sleep 10

# Obtain an API key from Grafana
GRAFANA_USERNAME="consul"
GRAFANA_PASSWORD="hashicorp"
GRAFANA_URL="http://${GRAFANA_USERNAME}:${GRAFANA_USERNAME}@localhost:3000"
BASIC_AUTH="${GRAFANA_USERNAME}:${GRAFANA_PASSWORD}"

# reset_admin_password "$GRAFANA_USERNAME" "$GRAFANA_PASSWORD}"

# Log in and create an organization and HTTP API key
echo "grafana: creating 'hashicorp' org"
curl -X POST \
  -s \
  -H "Content-Type: application/json" \
  -d '{"name":"hashicorp"}' \
  -u "${BASIC_AUTH}" "${GRAFANA_URL}"/api/orgs 1>/dev/null

# Get "hashicorp" org id
ORG_ID="$( curl \
  -s \
  -u "${BASIC_AUTH}" "${GRAFANA_URL}"/api/orgs | \
  jq -r '.[] | select(.name=="hashicorp").id' )" 1>/dev/null
echo "grafana: hashicorp org_id - $ORG_ID"


# Switch consul user to hashicorp organization
echo "grafana: setting consul admin user to hashicorp org"
curl \
  -X POST \
  -s \
  -u "${BASIC_AUTH}" \
  "${GRAFANA_URL}"/api/user/using/"${ORG_ID}" 1>/dev/null

OLD_API_KEY_ID="$( curl \
  -X GET \
  -s \
  -H "Content-Type: application/json" \
  -u "${BASIC_AUTH}" "${GRAFANA_URL}"/api/auth/keys | jq -r '.[] | select(.name=="consul_metrics").id' )"

if [[ -n "${OLD_API_KEY_ID}" ]]; then
  echo "grafana: previous API key found => $OLD_API_KEY_ID, deleting..."
  DELETED_CODE="$( curl \
    -X DELETE \
    -s \
    -H "Content-Type: application/json" \
    -u "${BASIC_AUTH}" \
    -o /dev/null \
    -w "%{http_code}" \
    "${GRAFANA_URL}"/api/auth/keys/"${OLD_API_KEY_ID}" )"
    [[ "${DELETED_CODE}" == "200" ]] || { echo "grafana: failed to remove previous API key, exit code: $DELETED_CODE"; }
fi

echo "grafana: generating hashicorp org API auth key"
API_KEY=$(curl \
  -X POST \
  -s \
  -H "Content-Type: application/json" \
  -d '{"name":"consul_metrics", "role": "Admin", "secondsToLive": 86400}' \
  -u "${BASIC_AUTH}" "${GRAFANA_URL}"/api/auth/keys | jq -r .key) 1>/dev/null

EXPIRATION=$(curl \
  -s \
  -X GET \
  -H "Content-Type: application/json" \
  -u "${BASIC_AUTH}" "${GRAFANA_URL}"/api/auth/keys | jq -r '.[] | .expiration' ) 1>/dev/null

if [[ -n $API_KEY ]] && [[ $API_KEY != 'null' ]]; then
  echo "hashicorp org_id: $ORG_ID | grafana API key: $API_KEY"
  echo "  => API key expiration date: $EXPIRATION"
  echo "$API_KEY" > scripts/files/grafana-api-key
  wait_timeout 5
else
  echo "failed to create API_KEY for grafana configuration"
  exit 1
fi

# Configure consul debug metrics datasource from influxDB
echo "grafana: updating grafana-config.json with influxdb token and hashicorp org"
DATASOURCE_CONFIG="scripts/files/grafana-config.json"
INFLUX_TOKEN=$(grep -E '^\s*token\s*=' "${HOME}/.influxdbv2/configs" | awk '{printf $3}' | tr -d '"')
jq \
  --arg ORG_ID "$ORG_ID" \
  --arg INFLUX_TOKEN "$INFLUX_TOKEN" \
  '.secureJsonData.token = $INFLUX_TOKEN | .orgId = $ORG_ID' \
  "$DATASOURCE_CONFIG" > tmpfile && mv tmpfile "$DATASOURCE_CONFIG"

echo "grafana: creating consul-debug-metrics-influx datasource"
status_code=$(curl -X POST -s \
   -H "Authorization: Bearer $API_KEY" \
   -H "Content-Type: application/json" \
   -d "@$DATASOURCE_CONFIG" \
   -o /dev/null -w "%{http_code}" \
   "${GRAFANA_URL}"/api/datasources)

if [ "$status_code" = '200' ]; then
  curl --request GET -s \
      -H "Authorization: Bearer $API_KEY" \
      -H "Content-Type: application/json" \
      "${GRAFANA_URL}"/api/datasources/name/consul-debug-metrics | jq .
  echo ""
  echo "grafana: datasource configuration successful!"
else
  echo "grafana: failed to configure consul-debug-metrics datasource, exit code: $status_code"
  exit 1
fi
exit 0


### Updating grafana datasource configuration
ds_uid="$(curl --request GET -s \
 -H "Authorization: Bearer $API_KEY" \
 -H "Content-Type: application/json" \
 "${GRAFANA_URL}"/api/datasources/name/consul-debug-metrics | jq -r .uid)"

status_code=$(curl -X PUT -s \
   -H "Authorization: Bearer $API_KEY" \
   -H "Content-Type: application/json" \
   -d "@$DATASOURCE_CONFIG" \
   -o /dev/null \
   -w "%{http_code}" \
   "${GRAFANA_URL}"/api/datasources/uid/"$ds_uid")


# Get auth keys
curl \
  -X GET \
  -s \
  -H "Content-Type: application/json" \
  -u "${BASIC_AUTH}" "${GRAFANA_URL}"/api/auth/keys