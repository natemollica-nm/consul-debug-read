# consul-debug-read
a simple cli tool for parsing consul-debug bundles to readable format


## Table of Contents

* [Getting Started](#Getting-Started)
* [Working with Debug Bundles](#Working-with-Debug-Bundles)
  * [Extract and Set Using CLI](#Extract-and-set-debug-path-using-CLI)
  * [Setting Debug Path](#settingchanging-debug-path)
  * [Using environment variable `CONSUL_DEBUG_PATH`](#Using-environment-variable)
* [Usage](#Usage)
  * [Consul Serf Membership](#Consul-Serf-Membership)
  * [Consul Raft Configuration](#Consul-Raft-Configuration)
  * [Consul Metrics Summary](#Consul-Metrics-Summary)
  * [Consul Metrics by Name](#Consul-Metrics-by-Name)
  * [Consul Host Metrics](#Consul-Host-Metrics)

## Getting Started

1. Clone this repository: 
  `$ git clone https://github.com/natemollica-nm/consul-debug-read.git`
2. Change to repo directory:
  `$ cd consul-debug-read`
3. Build and install binary: 
  `$ go install`
4. Test binary installed in path: 
  `$ consul-debug-read --help`

## Working with Debug Bundles

This tool uses the contents from the extracted bundle path to deliver a more useful and readable interpretation of the bundle.
The following sections explain how to point the tool to the right place using one of the three options:

* [Extract and set debug path using CLI](#Extract-and-set-debug-path-using-CLI) 
* [Set path to previously extracted bundle](#settingchanging-debug-path)
* [Using environment variable `CONSUL_DEBUG_PATH`](#Using-environment-variable)

### Extract and set debug path using CLI
1. Create and place the consul-debug.tar.gz file in a known location. For example

    ```shell
    # Create bundles directory and copy desired bundle tar.gz to dir
    $ mkdir -p ./bundles
    $ cp ~/Downloads/124722consul-debug-2023-10-04T18-29-47Z.tar.gz ./bundles/
    ```   
2. Run `consul-debug-read set-debug-path` using `--file`flag to both extract and set the debug directory to the extracted contents:

    ```shell
    $ consul-debug-read set-debug-path --file bundles/124722consul-debug-2023-10-04T18-29-47Z.tar.gz  
      2023/10/19 09:56:07 file passed in for extraction: bundles/124722consul-debug-2023-10-04T18-29-47Z.tar.gz
      2023/10/19 09:56:07 Extracting: bundles/124722consul-debug-2023-10-04T18-29-47Z.tar.gz
      2023/10/19 09:56:07 Destination File Extract Path: bundles/consul-debug-2023-10-04T18-29-47Z
      2023/10/19 09:56:08 Extraction of bundles/124722consul-debug-2023-10-04T18-29-47Z.tar.gz completed successfully.
      2023/10/19 09:56:08 set-debug-path: consul-debug-read debug-path has been set => bundles/consul-debug-2023-10-04T18-29-47Z
    ```
   
### Setting/changing debug path

1. Run `consul-debug-read set-debug-path` using `--path` flag to set the running config debug directory:
   ```shell
   $ consul-debug-read set-debug-path --path bundles/consul-debug-2023-10-04T18-29-47Z
     2023/10/19 10:09:15 set-debug-path: consul-debug-read debug-path has been set => bundles/consul-debug-2023-10-04T18-29-47Z
   ```

### Using environment variable

1. Export your terminal/shell session `CONSUL_DEBUG_PATH` variable:

    ```shell
    $ export CONSUL_DEBUG_PATH=bundles/consul-debug-2023-10-04T18-29-47Z
    $ consul-debug-read show-debug-path
      2023/10/19 14:46:31 using environment variable CONSUL_DEBUG_PATH - bundles/consul-debug-2023-10-04T18-29-47Z
      2023/10/19 14:46:31 debug-path => 'bundles/consul-debug-2023-10-04T18-29-47Z'
    ```

## Usage

1. Extract (if applicable) and set debug file path as outlined in [Working with Debug Bundles](#Working-with-debug-bundles) section above.
2. Explore bundle return options using `consul-debug-read --help`

### Consul Serf Membership

_consul debug only runs a cached members scrape to the `/v1/catalog/members?wan` endpoint_ 

Run: `consul-debug-read agent members`

```shell
# Example member return
Node                      Address             Status Type   Build      Protocol DC
hashi-i-01582cee96cd7dc0d 10.34.24.112:8302   Alive  server 1.15.3+ent 2        eu-01-stag
hashi-i-071c21a8d67edfe0d 10.34.23.73:8302    Alive  server 1.15.3+ent 2        eu-01-stag
hashi-i-073d7d2439f2e180f 10.34.45.248:8302   Alive  server 1.15.3+ent 2        eu-01-stag
hashi-i-0cc823f00596e8804 10.34.37.210:8302   Alive  server 1.15.3+ent 2        eu-01-stag
hashi-i-0dd679e8a2f054ac1 10.34.12.122:8302   Alive  server 1.15.3+ent 2        eu-01-stag
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

Run: `consul-debug-read raft-configuration`

```shell
# Example raft configuration return
Node              ID                                   Address             State    Voter
ip-10-135-25-56   c24d7789-af04-7bca-2649-42ebe6a227a3 10.135.25.56:8300   leader   true
ip-10-135-37-187  f20f69f5-3143-fdaa-3cd4-cde742808470 10.135.37.187:8300  follower true
ip-10-135-120-205 24128fc9-bc46-ef90-58b5-815ac343c12b 10.135.120.205:8300 follower true
ip-10-135-134-71  0060483d-9703-e017-087b-3f9635b462ab 10.135.134.71:8300  follower true
ip-10-135-78-52   3f8d935c-e08f-6d6f-2705-8c603fef1498 10.135.78.52:8300   follower true
```


### Consul Metrics Summary

Run: `consul-debug-read metrics`

```shell
# Example metrics summary overview return
Metrics Bundle Summary: bundles/consul-debug-2023-10-04T18-29-47Z/metrics.json
----------------------
Host Name: ip-10-135-37-187.ec2.internal
Agent Version: 1.15.6+ent
Interval: 30s
Duration: 5m2s
Capture Targets: [metrics logs pprof host agent members]
Raft State: Follower
```

### Consul Metrics by Name

Run: `consul-debug-read metrics --name consul.runtime.sys_bytes`


```shell
# Example return
Timestamp                     Metric                   Value   
2023-10-11 17:33:50 +0000 UTC consul.runtime.sys_bytes 4.62 GB 
2023-10-11 17:34:00 +0000 UTC consul.runtime.sys_bytes 4.62 GB 
2023-10-11 17:34:10 +0000 UTC consul.runtime.sys_bytes 4.62 GB 
2023-10-11 17:34:20 +0000 UTC consul.runtime.sys_bytes 4.62 GB 
2023-10-11 17:34:30 +0000 UTC consul.runtime.sys_bytes 4.62 GB 
2023-10-11 17:34:40 +0000 UTC consul.runtime.sys_bytes 4.62 GB 
2023-10-11 17:34:50 +0000 UTC consul.runtime.sys_bytes 4.62 GB 
2023-10-11 17:35:00 +0000 UTC consul.runtime.sys_bytes 4.62 GB 
2023-10-11 17:35:10 +0000 UTC consul.runtime.sys_bytes 4.62 GB 
2023-10-11 17:35:20 +0000 UTC consul.runtime.sys_bytes 4.62 GB 
2023-10-11 17:35:30 +0000 UTC consul.runtime.sys_bytes 4.62 GB 
2023-10-11 17:35:40 +0000 UTC consul.runtime.sys_bytes 4.62 GB 
```

### Consul Host Metrics

Run: `consul-debug-read metrics --host`

```shell
#Example Host Specific Metrics
Host Metrics Summary: bundles/consul-debug-2023-10-11T17-33-55Z/host.json
----------------------
OS: linux
Host Name hashi-i-073d7d2439f2e180f
Architecture: x86_64
Number of Cores: 8
CPU Vendor ID: GenuineIntel
CPU Model Name: Intel(R) Xeon(R) Platinum 8259CL CPU @ 2.50GHz
Platform: ubuntu | 20.04
Running Since: 2023-07-31 13:27:58 PDT
Uptime at Capture: 71 days, 21 hours, 5 minutes, 58 seconds

Host Memory Metrics Summary: bundles/consul-debug-2023-10-11T17-33-55Z/host.json
----------------------
Total: 30.89 GB
Used: 6.70 GB  (21.68%)
Total Available: 23.58 GB
VM Alloc Total: 32.00 TB
VM Alloc Used: 120.04 MB
Cached: 16.38 GB

Host Disk Metrics Summary: bundles/consul-debug-2023-10-11T17-33-55Z/host.json
----------------------
Used: 6.11 GB  (3.15%)
Free: 187.53 GB
Total: 193.65 GB
```

### Consul Agent Configuration

Run: `consul-debug-read agent --config`

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