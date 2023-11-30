# consul-debug-read
a simple cli tool for parsing consul-debug bundles to readable format

![](assets/consul-debug-read.gif)


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

### Download latest release

```shell
$ bash <(curl -sSL https://raw.githubusercontent.com/natemollica-nm/consul-debug-read/main/scripts/download.sh)
```

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
2. Run `consul-debug-read set-debug-path` using `--file` flag to both extract and set the debug directory to the extracted contents:

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

Run: `consul-debug-read raft-configuration`

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

Run: `consul-debug-read metrics`

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

### Consul Metrics by Name

Run: `consul-debug-read metrics --name consul.runtime.sys_bytes`


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

Run: `consul-debug-read metrics --host`

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
   `$ consul-debug-read --help`