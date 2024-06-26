#!/usr/bin/env bash

export CONSUL_HTTP_TOKEN=root-token
eval "$(cat scripts/logging.sh)"
eval "$(cat scripts/formatting.env)"

SECONDARY_COUNT=1
PRIMARY_DATACENTER_NAME=dc1

exit_code=0
trap 'cleanup' TERM ERR
cleanup() {
      local i dc CONSUL_PIDS
      local DC_COUNT

      mapfile -t CONSUL_PIDS < <(pgrep consul)
      DC_COUNT="${#CONSUL_PIDS[@]}"
      if [ "${#CONSUL_PIDS[@]}" -gt 0 ]; then
          print_msg "local-hashi: Killing" "${CONSUL_PIDS[@]}"
          while [ "${#CONSUL_PIDS[@]}" -gt 0 ]; do
              pkill consul
              mapfile -t CONSUL_PIDS < <(pgrep consul)
              sleep 1
          done
      fi

      for i in $(seq 2 "$((DC_COUNT))"); do
          dc=dc"${i}"
          rm -rf /tmp/consul-server-"$dc"/*
          log "local-hashi: Removing lo0 interfaces 127.0.0.$i"
          sudo ifconfig lo0 127.0.0."${i}" down
          sudo ifconfig lo0 alias 127.0.0."${i}" delete
      done

      # Ensure base lo0 is up post script run
      log "local-hashi: Ensuring lo0 127.0.0.1 address is online"
      sudo ifconfig lo0 127.0.0.1 up;
      sudo ifconfig lo0 alias 127.0.0.1 up
      log "local-hashi: Cleanup done!"
      exit "$exit_code"
}

run_primary() {
    log "local-hashi: Starting Primary Datacenter server agent"
    sudo rm -rf /tmp/consul-server-dc1/*
    consul agent \
        -server \
        -bootstrap \
        -datacenter dc1 \
        -node consul-server-dc1 \
        -bind 127.0.0.1 \
        -client 127.0.0.1 \
        -advertise 127.0.0.1 \
        -advertise-wan 127.0.0.1 \
        -hcl 'ui = true' \
        -hcl 'enable_debug = true' \
        -hcl 'acl { enabled = true }' \
        -hcl 'acl { enable_token_persistence = true }' \
        -hcl "acl { tokens { initial_management = \"$CONSUL_HTTP_TOKEN\" } }" \
        -data-dir /tmp/consul-server-dc1 \
        1>/tmp/consul-server-dc1.log &
    log "local-hashi: Started Primary Datacenter server agent! (PID: $!)"; sleep 15
    consul members >/dev/null 2>&1 || {
        err "local-hashi: Failed running consul members!"
        return 1
    }
}

run_secondaries() {
      local i dc addr
      local count="$1"
      count=$((count))
      
    for i in $(seq 2 $((count))); do
        dc=dc"${i}"
        addr=127.0.0."${i}"
        sudo ifconfig lo0 alias "$addr" up
        log "local-hashi: Starting Secondary $dc server agent on $addr"
        consul agent \
            -server \
            -bootstrap \
            -datacenter "$dc" \
            -bind "$addr" \
            -client "$addr" \
            -advertise "$addr" \
            -advertise-wan "$addr" \
            -node consul-server-"$dc" \
            -retry-join-wan 127.0.0.1:8302 \
            -data-dir /tmp/consul-server-"$dc" \
            -log-level trace \
            -hcl 'ui = true' \
            -hcl 'enable_debug = true' \
            -hcl 'acl { enabled = true }' \
            -hcl 'primary_datacenter = "dc1"' \
            -hcl 'acl { enable_token_persistence = true }' \
            -hcl 'acl { enable_token_replication = true }' \
            -hcl "acl { tokens { initial_management = \"$CONSUL_HTTP_TOKEN\" } }" \
            -hcl 'acl { tokens { replication = "root-token" } }' \
            -hcl "ports {serf_lan = $((8301 + i*1000))}" \
            -hcl "ports {serf_wan = $((8302 + i*1000))}" \
            -hcl "ports {server   = $((8300 + i*1000))}" \
            -hcl "ports {http     = $((8500 + i*1000))}" \
            -hcl "ports {https    = $((8501 + i*1000))}" \
            -hcl "ports {grpc     = $((8502 + i*1000))}" \
            -hcl "ports {grpc_tls = $((8503 + i*1000))}" \
            -hcl "ports {dns      = $((8600 + i*1000))}" 1>/tmp/consul-server-"$dc".log &
     log "local-hashi: Started Secondary Datacenter consul-server-$dc agent! (PID: $!)"; sleep 5
     CONSUL_HTTP_ADDR=http://"$addr":$((8500 + i*1000)) consul members >/dev/null 2>&1 || {
          err "local-hashi: Failed running consul members (CONSUL_HTTP_ADDR: http://$addr:$((8500 + i*1000)))!"
          return 1
      }
    done
    log "local-hashi: Started $((count)) Consul datacenters!"
}

main() {
  local num_of_secondaries="$1"
  if [ "$CLEANUP" != true ]; then
    run_primary || {
        err "local-hashi: Failed to start primary DC!"
        exit
    }
    run_secondaries "$num_of_secondaries" || {
        err "local-hashi: Failed to start secondary datacenter!"
        exit
    }
  else
    log "local-hashi: Running cleanup!"
    cleanup
  fi
}

while test $# -gt 0 ; do
	case "$1" in
		-n|--secondaries)
			shift
      SECONDARY_COUNT=$1
			shift
			;;
		-d|--clean)
		  CLEANUP=true
		  shift
		  ;;
		*)
			echo "Usage: $(basename "$0") [-n num_of_dcs]"
			echo ""
			echo " -h, --help            Print help"
			echo " -n, --datacenters     Number of Consul DCs to start, default is $SECONDARY_COUNT"
			echo " -d, --clean           Teardown running Consul processes (leverages 'pgrep consul' output)"
			exit 0
			;;
	esac
done

main "$SECONDARY_COUNT"