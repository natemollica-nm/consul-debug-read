# consul-debug-read
a simple cli tool for parsing consul-debug bundles to readable format

![](assets/consul-debug-read.gif)


## Table of Contents

* [Getting Started](#Getting-Started)
* [Working with Debug Bundles](#Working-with-Debug-Bundles)
  * [Extract and Set Using CLI](#Extract-and-set-path-using-CLI)
  * [Setting Debug Path](#settingchanging-debug-path)
  * [Using environment variable `CONSUL_DEBUG_PATH`](#Using-environment-variable)
* [Usage](#Usage)
  * [Consul Overall Summary](#consul-debug-overall-summary)
  * [Consul Log Parsing](#consul-log-parsing)
  * [Consul Serf Membership](#consul-serf-membership)
  * [Consul Raft Configuration](#consul-raft-configuration)
  * [Consul Metrics Summary](#consul-metrics-summary)
  * [Consul Metrics by Type](#consul-metrics-by-type)
  * [Consul Metrics by Name](#consul-metrics-by-name)
  * [Consul Host Metrics](#consul-host-metrics)

## Getting Started

### Download latest release

**linux|darwin architecture installs**

```shell
bash <(curl -sSL https://raw.githubusercontent.com/natemollica-nm/consul-debug-read/main/scripts/download.sh)
```

**_If you're having issues running the above scripted installation, you can enable more verbose output by running with '1' flag:_**

```shell
bash <(curl -sSL https://raw.githubusercontent.com/natemollica-nm/consul-debug-read/main/scripts/download.sh) 1
```


## Configuring consul-debug-read

This tool uses the contents from the extracted bundle path to deliver a more useful and readable interpretation of the bundle.
The following sections explain how to point the tool to the right place using one of the three options:

* Point to `.tar.gz` debug bundle to extract + set tool's location to the bundle's extracted root dir
* Point to a previously extracted bundle's root dir
* Set the `CONSUL_DEBUG_READ` environment variable and run the `config set-path` command

Run: `consul-debug-read config set-path [options]`

| Available Subcommands | Description                                                       |
|-----------------------|-------------------------------------------------------------------|
| `set-path`            | Configures CLI tool's extracted directory root of focus           |
| `current-path`        | Prints currently configured debug bundle extraction root of focus |

1. Identify the `.tar.gz` filepath or previously extracted root directory location you wish to examine.
2. Configure the tool to use this bundle's extract path:

```shell
# Example of set-path to directory where exact name of bundle is unknown | Selecting Option `1`
$ consul-debug-read config set-path -path $HOME/Downloads/

select a .tar.gz file to extract:
1: 124722consul-debug-2023-12-20T05-23-33Z.tar.gz  (27.00 MB)
2: 124722consul-debug-baseline-us-135-stag-default-2023-12-01T15-10-08-0500.tar.gz  (45.05 MB)
enter the number of the file to extract: 1
consul-debug-path set successfully
```

### Examining consul-debug-read config

Run: `consul-debug-read config [options]`

```shell
$ consul-debug-read config current-path                                                                                                                                                                                                                                                                 100%  
/Users/admin/hashi-bundles/consul-debug-2024-01-23T12-50-05-0500
```
    
### Using environment variable

1. Export your terminal/shell session `CONSUL_DEBUG_PATH` variable:

    ```shell
    $ export CONSUL_DEBUG_PATH=bundles/consul-debug-2023-10-04T18-29-47Z
    $ consul-debug-read set-path
      consul-debug-path set successfully using env var
    $ consul-debug-read current-path
      bundles/consul-debug-2023-10-04T18-29-47Z
    ```

## Usage

1. Extract (if applicable) and set debug directory path as outlined in [Configuring consul-debug-read](#configuring-consul-debug-read) section above.
2. Explore bundle return options using `consul-debug-read -help`


### Consul Debug Overall Summary

Run: `consul-debug-read summary`

```shell
# Example summary return
Consul Debug Bundle (2024-02-07 13:33:52): bundles/consul-debug-2024-02-07T12-40-42-0500
Debug Command Log Level: TRACE (Default) 
  * bundles/consul-debug-2024-02-07T12-40-42-0500/consul.log (1.57 MB)

Agent Configuration Summary:
----------------------------
Version: 1.17.1+ent
Server: true
Raft State: Leader
WAN Federation Status: true
WAN Federation Method: Basic (WAN Gossip)
WAN Member Count: 29
WAN Datacenter Count: 5
Datacenter: us-42-prod-default
Primary DC: us-east-primary
NodeName: ip-10-137-45-113
Supported Envoy Versions: [1.28.0 1.27.2 1.26.6 1.25.11 1.24.12 1.23.12 1.22.11]

Metrics Bundle Summary
----------------------
Datacenter: us-42-prod-default
Hostname: ip-10-28-78-42.ec2.internal
Agent Version: 1.17.2+ent
Raft State: Leader
Interval: 30s
Duration: 5m2s
Capture Targets: [metrics logs host agent members]
Total Captures: 30
Capture Time Start: 2024-02-07 17:40:40 +0000 UTC
Capture Time Stop: 2024-02-07 17:45:30 +0000 UTC

Host Summary:
-------------
OS: linux
Hostname: ip-10-28-78-42.ec2.internal
Architecture: aarch64
CPU Cores: 8
CPU Vendor ID: ARM
CPU Model: 
Platform: alpine | 3.17.7
Running Since: 1706635334
Total Uptime: 8 days, 0 hours, 18 minutes, 29 seconds

Host Memory Metrics Summary:
----------------------------
Used: 2.68 GB  (8.76%)
Total Available: 27.58 GB
Total: 30.63 GB

Host Disk Metrics Summary:
--------------------------
Used: 11.79 GB  (42.33%)
Total Available: 16.06 GB
Total: 29.36 GB
```

### Consul Log Parsing

Parse `INFO`, `WARN`, `ERROR`, `DEBUG`, and `TRACE` level log messages

Run: `consul-debug-read log <subcommand>`

| Available Subcommands | Description                                                                  |
|-----------------------|------------------------------------------------------------------------------|
| `parse-info`          | Returns all `[INFO]` messages                                                |
| `parse-warn`          | Returns all `[WARN]` messages                                                |
| `parse-error`         | Returns all `[ERROR]` messages                                               |
| `parse-debug`         | Returns all `[DEBUG]` messages                                               |
| `parse-trace`         | Returns all `[TRACE]` messages                                               |
| `parse-rpc-counts`    | Returns all `[TRACE]` messages pertaining to RPC rate limits in calls/minute |


| Available Options | Subcommand                                       | Description                                                                                                                        |
|-------------------|--------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------|
| `-message-count`  | `parse-[info\| warn \| error \| debug \| trace]` | Parse log for specific-level messages and return count sorted _(descending order)_ list of messages received                       |
| `-source-count`   | `parse-[info\| warn \| error \| debug \| trace]` | Parse log for specific-level messages and return count sorted _(descending order)_ list of messages received from specific sources |
| `-source`         | `parse-[info\| warn \| error \| debug \| trace]` | Capture specific-level messages from specific sources (e.g., "agent.http","agent.server", etc)                                     |
| `-method`         | `parse-rpc-counts`                               | Specify a specific RPC method for filtering RPC count results (e.g., "Catalog.NodeServiceList", "Health.ServiceNodes")             |

```shell
# Example using parse-debug sub-command
$ consul-debug-read log parse-debug                                                                                                                                                                                                                                                                     100%  
Timestamp            Source                                        Message                                                                                                                                                                                                  
2024-02-07T17:55:07Z agent.http                                    Request finished: method=GET url=/v1/catalog/datacenters from=10.137.149.96:56628 latency="427.384µs"                                                                                                    
2024-02-07T17:55:07Z agent.router                                  server in area left, skipping: server=consul-server-02.us-east area=wan func=GetDatacentersByDistance                                                                                           
2024-02-07T17:55:07Z agent.router                                  server in area left, skipping: server=consul-server-05.us-east area=wan func=GetDatacentersByDistance                                                                                           
2024-02-07T17:55:07Z agent.router                                  server in area left, skipping: server=consul-server-01.us-east area=wan func=GetDatacentersByDistance                                                                                           
2024-02-07T17:55:07Z agent.router                                  server in area left, skipping: server=consul-server-03.us-east area=wan func=GetDatacentersByDistance                                                                                           
2024-02-07T17:55:07Z agent.router                                  server in area left, skipping: server=consul-server-04.us-east area=wan func=GetDatacentersByDistance                                                                                           
2024-02-07T17:55:07Z agent.http                                    Request finished: method=GET url=/v1/status/leader from=10.137.127.240:51352 latency="103.746µs"                                                                                                         
2024-02-07T17:55:06Z agent.http                                    Request finished: method=GET url=/v1/status/leader from=10.137.89.52:33060 latency="83.135µs"
```

```shell
# Example using parse-debug -message-count
$ consul-debug-read log parse-debug -message-count                                                                                                                                                                                                                                                      100%  
Timestamp        Counts Source                                        Message                                                                                                                                                                                           
2024-02-07 17:54 48     agent.router                                  server in area left, skipping: server=consul-server-02.us-east area=wan func=GetDatacentersByDistance                                                                                    
2024-02-07 17:51 48     agent.router                                  server in area left, skipping: server=consul-server-01.us-east area=wan func=GetDatacentersByDistance                                                                                    
2024-02-07 17:54 48     agent.router                                  server in area left, skipping: server=consul-server-05.us-east area=wan func=GetDatacentersByDistance                                                                                    
2024-02-07 17:51 48     agent.router                                  server in area left, skipping: server=consul-server-03.us-east area=wan func=GetDatacentersByDistance                                                                                    
2024-02-07 17:51 48     agent.router                                  server in area left, skipping: server=consul-server-04.us-east area=wan func=GetDatacentersByDistance                                                                                    
2024-02-07 17:54 48     agent.router                                  server in area left, skipping: server=consul-server-01.us-east area=wan func=GetDatacentersByDistance
2024-02-07 17:52 1      agent.grpc-api.server-discovery.watch-servers starting stream: request_id=0813831d-145d-cfb7-5c95-cef228ff1cc3                                                                                                                                  
2024-02-07 17:54 1      agent.http                                    Request finished: method=GET url=/v1/catalog/datacenters from=10.137.29.125:55184 latency="432.267µs"                                                                                             
2024-02-07 17:51 1      agent.http                                    Request finished: method=GET url="/v1/health/service/queue-consumer?dc=us-east&index=2446522905&passing=1&stale=&wait=60000ms" from=10.137.123.50:34598 latency=1m0.000417706s                    
2024-02-07 17:52 1      agent.http                                    Request finished: method=GET url=/v1/status/leader from=10.137.56.240:1213 latency="80.772µs"                                                                                                     
2024-02-07 17:53 1      agent.http                                    Request finished: method=GET url="/v1/health/service/queue-consumer?dc=us-137-prod-default&passing=1&stale=&wait=60000ms" from=10.137.154.62:38878 latency="151.353µs"                            
2024-02-07 17:54 1      agent.http                                    Request finished: method=GET url=/v1/status/leader from=10.137.127.240:38958 latency="89.814µs"                                                                                                   
2024-02-07 17:54 1      agent.http                                    Request finished: method=GET url=/v1/status/leader from=10.137.127.240:55117 latency="78.409µs"
```

#### Consul RPC Rate Limiting Method Calls/Minute

Run: `consul-debug-read parse-rpc-counts`

```shell
# Example rpc-counts return
$ consul-debug-read log parse-rpc-counts                                                                                                                                                                                                                                                                100%  
Method                  Minute-Interval  Counts 
Health.ServiceNodes     2024-02-07 17:50 443    
Status.RaftStats        2024-02-07 17:51 300    
Status.RaftStats        2024-02-07 17:52 300    
Status.RaftStats        2024-02-07 17:53 300    
Status.RaftStats        2024-02-07 17:54 300    
Health.ServiceNodes     2024-02-07 17:54 293    
Status.RaftStats        2024-02-07 17:50 270    
Health.ServiceNodes     2024-02-07 17:53 158    
Health.ServiceNodes     2024-02-07 17:52 140    
Health.ServiceNodes     2024-02-07 17:51 116    
Status.Leader           2024-02-07 17:54 64     
Status.Leader           2024-02-07 17:53 64     
Status.Leader           2024-02-07 17:52 64     
Status.Leader           2024-02-07 17:51 64     
Status.Leader           2024-02-07 17:50 58     
ConnectCA.Roots         2024-02-07 17:50 51     
Catalog.ListDatacenters 2024-02-07 17:53 46     
Catalog.ListDatacenters 2024-02-07 17:51 45     
Catalog.ListDatacenters 2024-02-07 17:54 44     
Catalog.ListDatacenters 2024-02-07 17:52 44     
Status.RaftStats        2024-02-07 17:55 40     
Catalog.ListDatacenters 2024-02-07 17:50 34     
Health.ServiceNodes     2024-02-07 17:55 25     
ConnectCA.Roots         2024-02-07 17:51 22     
PreparedQuery.Execute   2024-02-07 17:50 12
```

### Consul Serf Membership

_consul debug only runs a cached members scrape to the `/v1/catalog/members?wan` endpoint_ 

Run: `consul-debug-read agent members`

```shell
# Example member return
Node                      Address             Status Type   Build      Protocol DC
consul-i-01582cee96cd7dc0d 10.34.24.112:8302   Alive  server 1.15.3+ent 2        eu-01-stag
consul-i-071c21a8d67edfe0d 10.34.23.73:8302    Alive  server 1.15.3+ent 2        eu-01-stag
consul-i-073d7d2439f2e180f 10.34.45.248:8302   Alive  server 1.15.3+ent 2        eu-01-stag
consul-i-0cc823f00596e8804 10.34.37.210:8302   Alive  server 1.15.3+ent 2        eu-01-stag
consul-i-0dd679e8a2f054ac1 10.34.12.122:8302   Alive  server 1.15.3+ent 2        eu-01-stag
ip-10-133-22-121          10.133.22.121:8302  Alive  server 1.15.6+ent 2        eu-133-stag-default
ip-10-133-45-202          10.133.45.202:8302  Alive  server 1.15.6+ent 2        eu-133-stag-default
ip-10-133-53-253          10.133.53.253:8302  Alive  server 1.15.6+ent 2        eu-133-stag-default
ip-10-133-92-59           10.133.92.59:8302   Alive  server 1.15.6+ent 2        eu-133-stag-default
ip-10-133-94-169          10.133.94.169:8302  Alive  server 1.15.6+ent 2        eu-133-stag-default
ip-10-135-120-205         10.135.120.205:8302 Alive  server 1.15.6+ent 2        us-135-stag-default
ip-10-135-134-71          10.135.134.71:8302  Alive  server 1.15.6+ent 2        us-135-stag-default
ip-10-135-25-56           10.135.25.56:8302   Alive  server 1.15.6+ent 2        us-135-stag-default
ip-10-135-37-187          10.135.37.187:8302  Alive  server 1.15.6+ent 2        us-135-stag-default
ip-10-135-78-52           10.135.78.52:8302   Alive  server 1.15.6+ent 2        us-135-stag-default
```

### Consul Raft Configuration

Run: `consul-debug-read agent raft-configuration`

```shell
# Example raft configuration return
Node                      ID                                   Address           State    Voter AppliedIndex CommitIndex
consul-i-0aa97949095868769 666e152f-7316-81aa-848b-3f4719564404 10.2.101.211:8300 leader   true  2696106337   2696106337
consul-i-0ba0dff4180ec2dc7 4f36f7ab-240a-61a6-c5e1-b78ce62813a2 10.2.4.230:8300   follower true  -            -
consul-i-08e67d882fe525809 1baa8d56-a9ae-adf7-5309-b12460c3e6c5 10.2.64.253:8300  follower true  -            -
consul-i-05a474f75fea384bb 263fd5e5-fbd7-90b1-a904-4ab3c53b74f7 10.2.17.109:8300  follower true  -            -
consul-i-06033dd57876bf1a7 eca79896-dad9-1713-94a2-c2b35a37d7df 10.2.4.89:8300    follower true  -            -
```


### Consul Metrics Summary

Run: `consul-debug-read metrics -summary`

```shell
# Example metrics summary overview return
Metrics Bundle Summary: bundles/consul-debug-2023-10-04T18-29-47Z/metrics.json
----------------------
Host Name: ip-10-135-37-187.ec2.internal
Agent Version: 1.15.6+ent
Raft State: Leader
Interval: 30s
Duration: 5m2s
Capture Targets: [metrics logs host agent members]
Total Captures: 30
Capture Time Start: 2023-10-23 15:03:40 +0000 UTC
Capture Time Stop: 2023-10-23 15:08:30 +0000 UTC
```

### Consul Metrics by Type

`consul-debug-read` comes with prebuilt variable groupings for each key-metric data type 
to aid in quickly assessing more important metrics to the scenario

| Command Option         | Description                                                                                             |
|------------------------|---------------------------------------------------------------------------------------------------------|
| `-auto-pilot`          | Retrieve key autopilot related metric values                                                           |
| `-bolt-db`             | Retrieve key boltDB related metric values for Consul from debug bundle                                  |
| `-dataplane-health`    | Retrieve key dataplane-related metric values for Consul from debug bundle                               |
| `-federation-health`   | Retrieve key secondary datacenter federation metric values for Consul from debug bundle                |
| `-host`                | Retrieve Host specific metrics                                                                          |
| `-key-metrics`         | Retrieve key metric values for Consul from debug bundle                                                 |
| `-leadership-health`   | Retrieve key raft leadership stability metric values for Consul from debug bundle                      |
| `-list-available-telemetry` | List available metric names as retrieved from consul telemetry docs                                |
| `-memory`              | Retrieve key memory metric values for Consul from debug bundle                                           |
| `-name=<string>`       | Retrieve specific metric timestamped values by name                                                      |
| `-network`             | Retrieve key network metric values for Consul from debug bundle                                           |
| `-raft-thread-health`  | Retrieve key raft thread saturation metric values for Consul from debug bundle                            |
| `-rate-limiting`       | Retrieve key rate limit metric values for Consul from debug bundle                                        |
| `-transaction-timing`  | Retrieve key transaction timing metric values for Consul from debug bundle                               |


```shell
# Example using -memory flag
$ consul-debug-read metrics -memory

                              consul.runtime.alloc_bytes 
                              -------------------------- 
Timestamp                     Metric                     Type  Unit  Value    
2024-02-07 17:44:40 +0000 UTC consul.runtime.alloc_bytes gauge bytes 55.17 MB 
2024-02-07 17:44:50 +0000 UTC consul.runtime.alloc_bytes gauge bytes 35.74 MB 
2024-02-07 17:45:00 +0000 UTC consul.runtime.alloc_bytes gauge bytes 45.64 MB 
2024-02-07 17:45:10 +0000 UTC consul.runtime.alloc_bytes gauge bytes 55.04 MB 
2024-02-07 17:45:20 +0000 UTC consul.runtime.alloc_bytes gauge bytes 65.78 MB 
2024-02-07 17:45:30 +0000 UTC consul.runtime.alloc_bytes gauge bytes 41.88 MB 

                              consul.runtime.heap_objects 
                              --------------------------- 
Timestamp                     Metric                      Type  Unit              Value  
2024-02-07 17:44:40 +0000 UTC consul.runtime.heap_objects gauge number of objects 365517 
2024-02-07 17:44:50 +0000 UTC consul.runtime.heap_objects gauge number of objects 162247 
2024-02-07 17:45:00 +0000 UTC consul.runtime.heap_objects gauge number of objects 252355 
2024-02-07 17:45:10 +0000 UTC consul.runtime.heap_objects gauge number of objects 341498 
2024-02-07 17:45:20 +0000 UTC consul.runtime.heap_objects gauge number of objects 434737 
2024-02-07 17:45:30 +0000 UTC consul.runtime.heap_objects gauge number of objects 217709 
                                                     
                                                  
                              consul.runtime.sys_bytes 
                              ------------------------ 
Timestamp                     Metric                   Type  Unit  Value     
2024-02-07 17:44:40 +0000 UTC consul.runtime.sys_bytes gauge bytes 191.30 MB 
2024-02-07 17:44:50 +0000 UTC consul.runtime.sys_bytes gauge bytes 191.30 MB 
2024-02-07 17:45:00 +0000 UTC consul.runtime.sys_bytes gauge bytes 191.30 MB 
2024-02-07 17:45:10 +0000 UTC consul.runtime.sys_bytes gauge bytes 191.30 MB 
2024-02-07 17:45:20 +0000 UTC consul.runtime.sys_bytes gauge bytes 191.30 MB 
2024-02-07 17:45:30 +0000 UTC consul.runtime.sys_bytes gauge bytes 191.30 MB 

                                                          
                              consul.runtime.total_gc_pause_ns 
                              -------------------------------- 
Timestamp                     Metric                           Type  Unit Value gc/min       
2024-02-07 17:44:40 +0000 UTC consul.runtime.total_gc_pause_ns gauge ns   1.77s -            
2024-02-07 17:44:50 +0000 UTC consul.runtime.total_gc_pause_ns gauge ns   1.77s 1.18ms/min   
2024-02-07 17:45:00 +0000 UTC consul.runtime.total_gc_pause_ns gauge ns   1.77s 0.0000ns/min 
2024-02-07 17:45:10 +0000 UTC consul.runtime.total_gc_pause_ns gauge ns   1.77s 0.0000ns/min 
2024-02-07 17:45:20 +0000 UTC consul.runtime.total_gc_pause_ns gauge ns   1.77s 0.0000ns/min 
2024-02-07 17:45:30 +0000 UTC consul.runtime.total_gc_pause_ns gauge ns   1.77s 1.01ms/min   
```

### Consul Metrics by Name

Run: `consul-debug-read metrics -name consul.runtime.sys_bytes`


```shell
# Example return
                              consul.runtime.sys_bytes           
                              ------------------------           
Timestamp                     Metric                             Type  Unit  Value    
2023-10-23 15:03:40 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:03:50 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:04:00 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:04:10 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:04:20 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:04:30 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:04:40 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:04:50 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:05:00 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:05:10 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:05:20 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:05:30 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:05:40 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:05:50 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:06:00 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:06:10 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:06:20 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:06:30 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:06:40 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:06:50 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:07:00 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:07:10 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:07:20 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:07:30 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:07:40 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:07:50 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:08:00 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:08:10 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:08:20 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB 
2023-10-23 15:08:30 +0000 UTC consul.runtime.sys_bytes.sys_bytes gauge bytes 16.25 GB
```

### Consul Host Metrics

Run: `consul-debug-read metrics -host`

```shell
#Example Host Specific Metrics
Host Metrics Summary: bundles/consul-debug-2023-10-23T11-03-40-0400/host.json
----------------------
OS: linux
Host Name consul-i-0aa97949095868769.node.consul
Architecture: x86_64
Number of Cores: 36
CPU Vendor ID: GenuineIntel
CPU Model Name: Intel(R) Xeon(R) Platinum 8124M CPU @ 3.00GHz
Platform: ubuntu | 20.04
Running Since: 2023-10-13 11:33:49 PDT
Uptime at Capture: 9 days, 20 hours, 29 minutes, 52 seconds

Host Memory Metrics Summary:
----------------------
Used: 15.92 GB  (23.21%)
Total Available: 51.75 GB
Total: 68.57 GB

Host Disk Metrics Summary:
----------------------
Used: 5.96 GB  (3.08%)
Free: 187.67 GB
Total: 193.65 GB
```

### Consul Agent Configuration

Run: `consul-debug-read agent -config`

```hcl
ACLEnableKeyListPolicy = false
ACLInitialManagementToken = hidden
ACLResolverSettings {
  ACLDefaultPolicy = allow
  ACLDownPolicy = extend-cache
  ACLPolicyTTL = 30s
  ACLRoleTTL = 0s
  ACLTokenTTL = 30s
  ACLsEnabled = true
  Datacenter = us-135-stag-default
  EnterpriseMeta {
    Namespace = 
    Partition = default
  }
  NodeName = ip-10-135-37-187
}
ACLTokenReplication = true
ACLTokens {
  ACLAgentRecoveryToken = hidden
  ACLAgentToken = hidden
  ACLConfigFileRegistrationToken = hidden
  ACLDefaultToken = hidden
  ACLReplicationToken = hidden
  DataDir = /data/consul
  EnablePersistence = true
  EnterpriseConfig {
    ACLServiceProviderTokens = []
  }
# .... cut ....
```

### Building and installing locally with Go

**Install golang**

Follow the [Download and Installation Instructions](https://go.dev/doc/install#tarball_non_standard) for installing go for your platform.

**Setup your **GOPATH** and **GOROOT** (if applicable) appropriately**

_If using macbook and installing go with homebrew set GOPATH and GOROOT as outlined in the examples below._

**GOPATH** is discussed in the [cmd/go documentation](http://golang.org/cmd/go/#hdr-GOPATH_environment_variable):

> The **GOPATH** environment variable lists places to look for Go code. On Unix, the value is a colon-separated string. On Windows, the value is a semicolon-separated string. On Plan 9, the value is a list.
**GOPATH** must be set to get, build and install packages outside the standard Go tree.

```shell
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

**GOROOT** is discussed in the [installation instructions](http://golang.org/doc/install#tarball_non_standard):

> The Go binary distributions assume they will be installed in /usr/local/go (or c:\Go under Windows), but it is possible to install the Go tools to a different location. In this case you must set the **GOROOT** environment variable to point to the directory in which it was installed.
For example, if you installed Go to your home directory you should add the following commands to $HOME/.profile:
```shell
# Example setting GOROOT with homebrew go installation.
export GOROOT=/opt/homebrew/opt/go/libexec
export PATH=$PATH:$GOROOT/bin
```
> **Note:** **GOROOT** must be set only when installing to a custom location.

1. Clone this repository:
   `$ git clone https://github.com/natemollica-nm/consul-debug-read.git`
2. Change to repo directory:
   `$ cd consul-debug-read`
3. Build and install binary:
   `$ go install`
4. Test binary installed in path:
   `$ consul-debug-read -help`