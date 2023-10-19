# consul-debug-reader
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
Host Name hashi-i-073d7d2439f2e180f.node.seatgeek.stag
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
}
ACLsEnabled = true
AEInterval = 1m0s
AdvertiseAddrLAN = 10.135.37.187
AdvertiseAddrWAN = 10.135.37.187
AdvertiseReconnectTimeout = 0s
AllowWriteHTTPFrom = []
AutoConfig {
  Authorizer {
    AllowReuse = false
    AuthMethod {
      ACLAuthMethodEnterpriseFields {
        NamespaceRules = []
      }
      Config {
        BoundAudiences = <nil>
        BoundIssuer = 
        ClaimMappings = <nil>
        ClockSkewLeeway = 0
        ExpirationLeeway = 0
        JWKSCACert = 
        JWKSURL = 
        JWTSupportedAlgs = <nil>
        JWTValidationPubKeys = <nil>
        ListClaimMappings = <nil>
        NotBeforeLeeway = 0
        OIDCDiscoveryCACert = 
        OIDCDiscoveryURL = 
      }
      Description = 
      DisplayName = 
      EnterpriseMeta {
        Namespace = default
        Partition = default
      }
      MaxTokenTTL = 0s
      Name = Auto Config Authorizer
      RaftIndex {
        CreateIndex = 0
        ModifyIndex = 0
      }
      TokenLocality = 
      Type = jwt
    }
    ClaimAssertions = []
    Enabled = false
  }
  DNSSANs = []
  Enabled = false
  IPSANs = []
  IntroToken = hidden
  IntroTokenFile = 
  ServerAddresses = []
}
AutoEncryptAllowTLS = true
AutoEncryptDNSSAN = []
AutoEncryptIPSAN = []
AutoEncryptTLS = false
AutoReloadConfig = false
AutoReloadConfigCoalesceInterval = 1s
AutopilotCleanupDeadServers = true
AutopilotDisableUpgradeMigration = false
AutopilotLastContactThreshold = 200ms
AutopilotMaxTrailingLogs = 250
AutopilotMinQuorum = 0
AutopilotRedundancyZoneTag = 
AutopilotServerStabilizationTime = 10s
AutopilotUpgradeVersionTag = 
BindAddr = 0.0.0.0
Bootstrap = false
BootstrapExpect = 5
BuildDate = 2023-09-19 17:10:00 +0000 UTC
Cache {
  EntryFetchMaxBurst = 2
  EntryFetchRate = 1.7976931348623157e+308
  Logger = <nil>
}
CheckDeregisterIntervalMin = 1m0s
CheckOutputMaxSize = 4096
CheckReapInterval = 30s
CheckUpdateInterval = 15s
Checks = []
ClientAddrs = [0.0.0.0]
Cloud {
  AuthURL = 
  ClientID = 
  ClientSecret = hidden
  Hostname = 
  ManagementToken = hidden
  NodeID = 
  NodeName = ip-10-135-37-187
  ResourceID = 
  ScadaAddress = 
  TLSConfig = <nil>
}
ConfigEntryBootstrap = []
ConnectCAConfig {
}
ConnectCAProvider = 
ConnectEnabled = true
ConnectMeshGatewayWANFederationEnabled = false
ConnectSidecarMaxPort = 21255
ConnectSidecarMinPort = 21000
ConnectTestCALeafRootChangeSpread = 0s
ConsulCoordinateUpdateBatchSize = 128
ConsulCoordinateUpdateMaxBatches = 5
ConsulCoordinateUpdatePeriod = 5s
ConsulRaftElectionTimeout = 1s
ConsulRaftHeartbeatTimeout = 1s
ConsulRaftLeaderLeaseTimeout = 500ms
ConsulServerHealthInterval = 2s
DNSARecordLimit = 0
DNSAddrs = [tcp://0.0.0.0:8600 udp://0.0.0.0:8600]
DNSAllowStale = true
DNSAltDomain = 
DNSCacheMaxAge = 0s
DNSDisableCompression = false
DNSDomain = seatgeek.stag
DNSEnableTruncate = false
DNSMaxStale = 87600h0m0s
DNSNodeMetaTXT = true
DNSNodeTTL = 1m0s
DNSOnlyPassing = true
DNSPort = 8600
DNSRecursorStrategy = sequential
DNSRecursorTimeout = 2s
DNSRecursors = [169.254.169.253]
DNSSOA {
  Expire = 86400
  Minttl = 0
  Refresh = 3600
  Retry = 600
}
DNSServiceTTL {
  * = 10s
}
DNSUDPAnswerLimit = 3
DNSUseCache = false
DataDir = /data/consul
Datacenter = us-135-stag-default
DefaultQueryTime = 5m0s
DevMode = false
DisableAnonymousSignature = false
DisableCoordinates = false
DisableHTTPUnprintableCharFilter = false
DisableHostNodeID = true
DisableKeyringFile = false
DisableRemoteExec = true
DisableUpdateCheck = true
DiscardCheckOutput = false
DiscoveryMaxStale = 0s
EnableAgentTLSForChecks = false
EnableCentralServiceConfig = true
EnableDebug = true
EnableLocalScriptChecks = false
EnableRemoteScriptChecks = false
EncryptKey = hidden
EnterpriseRuntimeConfig {
  ACLMSPDisableBootstrap = false
  AuditEnabled = true
  AuditSinks = [{best-effort  json 292 log_file /data/logs/audit.json 1000000000 24h0m0s 3 file}]
  DNSPreferNamespace = false
  LicensePath = 
  LicensePollBaseTime = 6h0m0s
  LicensePollMaxTime = 24h0m0s
  LicenseUpdateBaseTime = 1m0s
  LicenseUpdateMaxTime = 20m0s
  Partition = 
}
ExposeMaxPort = 21755
ExposeMinPort = 21500
GRPCAddrs = []
GRPCPort = -1
GRPCTLSAddrs = [tcp://0.0.0.0:8503]
GRPCTLSPort = 8503
GossipLANGossipInterval = 200ms
GossipLANGossipNodes = 3
GossipLANProbeInterval = 1s
GossipLANProbeTimeout = 500ms
GossipLANRetransmitMult = 4
GossipLANSuspicionMult = 4
GossipWANGossipInterval = 500ms
GossipWANGossipNodes = 3
GossipWANProbeInterval = 5s
GossipWANProbeTimeout = 3s
GossipWANRetransmitMult = 4
GossipWANSuspicionMult = 6
HTTPAddrs = [tcp://0.0.0.0:8500]
HTTPBlockEndpoints = []
HTTPMaxConnsPerClient = 128
HTTPMaxHeaderBytes = 0
HTTPPort = 8500
HTTPResponseHeaders {
}
HTTPSAddrs = [tcp://0.0.0.0:8501]
HTTPSHandshakeTimeout = 5s
HTTPSPort = 8501
HTTPUseCache = true
KVMaxValueSize = 524288
LeaveDrainTime = 5s
LeaveOnTerm = false
LocalProxyConfigResyncInterval = 30s
Logging {
  EnableSyslog = false
  LogFilePath = /data/logs/
  LogJSON = true
  LogLevel = INFO
  LogRotateBytes = 1000000000
  LogRotateDuration = 24h0m0s
  LogRotateMaxFiles = 3
  Name = 
  SyslogFacility = LOCAL0
}
MaxQueryTime = 10m0s
NodeID = f20f69f5-3143-fdaa-3cd4-cde742808470
NodeMeta {
}
NodeName = ip-10-135-37-187
PeeringEnabled = false
PeeringTestAllowPeerRegistrations = false
PidFile = 
PrimaryDatacenter = us-east-stag
PrimaryGateways = []
PrimaryGatewaysInterval = 30s
RPCAdvertiseAddr = tcp://10.135.37.187:8300
RPCBindAddr = tcp://0.0.0.0:8300
RPCClientTimeout = 1m0s
RPCConfig {
  EnableStreaming = true
}
RPCHandshakeTimeout = 5s
RPCHoldTimeout = 7s
RPCMaxBurst = 1000
RPCMaxConnsPerClient = 100
RPCProtocol = 2
RPCRateLimit = 1.7976931348623157e+308
RaftLogStoreConfig {
  Backend = boltdb
  BoltDB {
    NoFreelistSync = false
  }
  DisableLogCache = false
  Verification {
    Enabled = false
    Interval = 0s
  }
  Wal {
    SegmentSize = 67108864
  }
}
RaftProtocol = 3
RaftSnapshotInterval = 30s
RaftSnapshotThreshold = 16384
RaftTrailingLogs = 10240
ReadReplica = false
ReconnectTimeoutLAN = 0s
ReconnectTimeoutWAN = 0s
RejoinAfterLeave = true
Reporting {
  License {
    Enabled = true
  }
}
RequestLimitsMode = 0
RequestLimitsReadRate = 1.7976931348623157e+308
RequestLimitsWriteRate = 1.7976931348623157e+308
RetryJoinIntervalLAN = 30s
RetryJoinIntervalWAN = 30s
RetryJoinLAN = [provider=aws service=ecs ecs_cluster=us-135-consul-server region=us-east-1 ecs_family=us-135-consul-server tag_key=hidden tag_value=consul-server]
RetryJoinMaxAttemptsLAN = 0
RetryJoinMaxAttemptsWAN = 0
RetryJoinWAN = [10.2.210.75 10.2.215.5 10.2.194.210 10.2.198.133 10.2.201.160 10.2.204.197]
Revision = 41f056ed
SegmentLimit = 64
SegmentName = 
SegmentNameLimit = 64
Segments = []
SerfAdvertiseAddrLAN = tcp://10.135.37.187:8301
SerfAdvertiseAddrWAN = tcp://10.135.37.187:8302
SerfAllowedCIDRsLAN = []
SerfAllowedCIDRsWAN = []
SerfBindAddrLAN = tcp://0.0.0.0:8301
SerfBindAddrWAN = tcp://0.0.0.0:8302
SerfPortLAN = 8301
SerfPortWAN = 8302
ServerMode = true
ServerName = 
ServerPort = 8300
ServerRejoinAgeMax = 168h0m0s
Services = []
SessionTTLMin = 0s
SkipLeaveOnInt = false
StaticRuntimeConfig {
  EncryptVerifyIncoming = true
  EncryptVerifyOutgoing = true
}
SyncCoordinateIntervalMin = 15s
SyncCoordinateRateTarget = 64
TLS {
  AutoTLS = false
  Domain = seatgeek.stag
  EnableAgentTLSForChecks = false
  Grpc {
    CAFile = /data/consul/consul-agent-ca.pem
    CAPath = 
    CertFile = /data/consul/us-135-stag-default-server-consul-0.pem
    CipherSuites = []
    KeyFile = hidden
    TLSMinVersion = TLSv1_2
    UseAutoCert = false
    VerifyIncoming = false
    VerifyOutgoing = false
    VerifyServerHostname = false
  }
  HTTPS {
    CAFile = /data/consul/consul-agent-ca.pem
    CAPath = 
    CertFile = /data/consul/us-135-stag-default-server-consul-0.pem
    CipherSuites = []
    KeyFile = hidden
    TLSMinVersion = TLSv1_2
    UseAutoCert = false
    VerifyIncoming = false
    VerifyOutgoing = false
    VerifyServerHostname = false
  }
  InternalRPC {
    CAFile = /data/consul/consul-agent-ca.pem
    CAPath = 
    CertFile = /data/consul/us-135-stag-default-server-consul-0.pem
    CipherSuites = []
    KeyFile = hidden
    TLSMinVersion = TLSv1_2
    UseAutoCert = false
    VerifyIncoming = false
    VerifyOutgoing = false
    VerifyServerHostname = false
  }
  NodeName = ip-10-135-37-187
  ServerMode = true
  ServerName = 
}
TaggedAddresses {
  lan = 10.135.37.187
  lan_ipv4 = 10.135.37.187
  wan = 10.135.37.187
  wan_ipv4 = 10.135.37.187
}
Telemetry {
  AllowedPrefixes = []
  BlockedPrefixes = [consul.rpc.server.call]
  CirconusAPIApp = 
  CirconusAPIToken = hidden
  CirconusAPIURL = 
  CirconusBrokerID = 
  CirconusBrokerSelectTag = 
  CirconusCheckDisplayName = 
  CirconusCheckForceMetricActivation = 
  CirconusCheckID = 
  CirconusCheckInstanceID = 
  CirconusCheckSearchTag = 
  CirconusCheckTags = 
  CirconusSubmissionInterval = 
  CirconusSubmissionURL = 
  Disable = false
  DisableHostname = false
  DogstatsdAddr = 127.0.0.1:8125
  DogstatsdTags = [vdc:us-135 env:stag environment:staging]
  EnableHostMetrics = false
  FilterDefault = true
  MetricsPrefix = consul
  PrometheusOpts {
    CounterDefinitions = []
    Expiration = 0s
    GaugeDefinitions = []
    Name = consul
    Registerer = <nil>
    SummaryDefinitions = []
  }
  RetryFailedConfiguration = true
  StatsdAddr = 
  StatsiteAddr = 
}
TranslateWANAddrs = false
TxnMaxReqLen = 524288
UIConfig {
  ContentPath = /ui/
  DashboardURLTemplates {
  }
  Dir = 
  Enabled = true
  HCPEnabled = false
  MetricsProvider = 
  MetricsProviderFiles = []
  MetricsProviderOptionsJSON = 
  MetricsProxy {
    AddHeaders = []
    BaseURL = 
    PathAllowlist = []
  }
}
UnixSocketGroup = 
UnixSocketMode = 
UnixSocketUser = 
UseStreamingBackend = true
Version = 1.15.6
VersionMetadata = ent
VersionPrerelease = 
Watches = []
XDSUpdateRateLimit = 250
```