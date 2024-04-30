package read

import (
	"encoding/json"
	"fmt"
	"github.com/ryanuber/columnize"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	BuildDate         time.Time `json:"BuildDate"`
	Datacenter        string    `json:"Datacenter"`
	NodeID            string    `json:"NodeID"`
	NodeName          string    `json:"NodeName"`
	PrimaryDatacenter string    `json:"PrimaryDatacenter"`
	Revision          string    `json:"Revision"`
	Server            bool      `json:"Server"`
	Version           string    `json:"Version"`
}

type Coord struct {
	Adjustment float64   `json:"Adjustment"`
	Error      float64   `json:"Error"`
	Height     float64   `json:"Height"`
	Vec        []float64 `json:"Vec"`
}

type DebugConfig struct {
	ACLEnableKeyListPolicy    bool   `json:"ACLEnableKeyListPolicy"`
	ACLInitialManagementToken string `json:"ACLInitialManagementToken"`
	ACLResolverSettings       struct {
		ACLDefaultPolicy string `json:"ACLDefaultPolicy"`
		ACLDownPolicy    string `json:"ACLDownPolicy"`
		ACLPolicyTTL     string `json:"ACLPolicyTTL"`
		ACLRoleTTL       string `json:"ACLRoleTTL"`
		ACLTokenTTL      string `json:"ACLTokenTTL"`
		ACLsEnabled      bool   `json:"ACLsEnabled"`
		Datacenter       string `json:"Datacenter"`
		EnterpriseMeta   struct {
			Namespace string `json:"Namespace"`
			Partition string `json:"Partition"`
		} `json:"EnterpriseMeta"`
		NodeName string `json:"NodeName"`
	} `json:"ACLResolverSettings"`
	ACLTokenReplication bool `json:"ACLTokenReplication"`
	ACLTokens           struct {
		ACLAgentRecoveryToken          string `json:"ACLAgentRecoveryToken"`
		ACLAgentToken                  string `json:"ACLAgentToken"`
		ACLConfigFileRegistrationToken string `json:"ACLConfigFileRegistrationToken"`
		ACLDefaultToken                string `json:"ACLDefaultToken"`
		ACLReplicationToken            string `json:"ACLReplicationToken"`
		DataDir                        string `json:"DataDir"`
		EnablePersistence              bool   `json:"EnablePersistence"`
		EnterpriseConfig               struct {
			ACLServiceProviderTokens []any `json:"ACLServiceProviderTokens"`
		} `json:"EnterpriseConfig"`
	} `json:"ACLTokens"`
	ACLsEnabled               bool   `json:"ACLsEnabled"`
	AEInterval                string `json:"AEInterval"`
	AdvertiseAddrLAN          string `json:"AdvertiseAddrLAN"`
	AdvertiseAddrWAN          string `json:"AdvertiseAddrWAN"`
	AdvertiseReconnectTimeout string `json:"AdvertiseReconnectTimeout"`
	AllowWriteHTTPFrom        []any  `json:"AllowWriteHTTPFrom"`
	AutoConfig                struct {
		Authorizer struct {
			AllowReuse bool `json:"AllowReuse"`
			AuthMethod struct {
				ACLAuthMethodEnterpriseFields struct {
					NamespaceRules []any `json:"NamespaceRules"`
				} `json:"ACLAuthMethodEnterpriseFields"`
				Config struct {
					BoundAudiences       any    `json:"BoundAudiences"`
					BoundIssuer          string `json:"BoundIssuer"`
					ClaimMappings        any    `json:"ClaimMappings"`
					ClockSkewLeeway      int    `json:"ClockSkewLeeway"`
					ExpirationLeeway     int    `json:"ExpirationLeeway"`
					JWKSCACert           string `json:"JWKSCACert"`
					Jwksurl              string `json:"JWKSURL"`
					JWTSupportedAlgs     any    `json:"JWTSupportedAlgs"`
					JWTValidationPubKeys any    `json:"JWTValidationPubKeys"`
					ListClaimMappings    any    `json:"ListClaimMappings"`
					NotBeforeLeeway      int    `json:"NotBeforeLeeway"`
					OIDCDiscoveryCACert  string `json:"OIDCDiscoveryCACert"`
					OIDCDiscoveryURL     string `json:"OIDCDiscoveryURL"`
				} `json:"Config"`
				Description    string `json:"Description"`
				DisplayName    string `json:"DisplayName"`
				EnterpriseMeta struct {
					Namespace string `json:"Namespace"`
					Partition string `json:"Partition"`
				} `json:"EnterpriseMeta"`
				MaxTokenTTL string `json:"MaxTokenTTL"`
				Name        string `json:"Name"`
				RaftIndex   struct {
					CreateIndex int `json:"CreateIndex"`
					ModifyIndex int `json:"ModifyIndex"`
				} `json:"RaftIndex"`
				TokenLocality string `json:"TokenLocality"`
				Type          string `json:"Type"`
			} `json:"AuthMethod"`
			ClaimAssertions []any `json:"ClaimAssertions"`
			Enabled         bool  `json:"Enabled"`
		} `json:"Authorizer"`
		DNSSANs         []any  `json:"DNSSANs"`
		Enabled         bool   `json:"Enabled"`
		IPSANs          []any  `json:"IPSANs"`
		IntroToken      string `json:"IntroToken"`
		IntroTokenFile  string `json:"IntroTokenFile"`
		ServerAddresses []any  `json:"ServerAddresses"`
	} `json:"AutoConfig"`
	AutoEncryptAllowTLS              bool     `json:"AutoEncryptAllowTLS"`
	AutoEncryptDNSSAN                []string `json:"AutoEncryptDNSSAN"`
	AutoEncryptIPSAN                 []string `json:"AutoEncryptIPSAN"`
	AutoEncryptTLS                   bool     `json:"AutoEncryptTLS"`
	AutoReloadConfig                 bool     `json:"AutoReloadConfig"`
	AutoReloadConfigCoalesceInterval string   `json:"AutoReloadConfigCoalesceInterval"`
	AutopilotCleanupDeadServers      bool     `json:"AutopilotCleanupDeadServers"`
	AutopilotDisableUpgradeMigration bool     `json:"AutopilotDisableUpgradeMigration"`
	AutopilotLastContactThreshold    string   `json:"AutopilotLastContactThreshold"`
	AutopilotMaxTrailingLogs         int      `json:"AutopilotMaxTrailingLogs"`
	AutopilotMinQuorum               int      `json:"AutopilotMinQuorum"`
	AutopilotRedundancyZoneTag       string   `json:"AutopilotRedundancyZoneTag"`
	AutopilotServerStabilizationTime string   `json:"AutopilotServerStabilizationTime"`
	AutopilotUpgradeVersionTag       string   `json:"AutopilotUpgradeVersionTag"`
	BindAddr                         string   `json:"BindAddr"`
	Bootstrap                        bool     `json:"Bootstrap"`
	BootstrapExpect                  int      `json:"BootstrapExpect"`
	BuildDate                        string   `json:"BuildDate"`
	Cache                            struct {
		EntryFetchMaxBurst int     `json:"EntryFetchMaxBurst"`
		EntryFetchRate     float64 `json:"EntryFetchRate"`
		Logger             any     `json:"Logger"`
	} `json:"Cache"`
	CheckDeregisterIntervalMin string   `json:"CheckDeregisterIntervalMin"`
	CheckOutputMaxSize         int      `json:"CheckOutputMaxSize"`
	CheckReapInterval          string   `json:"CheckReapInterval"`
	CheckUpdateInterval        string   `json:"CheckUpdateInterval"`
	Checks                     []any    `json:"Checks"`
	ClientAddrs                []string `json:"ClientAddrs"`
	Cloud                      struct {
		AuthURL         string `json:"AuthURL"`
		ClientID        string `json:"ClientID"`
		ClientSecret    string `json:"ClientSecret"`
		Hostname        string `json:"Hostname"`
		ManagementToken string `json:"ManagementToken"`
		NodeID          string `json:"NodeID"`
		NodeName        string `json:"NodeName"`
		ResourceID      string `json:"ResourceID"`
		ScadaAddress    string `json:"ScadaAddress"`
		TLSConfig       any    `json:"TLSConfig"`
	} `json:"Cloud"`
	ConfigEntryBootstrap []any `json:"ConfigEntryBootstrap"`
	ConnectCAConfig      struct {
	} `json:"ConnectCAConfig"`
	ConnectCAProvider                      string   `json:"ConnectCAProvider"`
	ConnectEnabled                         bool     `json:"ConnectEnabled"`
	ConnectMeshGatewayWANFederationEnabled bool     `json:"ConnectMeshGatewayWANFederationEnabled"`
	ConnectSidecarMaxPort                  int      `json:"ConnectSidecarMaxPort"`
	ConnectSidecarMinPort                  int      `json:"ConnectSidecarMinPort"`
	ConnectTestCALeafRootChangeSpread      string   `json:"ConnectTestCALeafRootChangeSpread"`
	ConsulCoordinateUpdateBatchSize        int      `json:"ConsulCoordinateUpdateBatchSize"`
	ConsulCoordinateUpdateMaxBatches       int      `json:"ConsulCoordinateUpdateMaxBatches"`
	ConsulCoordinateUpdatePeriod           string   `json:"ConsulCoordinateUpdatePeriod"`
	ConsulRaftElectionTimeout              string   `json:"ConsulRaftElectionTimeout"`
	ConsulRaftHeartbeatTimeout             string   `json:"ConsulRaftHeartbeatTimeout"`
	ConsulRaftLeaderLeaseTimeout           string   `json:"ConsulRaftLeaderLeaseTimeout"`
	ConsulServerHealthInterval             string   `json:"ConsulServerHealthInterval"`
	DNSARecordLimit                        int      `json:"DNSARecordLimit"`
	DNSAddrs                               []string `json:"DNSAddrs"`
	DNSAllowStale                          bool     `json:"DNSAllowStale"`
	DNSAltDomain                           string   `json:"DNSAltDomain"`
	DNSCacheMaxAge                         string   `json:"DNSCacheMaxAge"`
	DNSDisableCompression                  bool     `json:"DNSDisableCompression"`
	DNSDomain                              string   `json:"DNSDomain"`
	DNSEnableTruncate                      bool     `json:"DNSEnableTruncate"`
	DNSMaxStale                            string   `json:"DNSMaxStale"`
	DNSNodeMetaTXT                         bool     `json:"DNSNodeMetaTXT"`
	DNSNodeTTL                             string   `json:"DNSNodeTTL"`
	DNSOnlyPassing                         bool     `json:"DNSOnlyPassing"`
	DNSPort                                int      `json:"DNSPort"`
	DNSRecursorStrategy                    string   `json:"DNSRecursorStrategy"`
	DNSRecursorTimeout                     string   `json:"DNSRecursorTimeout"`
	DNSRecursors                           []string `json:"DNSRecursors"`
	DNSSOA                                 struct {
		Expire  int `json:"Expire"`
		Minttl  int `json:"Minttl"`
		Refresh int `json:"Refresh"`
		Retry   int `json:"Retry"`
	} `json:"DNSSOA"`
	DNSServiceTTL                    map[string]string `json:"DNSServiceTTL"`
	DNSUDPAnswerLimit                int               `json:"DNSUDPAnswerLimit"`
	DNSUseCache                      bool              `json:"DNSUseCache"`
	DataDir                          string            `json:"DataDir"`
	Datacenter                       string            `json:"Datacenter"`
	DefaultQueryTime                 string            `json:"DefaultQueryTime"`
	DevMode                          bool              `json:"DevMode"`
	DisableAnonymousSignature        bool              `json:"DisableAnonymousSignature"`
	DisableCoordinates               bool              `json:"DisableCoordinates"`
	DisableHTTPUnprintableCharFilter bool              `json:"DisableHTTPUnprintableCharFilter"`
	DisableHostNodeID                bool              `json:"DisableHostNodeID"`
	DisableKeyringFile               bool              `json:"DisableKeyringFile"`
	DisableRemoteExec                bool              `json:"DisableRemoteExec"`
	DisableUpdateCheck               bool              `json:"DisableUpdateCheck"`
	DiscardCheckOutput               bool              `json:"DiscardCheckOutput"`
	DiscoveryMaxStale                string            `json:"DiscoveryMaxStale"`
	EnableAgentTLSForChecks          bool              `json:"EnableAgentTLSForChecks"`
	EnableCentralServiceConfig       bool              `json:"EnableCentralServiceConfig"`
	EnableDebug                      bool              `json:"EnableDebug"`
	EnableLocalScriptChecks          bool              `json:"EnableLocalScriptChecks"`
	EnableRemoteScriptChecks         bool              `json:"EnableRemoteScriptChecks"`
	EncryptKey                       string            `json:"EncryptKey"`
	EnterpriseRuntimeConfig          struct {
		ACLMSPDisableBootstrap bool `json:"ACLMSPDisableBootstrap"`
		AuditEnabled           bool `json:"AuditEnabled"`
		AuditSinks             []struct {
			DeliveryGuarantee string `json:"DeliveryGuarantee"`
			FileName          string `json:"FileName"`
			Format            string `json:"Format"`
			Mode              int    `json:"Mode"`
			Name              string `json:"Name"`
			Path              string `json:"Path"`
			RotateBytes       int    `json:"RotateBytes"`
			RotateDuration    string `json:"RotateDuration"`
			RotateMaxFiles    int    `json:"RotateMaxFiles"`
			Type              string `json:"Type"`
		} `json:"AuditSinks"`
		DNSPreferNamespace    bool   `json:"DNSPreferNamespace"`
		LicensePath           string `json:"LicensePath"`
		LicensePollBaseTime   string `json:"LicensePollBaseTime"`
		LicensePollMaxTime    string `json:"LicensePollMaxTime"`
		LicenseUpdateBaseTime string `json:"LicenseUpdateBaseTime"`
		LicenseUpdateMaxTime  string `json:"LicenseUpdateMaxTime"`
		Partition             string `json:"Partition"`
	} `json:"EnterpriseRuntimeConfig"`
	ExposeMaxPort           int      `json:"ExposeMaxPort"`
	ExposeMinPort           int      `json:"ExposeMinPort"`
	GRPCAddrs               []string `json:"GRPCAddrs"`
	GRPCPort                int      `json:"GRPCPort"`
	GRPCTLSAddrs            []string `json:"GRPCTLSAddrs"`
	GRPCTLSPort             int      `json:"GRPCTLSPort"`
	GossipLANGossipInterval string   `json:"GossipLANGossipInterval"`
	GossipLANGossipNodes    int      `json:"GossipLANGossipNodes"`
	GossipLANProbeInterval  string   `json:"GossipLANProbeInterval"`
	GossipLANProbeTimeout   string   `json:"GossipLANProbeTimeout"`
	GossipLANRetransmitMult int      `json:"GossipLANRetransmitMult"`
	GossipLANSuspicionMult  int      `json:"GossipLANSuspicionMult"`
	GossipWANGossipInterval string   `json:"GossipWANGossipInterval"`
	GossipWANGossipNodes    int      `json:"GossipWANGossipNodes"`
	GossipWANProbeInterval  string   `json:"GossipWANProbeInterval"`
	GossipWANProbeTimeout   string   `json:"GossipWANProbeTimeout"`
	GossipWANRetransmitMult int      `json:"GossipWANRetransmitMult"`
	GossipWANSuspicionMult  int      `json:"GossipWANSuspicionMult"`
	HTTPAddrs               []string `json:"HTTPAddrs"`
	HTTPBlockEndpoints      []any    `json:"HTTPBlockEndpoints"`
	HTTPMaxConnsPerClient   int      `json:"HTTPMaxConnsPerClient"`
	HTTPMaxHeaderBytes      int      `json:"HTTPMaxHeaderBytes"`
	HTTPPort                int      `json:"HTTPPort"`
	HTTPResponseHeaders     struct {
	} `json:"HTTPResponseHeaders"`
	HTTPSAddrs                     []string `json:"HTTPSAddrs"`
	HTTPSHandshakeTimeout          string   `json:"HTTPSHandshakeTimeout"`
	HTTPSPort                      int      `json:"HTTPSPort"`
	HTTPUseCache                   bool     `json:"HTTPUseCache"`
	KVMaxValueSize                 int      `json:"KVMaxValueSize"`
	LeaveDrainTime                 string   `json:"LeaveDrainTime"`
	LeaveOnTerm                    bool     `json:"LeaveOnTerm"`
	LocalProxyConfigResyncInterval string   `json:"LocalProxyConfigResyncInterval"`
	Logging                        struct {
		EnableSyslog      bool   `json:"EnableSyslog"`
		LogFilePath       string `json:"LogFilePath"`
		LogJSON           bool   `json:"LogJSON"`
		LogLevel          string `json:"LogLevel"`
		LogRotateBytes    int    `json:"LogRotateBytes"`
		LogRotateDuration string `json:"LogRotateDuration"`
		LogRotateMaxFiles int    `json:"LogRotateMaxFiles"`
		Name              string `json:"Name"`
		SyslogFacility    string `json:"SyslogFacility"`
	} `json:"Logging"`
	MaxQueryTime string `json:"MaxQueryTime"`
	NodeID       string `json:"NodeID"`
	NodeMeta     struct {
	} `json:"NodeMeta"`
	NodeName                          string `json:"NodeName"`
	PeeringEnabled                    bool   `json:"PeeringEnabled"`
	PeeringTestAllowPeerRegistrations bool   `json:"PeeringTestAllowPeerRegistrations"`
	PidFile                           string `json:"PidFile"`
	PrimaryDatacenter                 string `json:"PrimaryDatacenter"`
	PrimaryGateways                   []any  `json:"PrimaryGateways"`
	PrimaryGatewaysInterval           string `json:"PrimaryGatewaysInterval"`
	RPCAdvertiseAddr                  string `json:"RPCAdvertiseAddr"`
	RPCBindAddr                       string `json:"RPCBindAddr"`
	RPCClientTimeout                  string `json:"RPCClientTimeout"`
	RPCConfig                         struct {
		EnableStreaming bool `json:"EnableStreaming"`
	} `json:"RPCConfig"`
	RPCHandshakeTimeout  string  `json:"RPCHandshakeTimeout"`
	RPCHoldTimeout       string  `json:"RPCHoldTimeout"`
	RPCMaxBurst          int     `json:"RPCMaxBurst"`
	RPCMaxConnsPerClient int     `json:"RPCMaxConnsPerClient"`
	RPCProtocol          int     `json:"RPCProtocol"`
	RPCRateLimit         float64 `json:"RPCRateLimit"`
	RaftLogStoreConfig   struct {
		Backend string `json:"Backend"`
		BoltDB  struct {
			NoFreelistSync bool `json:"NoFreelistSync"`
		} `json:"BoltDB"`
		DisableLogCache bool `json:"DisableLogCache"`
		Verification    struct {
			Enabled  bool   `json:"Enabled"`
			Interval string `json:"Interval"`
		} `json:"Verification"`
		Wal struct {
			SegmentSize int `json:"SegmentSize"`
		} `json:"WAL"`
	} `json:"RaftLogStoreConfig"`
	RaftProtocol          int    `json:"RaftProtocol"`
	RaftSnapshotInterval  string `json:"RaftSnapshotInterval"`
	RaftSnapshotThreshold int    `json:"RaftSnapshotThreshold"`
	RaftTrailingLogs      int    `json:"RaftTrailingLogs"`
	ReadReplica           bool   `json:"ReadReplica"`
	ReconnectTimeoutLAN   string `json:"ReconnectTimeoutLAN"`
	ReconnectTimeoutWAN   string `json:"ReconnectTimeoutWAN"`
	RejoinAfterLeave      bool   `json:"RejoinAfterLeave"`
	Reporting             struct {
		License struct {
			Enabled bool `json:"Enabled"`
		} `json:"License"`
	} `json:"Reporting"`
	RequestLimitsMode       int      `json:"RequestLimitsMode"`
	RequestLimitsReadRate   float64  `json:"RequestLimitsReadRate"`
	RequestLimitsWriteRate  float64  `json:"RequestLimitsWriteRate"`
	RetryJoinIntervalLAN    string   `json:"RetryJoinIntervalLAN"`
	RetryJoinIntervalWAN    string   `json:"RetryJoinIntervalWAN"`
	RetryJoinLAN            []string `json:"RetryJoinLAN"`
	RetryJoinMaxAttemptsLAN int      `json:"RetryJoinMaxAttemptsLAN"`
	RetryJoinMaxAttemptsWAN int      `json:"RetryJoinMaxAttemptsWAN"`
	RetryJoinWAN            []string `json:"RetryJoinWAN"`
	Revision                string   `json:"Revision"`
	SegmentLimit            int      `json:"SegmentLimit"`
	SegmentName             string   `json:"SegmentName"`
	SegmentNameLimit        int      `json:"SegmentNameLimit"`
	Segments                []any    `json:"Segments"`
	SerfAdvertiseAddrLAN    string   `json:"SerfAdvertiseAddrLAN"`
	SerfAdvertiseAddrWAN    string   `json:"SerfAdvertiseAddrWAN"`
	SerfAllowedCIDRsLAN     []any    `json:"SerfAllowedCIDRsLAN"`
	SerfAllowedCIDRsWAN     []any    `json:"SerfAllowedCIDRsWAN"`
	SerfBindAddrLAN         string   `json:"SerfBindAddrLAN"`
	SerfBindAddrWAN         string   `json:"SerfBindAddrWAN"`
	SerfPortLAN             int      `json:"SerfPortLAN"`
	SerfPortWAN             int      `json:"SerfPortWAN"`
	ServerMode              bool     `json:"ServerMode"`
	ServerName              string   `json:"ServerName"`
	ServerPort              int      `json:"ServerPort"`
	ServerRejoinAgeMax      string   `json:"ServerRejoinAgeMax"`
	Services                []any    `json:"Services"`
	SessionTTLMin           string   `json:"SessionTTLMin"`
	SkipLeaveOnInt          bool     `json:"SkipLeaveOnInt"`
	StaticRuntimeConfig     struct {
		EncryptVerifyIncoming bool `json:"EncryptVerifyIncoming"`
		EncryptVerifyOutgoing bool `json:"EncryptVerifyOutgoing"`
	} `json:"StaticRuntimeConfig"`
	SyncCoordinateIntervalMin string `json:"SyncCoordinateIntervalMin"`
	SyncCoordinateRateTarget  int    `json:"SyncCoordinateRateTarget"`
	TLS                       struct {
		AutoTLS                 bool   `json:"AutoTLS"`
		Domain                  string `json:"Domain"`
		EnableAgentTLSForChecks bool   `json:"EnableAgentTLSForChecks"`
		Grpc                    struct {
			CAFile               string `json:"CAFile"`
			CAPath               string `json:"CAPath"`
			CertFile             string `json:"CertFile"`
			CipherSuites         []any  `json:"CipherSuites"`
			KeyFile              string `json:"KeyFile"`
			TLSMinVersion        string `json:"TLSMinVersion"`
			UseAutoCert          bool   `json:"UseAutoCert"`
			VerifyIncoming       bool   `json:"VerifyIncoming"`
			VerifyOutgoing       bool   `json:"VerifyOutgoing"`
			VerifyServerHostname bool   `json:"VerifyServerHostname"`
		} `json:"GRPC"`
		HTTPS struct {
			CAFile               string `json:"CAFile"`
			CAPath               string `json:"CAPath"`
			CertFile             string `json:"CertFile"`
			CipherSuites         []any  `json:"CipherSuites"`
			KeyFile              string `json:"KeyFile"`
			TLSMinVersion        string `json:"TLSMinVersion"`
			UseAutoCert          bool   `json:"UseAutoCert"`
			VerifyIncoming       bool   `json:"VerifyIncoming"`
			VerifyOutgoing       bool   `json:"VerifyOutgoing"`
			VerifyServerHostname bool   `json:"VerifyServerHostname"`
		} `json:"HTTPS"`
		InternalRPC struct {
			CAFile               string `json:"CAFile"`
			CAPath               string `json:"CAPath"`
			CertFile             string `json:"CertFile"`
			CipherSuites         []any  `json:"CipherSuites"`
			KeyFile              string `json:"KeyFile"`
			TLSMinVersion        string `json:"TLSMinVersion"`
			UseAutoCert          bool   `json:"UseAutoCert"`
			VerifyIncoming       bool   `json:"VerifyIncoming"`
			VerifyOutgoing       bool   `json:"VerifyOutgoing"`
			VerifyServerHostname bool   `json:"VerifyServerHostname"`
		} `json:"InternalRPC"`
		NodeName   string `json:"NodeName"`
		ServerMode bool   `json:"ServerMode"`
		ServerName string `json:"ServerName"`
	} `json:"TLS"`
	TaggedAddresses struct {
		Lan     string `json:"lan"`
		LanIpv4 string `json:"lan_ipv4"`
		Wan     string `json:"wan"`
		WanIpv4 string `json:"wan_ipv4"`
	} `json:"TaggedAddresses"`
	Telemetry struct {
		AllowedPrefixes                    []string `json:"AllowedPrefixes"`
		BlockedPrefixes                    []string `json:"BlockedPrefixes"`
		CirconusAPIApp                     string   `json:"CirconusAPIApp"`
		CirconusAPIToken                   string   `json:"CirconusAPIToken"`
		CirconusAPIURL                     string   `json:"CirconusAPIURL"`
		CirconusBrokerID                   string   `json:"CirconusBrokerID"`
		CirconusBrokerSelectTag            string   `json:"CirconusBrokerSelectTag"`
		CirconusCheckDisplayName           string   `json:"CirconusCheckDisplayName"`
		CirconusCheckForceMetricActivation string   `json:"CirconusCheckForceMetricActivation"`
		CirconusCheckID                    string   `json:"CirconusCheckID"`
		CirconusCheckInstanceID            string   `json:"CirconusCheckInstanceID"`
		CirconusCheckSearchTag             string   `json:"CirconusCheckSearchTag"`
		CirconusCheckTags                  string   `json:"CirconusCheckTags"`
		CirconusSubmissionInterval         string   `json:"CirconusSubmissionInterval"`
		CirconusSubmissionURL              string   `json:"CirconusSubmissionURL"`
		Disable                            bool     `json:"Disable"`
		DisableHostname                    bool     `json:"DisableHostname"`
		DogstatsdAddr                      string   `json:"DogstatsdAddr"`
		DogstatsdTags                      []string `json:"DogstatsdTags"`
		EnableHostMetrics                  bool     `json:"EnableHostMetrics"`
		FilterDefault                      bool     `json:"FilterDefault"`
		MetricsPrefix                      string   `json:"MetricsPrefix"`
		PrometheusOpts                     struct {
			CounterDefinitions []any  `json:"CounterDefinitions"`
			Expiration         string `json:"Expiration"`
			GaugeDefinitions   []any  `json:"GaugeDefinitions"`
			Name               string `json:"Name"`
			Registerer         any    `json:"Registerer"`
			SummaryDefinitions []any  `json:"SummaryDefinitions"`
		} `json:"PrometheusOpts"`
		RetryFailedConfiguration bool   `json:"RetryFailedConfiguration"`
		StatsdAddr               string `json:"StatsdAddr"`
		StatsiteAddr             string `json:"StatsiteAddr"`
	} `json:"Telemetry"`
	TranslateWANAddrs bool `json:"TranslateWANAddrs"`
	TxnMaxReqLen      int  `json:"TxnMaxReqLen"`
	UIConfig          struct {
		ContentPath           string `json:"ContentPath"`
		DashboardURLTemplates struct {
		} `json:"DashboardURLTemplates"`
		Dir                        string `json:"Dir"`
		Enabled                    bool   `json:"Enabled"`
		HCPEnabled                 bool   `json:"HCPEnabled"`
		MetricsProvider            string `json:"MetricsProvider"`
		MetricsProviderFiles       []any  `json:"MetricsProviderFiles"`
		MetricsProviderOptionsJSON string `json:"MetricsProviderOptionsJSON"`
		MetricsProxy               struct {
			AddHeaders    []any  `json:"AddHeaders"`
			BaseURL       string `json:"BaseURL"`
			PathAllowlist []any  `json:"PathAllowlist"`
		} `json:"MetricsProxy"`
	} `json:"UIConfig"`
	UnixSocketGroup     string `json:"UnixSocketGroup"`
	UnixSocketMode      string `json:"UnixSocketMode"`
	UnixSocketUser      string `json:"UnixSocketUser"`
	UseStreamingBackend bool   `json:"UseStreamingBackend"`
	Version             string `json:"Version"`
	VersionMetadata     string `json:"VersionMetadata"`
	VersionPrerelease   string `json:"VersionPrerelease"`
	Watches             []any  `json:"Watches"`
	XDSUpdateRateLimit  int    `json:"XDSUpdateRateLimit"`
}

// AgentConfig is the standard agent.json structure for a user provided agent file.
// Reference: https://developer.hashicorp.com/consul/docs/agent/config/config-files#agents-configuration-file-reference
type AgentConfig struct {
	Domain            string   `json:"domain,omitempty"`
	Datacenter        string   `json:"datacenter,omitempty"`
	PrimaryDatacenter string   `json:"primary_datacenter,omitempty"`
	NodeName          string   `json:"node_name,omitempty"`
	DataDir           string   `json:"data_dir,omitempty"`
	Server            bool     `json:"server,omitempty"`
	Bootstrap         bool     `json:"bootstrap,omitempty"`
	BootstrapExpect   int      `json:"bootstrap_expect,omitempty"`
	RetryJoinLAN      []string `json:"retry_join,omitempty"`
	RetryJoinWAN      []string `json:"retry_join_wan,omitempty"`

	// Addresses
	Addresses struct {
		DNS     []string `json:"dns,omitempty"`
		HTTP    []string `json:"http,omitempty"`
		HTTPS   []string `json:"https,omitempty"`
		GRPC    []string `json:"grpc,omitempty"`
		GRPCTLS []string `json:"grpc_tls,omitempty"`
	} `json:"addresses,omitempty"`
	BindAddr                  string `json:"bind_addr,omitempty"`
	ClientAddr                string `json:"client_addr,omitempty"`
	AdvertiseAddrLAN          string `json:"advertise_addr,omitempty"`
	AdvertiseAddrLANIPv4      string `json:"advertise_addr_ipv4,omitempty"`
	AdvertiseAddrLANIPv6      string `json:"advertise_addr_ipv6,omitempty"`
	AdvertiseAddrWAN          string `json:"advertise_addr_wan,omitempty"`
	AdvertiseAddrWANIPv4      string `json:"advertise_addr_wan_ipv4,omitempty"`
	AdvertiseAddrWANIPv6      string `json:"advertise_addr_wan_ipv6,omitempty"`
	AdvertiseReconnectTimeout string `json:"advertise_reconnect_timeout,omitempty"`
	TranslateWANAddrs         bool   `json:"translate_wan_addrs,omitempty"`

	// Leave on Interrupt
	SkipLeaveOnInt bool `json:"skip_leave_on_interrupt" json:"skip_leave_on_interrupt,omitempty"`
	LeaveOnTerm    bool `json:"leave_on_terminate" json:"leave_on_terminate,omitempty"`

	DiscoveryMaxStale string `json:"discovery_max_stale" json:"discovery_max_stale,omitempty"`

	// Logging
	SyslogFacility string `json:"syslog_facility,omitempty"`
	LogLevel       string `json:"log_level,omitempty"`
	LogJSON        bool   `json:"log_json,omitempty"`
	LogFile        string `json:"log_file,omitempty"`

	// Ports
	Ports struct {
		DNS            int `json:"dns,omitempty"`
		HTTP           int `json:"http,omitempty"`
		HTTPS          int `json:"https,omitempty"`
		SerfLAN        int `json:"serf_lan,omitempty"`
		SerfWAN        int `json:"serf_wan,omitempty"`
		Server         int `json:"server,omitempty"`
		GRPC           int `json:"grpc,omitempty"`
		GRPCTLS        int `json:"grpc_tls,omitempty"`
		SidecarMinPort int `json:"sidecar_min_port,omitempty"`
		SidecarMaxPort int `json:"sidecar_max_port,omitempty"`
		ExposeMinPort  int `json:"expose_min_port,omitempty" `
		ExposeMaxPort  int `json:"expose_max_port,omitempty"`
	} `json:"ports,omitempty"`

	// ACL
	ACL struct {
		Enabled             bool   `json:"enabled,omitempty"`
		TokenReplication    bool   `json:"enable_token_replication,omitempty"`
		PolicyTTL           string `json:"policy_ttl,omitempty"`
		RoleTTL             string `json:"role_ttl,omitempty"`
		TokenTTL            string `json:"token_ttl,omitempty"`
		DownPolicy          string `json:"down_policy,omitempty"`
		DefaultPolicy       string `json:"default_policy,omitempty"`
		EnableKeyListPolicy bool   `json:"enable_key_list_policy,omitempty"`
		Tokens              struct {
			InitialManagement      string `json:"initial_management,omitempty"`
			Replication            string `json:"replication,omitempty"`
			AgentRecovery          string `json:"agent_recovery,omitempty"`
			Default                string `json:"default,omitempty"`
			Agent                  string `json:"agent,omitempty"`
			ConfigFileRegistration string `json:"config_file_service_registration,omitempty"`
			DNS                    string `json:"dns,omitempty"`
		} `json:"tokens,omitempty"`
		EnableTokenPersistence bool `json:"enable_token_persistence,omitempty"`
	} `json:"acl,omitempty"`

	// Gossip
	Encrypt               string `json:"encrypt,omitempty"`
	EncryptVerifyIncoming bool   `json:"encrypt_verify_incoming,omitempty"`
	EncryptVerifyOutgoing bool   `json:"encrypt_verify_outgoing,omitempty"`

	// Script/TLS Health Checks
	EnableAgentTLSForChecks    bool `json:"EnableAgentTLSForChecks,omitempty"`
	EnableCentralServiceConfig bool `json:"enable_central_service_config,omitempty"`

	// AutoEncrypt TLS
	AutoEncrypt struct {
		AllowTLS bool     `json:"allow_tls,omitempty"`
		TLS      bool     `json:"tls,omitempty"`
		DNSSAN   []string `json:"dns_san,omitempty"`
		IPSAN    []string `json:"ip_san,omitempty"`
	} `json:"auto_encrypt,omitempty"`

	// RPC TLS
	TLS struct {
		Defaults struct {
			CAFile               string `json:"ca_file,omitempty"`
			CAPath               string `json:"ca_path,omitempty"`
			CertFile             string `json:"cert_file,omitempty"`
			KeyFile              string `json:"key_file,omitempty"`
			TLSCipherSuites      []any  `json:"tls_cipher_suites,omitempty"`
			TLSMinVersion        string `json:"tls_min_version,omitempty"`
			VerifyIncoming       bool   `json:"verify_incoming,omitempty"`
			VerifyOutgoing       bool   `json:"verify_outgoing,omitempty"`
			VerifyServerHostname bool   `json:"verify_server_hostname,omitempty"`
		} `json:"defaults,omitempty"`
		Grpc struct {
			CAFile          string `json:"ca_file,omitempty"`
			CAPath          string `json:"ca_path,omitempty"`
			CertFile        string `json:"cert_file,omitempty"`
			KeyFile         string `json:"key_file,omitempty"`
			TLSCipherSuites []any  `json:"tls_cipher_suites,omitempty"`
			TLSMinVersion   string `json:"tls_min_version,omitempty"`
			VerifyIncoming  bool   `json:"verify_incoming,omitempty"`
			VerifyOutgoing  bool   `json:"verify_outgoing,omitempty"`
			UseAutoCert     bool   `json:"use_auto_cert,omitempty"`
		} `json:"grpc,omitempty"`
		HTTPS struct {
			CAFile          string `json:"ca_file,omitempty"`
			CAPath          string `json:"ca_path,omitempty"`
			CertFile        string `json:"cert_file,omitempty"`
			KeyFile         string `json:"key_file,omitempty"`
			TLSCipherSuites []any  `json:"tls_cipher_suites,omitempty"`
			TLSMinVersion   string `json:"tls_min_version,omitempty"`
			VerifyIncoming  bool   `json:"verify_incoming,omitempty"`
			VerifyOutgoing  bool   `json:"verify_outgoing,omitempty"`
		} `json:"https,omitempty"`
		InternalRPC struct {
			CAFile               string `json:"ca_file,omitempty"`
			CAPath               string `json:"ca_path,omitempty"`
			CertFile             string `json:"cert_file,omitempty"`
			KeyFile              string `json:"key_file,omitempty"`
			TLSCipherSuites      []any  `json:"tls_cipher_suites,omitempty"`
			TLSMinVersion        string `json:"tls_min_version,omitempty"`
			VerifyIncoming       bool   `json:"verify_incoming,omitempty"`
			VerifyOutgoing       bool   `json:"verify_outgoing,omitempty"`
			VerifyServerHostname bool   `json:"verify_server_hostname,omitempty"`
		} `json:"internal_rpc,omitempty"`
		ServerName string `json:"server_name,omitempty"`
		NodeName   string `json:"NodeName,omitempty"`
		ServerMode bool   `json:"ServerMode,omitempty"`
	} `json:"tls,omitempty"`

	// PPROF Debugging
	EnableDebug bool `json:"enable_debug" json:"enable_debug,omitempty"`

	// Telemetry
	Telemetry struct {
		CirconusAPIApp                     string   `json:"circonus_api_app,omitempty"`
		CirconusAPIToken                   string   `json:"circonus_api_token,omitempty"`
		CirconusAPIURL                     string   `json:"circonus_api_url,omitempty"`
		CirconusBrokerID                   string   `json:"circonus_broker_id,omitempty"`
		CirconusBrokerSelectTag            string   `json:"circonus_broker_select_tag,omitempty"`
		CirconusCheckDisplayName           string   `json:"circonus_check_display_name,omitempty"`
		CirconusCheckForceMetricActivation string   `json:"circonus_check_force_metric_activation,omitempty"`
		CirconusCheckID                    string   `json:"circonus_check_id,omitempty"`
		CirconusCheckInstanceID            string   `json:"circonus_check_instance_id,omitempty"`
		CirconusCheckSearchTag             string   `json:"circonus_check_search_tag,omitempty"`
		CirconusCheckTags                  string   `json:"circonus_check_tags,omitempty"`
		CirconusSubmissionInterval         string   `json:"circonus_submission_interval,omitempty"`
		CirconusSubmissionURL              string   `json:"circonus_submission_url,omitempty"`
		DisableHostname                    bool     `json:"disable_hostname,omitempty"`
		EnableHostMetrics                  bool     `json:"enable_host_metrics,omitempty"`
		DogstatsdAddr                      string   `json:"dogstatsd_addr,omitempty"`
		DogstatsdTags                      []string `json:"dogstatsd_tags,omitempty"`
		RetryFailedConfiguration           bool     `json:"retry_failed_connection,omitempty"`
		FilterDefault                      bool     `json:"filter_default,omitempty"`
		PrefixFilter                       []string `json:"prefix_filter,omitempty"`
		MetricsPrefix                      string   `json:"metrics_prefix,omitempty"`
		PrometheusRetentionTime            string   `json:"prometheus_retention_time,omitempty"`
		StatsdAddr                         string   `json:"statsd_address,omitempty"`
		StatsiteAddr                       string   `json:"statsite_address,omitempty"`
	} `json:"telemetry"`

	// DNS
	DNS struct {
		AllowStale         bool              `json:"allow_stale,omitempty"`
		ARecordLimit       int               `json:"a_record_limit,omitempty"`
		DisableCompression bool              `json:"disable_compression,omitempty"`
		EnableTruncate     bool              `json:"enable_truncate,omitempty"`
		MaxStale           string            `json:"max_stale,omitempty"`
		NodeTTL            string            `json:"node_ttl,omitempty"`
		OnlyPassing        bool              `json:"only_passing,omitempty"`
		RecursorStrategy   string            `json:"recursor_strategy,omitempty"`
		RecursorTimeout    string            `json:"recursor_timeout,omitempty"`
		ServiceTTL         map[string]string `json:"service_ttl,omitempty"`
		UDPAnswerLimit     int               `json:"udp_answer_limit,omitempty"`
		NodeMetaTXT        bool              `json:"enable_additional_node_meta_txt,omitempty"`
		SOA                struct {
			Refresh int `json:"refresh,omitempty"`
			Retry   int `json:"retry,omitempty"`
			Expire  int `json:"expire,omitempty"`
			Minttl  int `json:"min_ttl,omitempty"`
		} `json:"soa,omitempty"`
		UseCache    bool   `json:"use_cache"`
		CacheMaxAge string `json:"cache_max_age"`
	} `json:"dns,omitempty"`

	// Caching
	Cache struct {
		EntryFetchMaxBurst int     `json:"entry_fetch_max_burst"`
		EntryFetchRate     float64 `json:"entry_fetch_rate"`
	} `json:"cache,omitempty"`

	// RPC and HTTP Limits
	Limits struct {
		HTTPMaxConnsPerClient int    `json:"http_max_conns_per_client,omitempty"`
		HTTPSHandshakeTimeout string `json:"https_handshake_timeout,omitempty"`
		RequestLimits         struct {
			Mode      int     `json:"mode,omitempty"`
			ReadRate  float64 `json:"read_rate,omitempty"`
			WriteRate float64 `json:"write_rate,omitempty"`
		} `json:"request_limits,omitempty"`
		RPCClientTimeout     string  `json:"rpc_client_timeout,omitempty"`
		RPCHandshakeTimeout  string  `json:"rpc_handshake_timeout,omitempty"`
		RPCMaxBurst          int     `json:"rpc_max_burst,omitempty"`
		RPCMaxConnsPerClient int     `json:"rpc_max_conns_per_client,omitempty"`
		RPCRate              float64 `json:"rpc_rate,omitempty"`
		KVMaxValueSize       int     `json:"kv_max_value_size,omitempty"`
		TxnMaxReqLen         int     `json:"txn_max_req_len,omitempty"`
	} `json:"limits,omitempty"`

	// xDS Limits
	XDS struct {
		UpdateMaxPerSecond float64 `json:"update_max_per_second,omitempty"`
	} `json:"xds"`
}

type Member struct {
	Addr        string `json:"Addr"`
	DelegateCur int    `json:"DelegateCur"`
	DelegateMax int    `json:"DelegateMax"`
	DelegateMin int    `json:"DelegateMin"`
	Name        string `json:"Name"`
	Port        int    `json:"Port"`
	ProtocolCur int    `json:"ProtocolCur"`
	ProtocolMax int    `json:"ProtocolMax"`
	ProtocolMin int    `json:"ProtocolMin"`
	Status      int    `json:"Status"`
	Tags        struct {
		Acls        string `json:"acls"`
		Build       string `json:"build"`
		Dc          string `json:"dc"`
		Expect      string `json:"expect"`
		FtAdmpart   string `json:"ft_admpart"`
		FtFs        string `json:"ft_fs"`
		FtNs        string `json:"ft_ns"`
		FtSi        string `json:"ft_si"`
		GrpcTLSPort string `json:"grpc_tls_port"`
		ID          string `json:"id"`
		Port        string `json:"port"`
		RaftVsn     string `json:"raft_vsn"`
		Role        string `json:"role"`
		Segment     string `json:"segment"`
		UseTLS      string `json:"use_tls"`
		Vsn         string `json:"vsn"`
		VsnMax      string `json:"vsn_max"`
		VsnMin      string `json:"vsn_min"`
		WanJoinPort string `json:"wan_join_port"`
	} `json:"Tags"`
}

type Meta struct {
	ConsulNetworkSegment string `json:"consul-network-segment"`
}

type Stats struct {
	Agent struct {
		CheckMonitors string `json:"check_monitors"`
		CheckTtls     string `json:"check_ttls"`
		Checks        string `json:"checks"`
		Services      string `json:"services"`
	} `json:"agent"`
	Build struct {
		Prerelease      string `json:"prerelease"`
		Revision        string `json:"revision"`
		Version         string `json:"version"`
		VersionMetadata string `json:"version_metadata"`
	} `json:"build"`
	Consul struct {
		ACL              string `json:"acl"`
		Bootstrap        string `json:"bootstrap"`
		KnownDatacenters string `json:"known_datacenters"`
		Leader           string `json:"leader"`
		LeaderAddr       string `json:"leader_addr"`
		Server           string `json:"server"`
	} `json:"consul"`
	License struct {
		Customer       string `json:"customer"`
		ExpirationTime string `json:"expiration_time"`
		Features       string `json:"features"`
		ID             string `json:"id"`
		InstallID      string `json:"install_id"`
		IssueTime      string `json:"issue_time"`
		Modules        string `json:"modules"`
		Product        string `json:"product"`
		StartTime      string `json:"start_time"`
	} `json:"license"`
	Raft struct {
		AppliedIndex             string `json:"applied_index"`
		CommitIndex              string `json:"commit_index"`
		FsmPending               string `json:"fsm_pending"`
		LastContact              string `json:"last_contact"`
		LastLogIndex             string `json:"last_log_index"`
		LastLogTerm              string `json:"last_log_term"`
		LastSnapshotIndex        string `json:"last_snapshot_index"`
		LastSnapshotTerm         string `json:"last_snapshot_term"`
		LatestConfiguration      string `json:"latest_configuration"`
		LatestConfigurationIndex string `json:"latest_configuration_index"`
		NumPeers                 string `json:"num_peers"`
		ProtocolVersion          string `json:"protocol_version"`
		ProtocolVersionMax       string `json:"protocol_version_max"`
		ProtocolVersionMin       string `json:"protocol_version_min"`
		SnapshotVersionMax       string `json:"snapshot_version_max"`
		SnapshotVersionMin       string `json:"snapshot_version_min"`
		State                    string `json:"state"`
		Term                     string `json:"term"`
	} `json:"raft"`
	Runtime struct {
		Arch       string `json:"arch"`
		CPUCount   string `json:"cpu_count"`
		Goroutines string `json:"goroutines"`
		MaxProcs   string `json:"max_procs"`
		Os         string `json:"os"`
		Version    string `json:"version"`
	} `json:"runtime"`
	SerfLan struct {
		CoordinateResets string `json:"coordinate_resets"`
		Encrypted        string `json:"encrypted"`
		EventQueue       string `json:"event_queue"`
		EventTime        string `json:"event_time"`
		Failed           string `json:"failed"`
		HealthScore      string `json:"health_score"`
		IntentQueue      string `json:"intent_queue"`
		Left             string `json:"left"`
		MemberTime       string `json:"member_time"`
		Members          string `json:"members"`
		QueryQueue       string `json:"query_queue"`
		QueryTime        string `json:"query_time"`
	} `json:"serf_lan"`
	SerfWan struct {
		CoordinateResets string `json:"coordinate_resets"`
		Encrypted        string `json:"encrypted"`
		EventQueue       string `json:"event_queue"`
		EventTime        string `json:"event_time"`
		Failed           string `json:"failed"`
		HealthScore      string `json:"health_score"`
		IntentQueue      string `json:"intent_queue"`
		Left             string `json:"left"`
		MemberTime       string `json:"member_time"`
		Members          string `json:"members"`
		QueryQueue       string `json:"query_queue"`
		QueryTime        string `json:"query_time"`
	} `json:"serf_wan"`
}

type xDS struct {
	Port  int `json:"Port"`
	Ports struct {
		Plaintext int `json:"Plaintext"`
		TLS       int `json:"TLS"`
	} `json:"Ports"`
	SupportedProxies struct {
		Envoy []string `json:"envoy"`
	} `json:"SupportedProxies"`
}

type Agent struct {
	Config      Config      `json:"Config"`
	Coord       Coord       `json:"Coord"`
	DebugConfig DebugConfig `json:"DebugConfig"`
	Member      Member      `json:"Member"`
	Meta        Meta        `json:"Meta"`
	Stats       Stats       `json:"Stats"`
	XDS         xDS         `json:"xDS"`
	Members     []Member
}

// CompareVersion compares the Version field of Config with a given version string
// It returns 1 if the Config version is greater, -1 if the given version is greater,
// and 0 if they are equal.
func CompareVersion(c Config, givenVersion string) int {
	// Split the version strings into their components
	cParts := strings.Split(c.Version, ".")
	givenParts := strings.Split(givenVersion, ".")

	// Convert each component to integers and compare them
	for i := 0; i < len(cParts) && i < len(givenParts); i++ {
		cNum, _ := strconv.Atoi(cParts[i])
		givenNum, _ := strconv.Atoi(givenParts[i])

		if cNum > givenNum {
			return 1
		} else if cNum < givenNum {
			return -1
		}
	}

	// If all components are equal, check if one version has more components
	if len(cParts) > len(givenParts) {
		return 1
	} else if len(cParts) < len(givenParts) {
		return -1
	}

	return 0
}

// ByMemberName sorts members by name with a stable sort.
// 1. servers go at the top
// 2. group by datacenter name
// 3. sort by node name
type ByMemberName []Member

func (m ByMemberName) Len() int      { return len(m) }
func (m ByMemberName) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m ByMemberName) Less(i, j int) bool {
	tagsI := m[i].Tags
	tagsJ := m[j].Tags

	// put role=consul first
	switch {
	case tagsI.Role == "consul" && tagsJ.Role != "consul":
		return true
	case tagsI.Role != "consul" && tagsJ.Role == "consul":
		return false
	}

	// then by datacenter
	switch {
	case tagsI.Dc < tagsJ.Dc:
		return true
	case tagsI.Dc > tagsJ.Dc:
		return false
	}

	// finally by name
	return m[i].Name < m[j].Name
}

type RaftServer struct {
	// ID is the unique ID for the server. These are currently the same
	// as the address, but they will be changed to a real GUID in a future
	// release of Consul.
	ID string

	// Node is the node name of the server, as known by Consul, or this
	// will be set to "(unknown)" otherwise.
	Node string

	// Address is the IP:port of the server, used for Raft communications.
	Address string

	// Voter is true if this server has a vote in the cluster. This might
	// be false if the server is staging and still coming online, or if
	// it's a non-voting server, which will be added in a future release of
	// Consul.
	Voter bool
}

func (a *Agent) ParseDebugRaftConfig() string {
	raftConfig := ConvertToValidJSON(a.Stats.Raft.LatestConfiguration)
	return raftConfig
}

func (a *Agent) AgentConfigFull() (string, error) {
	var err error
	var userConfig *AgentConfig

	userConfig, err = a.toUserAgentConfig()
	if err != nil {
		return "", err
	}
	stringJson, err := json.MarshalIndent(userConfig, "", "    ")
	if err != nil {
		return "", err
	}
	return string(stringJson), nil
}

func (a *Agent) toUserAgentConfig() (*AgentConfig, error) {
	var agentConfig = &AgentConfig{
		Datacenter:        a.Config.Datacenter,
		PrimaryDatacenter: a.DebugConfig.PrimaryDatacenter,
		NodeName:          a.Config.NodeName,
		Server:            a.Config.Server,
		Bootstrap:         a.DebugConfig.Bootstrap,
		BootstrapExpect:   a.DebugConfig.BootstrapExpect,
		LogFile:           a.DebugConfig.Logging.LogFilePath,
		LogJSON:           a.DebugConfig.Logging.LogJSON,
		DataDir:           a.DebugConfig.DataDir,
		LogLevel:          a.DebugConfig.Logging.LogLevel,
		BindAddr:          a.DebugConfig.BindAddr,
		ClientAddr:        a.DebugConfig.ClientAddrs[0],
		DiscoveryMaxStale: a.DebugConfig.DiscoveryMaxStale,

		ACL: struct {
			Enabled             bool   `json:"enabled,omitempty"`
			TokenReplication    bool   `json:"enable_token_replication,omitempty"`
			PolicyTTL           string `json:"policy_ttl,omitempty"`
			RoleTTL             string `json:"role_ttl,omitempty"`
			TokenTTL            string `json:"token_ttl,omitempty"`
			DownPolicy          string `json:"down_policy,omitempty"`
			DefaultPolicy       string `json:"default_policy,omitempty"`
			EnableKeyListPolicy bool   `json:"enable_key_list_policy,omitempty"`
			Tokens              struct {
				InitialManagement      string `json:"initial_management,omitempty"`
				Replication            string `json:"replication,omitempty"`
				AgentRecovery          string `json:"agent_recovery,omitempty"`
				Default                string `json:"default,omitempty"`
				Agent                  string `json:"agent,omitempty"`
				ConfigFileRegistration string `json:"config_file_service_registration,omitempty"`
				DNS                    string `json:"dns,omitempty"`
			} `json:"tokens,omitempty"`
			EnableTokenPersistence bool `json:"enable_token_persistence,omitempty"`
		}(struct {
			Enabled             bool
			TokenReplication    bool
			PolicyTTL           string
			RoleTTL             string
			TokenTTL            string
			DownPolicy          string
			DefaultPolicy       string
			EnableKeyListPolicy bool
			Tokens              struct {
				InitialManagement      string
				Replication            string
				AgentRecovery          string
				Default                string
				Agent                  string
				ConfigFileRegistration string
				DNS                    string
			}
			EnableTokenPersistence bool
		}{Enabled: a.DebugConfig.ACLsEnabled, TokenReplication: a.DebugConfig.ACLTokenReplication,
			PolicyTTL:              a.DebugConfig.ACLResolverSettings.ACLPolicyTTL,
			RoleTTL:                a.DebugConfig.ACLResolverSettings.ACLRoleTTL,
			TokenTTL:               a.DebugConfig.ACLResolverSettings.ACLTokenTTL,
			DownPolicy:             a.DebugConfig.ACLResolverSettings.ACLDownPolicy,
			DefaultPolicy:          a.DebugConfig.ACLResolverSettings.ACLDefaultPolicy,
			EnableTokenPersistence: a.DebugConfig.ACLTokens.EnablePersistence,
			Tokens: struct {
				InitialManagement      string
				Replication            string
				AgentRecovery          string
				Default                string
				Agent                  string
				ConfigFileRegistration string
				DNS                    string
			}{InitialManagement: a.DebugConfig.ACLInitialManagementToken, Replication: a.DebugConfig.ACLTokens.ACLReplicationToken, AgentRecovery: a.DebugConfig.ACLTokens.ACLAgentRecoveryToken, Default: a.DebugConfig.ACLTokens.ACLDefaultToken, Agent: a.DebugConfig.ACLTokens.ACLAgentToken, ConfigFileRegistration: a.DebugConfig.ACLTokens.ACLConfigFileRegistrationToken}}),

		Addresses: struct {
			DNS     []string `json:"dns,omitempty"`
			HTTP    []string `json:"http,omitempty"`
			HTTPS   []string `json:"https,omitempty"`
			GRPC    []string `json:"grpc,omitempty"`
			GRPCTLS []string `json:"grpc_tls,omitempty"`
		}(struct {
			DNS     []string
			HTTP    []string
			HTTPS   []string
			GRPC    []string
			GRPCTLS []string
		}{DNS: a.DebugConfig.DNSAddrs, HTTP: a.DebugConfig.HTTPAddrs, HTTPS: a.DebugConfig.HTTPSAddrs, GRPC: a.DebugConfig.GRPCAddrs, GRPCTLS: a.DebugConfig.GRPCTLSAddrs}),
		Ports: struct {
			DNS            int `json:"dns,omitempty"`
			HTTP           int `json:"http,omitempty"`
			HTTPS          int `json:"https,omitempty"`
			SerfLAN        int `json:"serf_lan,omitempty"`
			SerfWAN        int `json:"serf_wan,omitempty"`
			Server         int `json:"server,omitempty"`
			GRPC           int `json:"grpc,omitempty"`
			GRPCTLS        int `json:"grpc_tls,omitempty"`
			SidecarMinPort int `json:"sidecar_min_port,omitempty"`
			SidecarMaxPort int `json:"sidecar_max_port,omitempty"`
			ExposeMinPort  int `json:"expose_min_port,omitempty" `
			ExposeMaxPort  int `json:"expose_max_port,omitempty"`
		}(struct {
			DNS            int
			HTTP           int
			HTTPS          int
			SerfLAN        int
			SerfWAN        int
			Server         int
			GRPC           int
			GRPCTLS        int
			SidecarMinPort int
			SidecarMaxPort int
			ExposeMinPort  int
			ExposeMaxPort  int
		}{DNS: a.DebugConfig.DNSPort, HTTP: a.DebugConfig.HTTPPort, HTTPS: a.DebugConfig.HTTPSPort, SerfLAN: a.DebugConfig.SerfPortLAN, SerfWAN: a.DebugConfig.SerfPortWAN, Server: a.DebugConfig.ServerPort, GRPC: a.DebugConfig.GRPCPort, GRPCTLS: a.DebugConfig.GRPCTLSPort, SidecarMinPort: a.DebugConfig.ConnectSidecarMinPort, SidecarMaxPort: a.DebugConfig.ConnectSidecarMaxPort, ExposeMinPort: a.DebugConfig.ExposeMinPort, ExposeMaxPort: a.DebugConfig.ExposeMaxPort}),

		AutoEncrypt: struct {
			AllowTLS bool     `json:"allow_tls,omitempty"`
			TLS      bool     `json:"tls,omitempty"`
			DNSSAN   []string `json:"dns_san,omitempty"`
			IPSAN    []string `json:"ip_san,omitempty"`
		}{AllowTLS: a.DebugConfig.AutoEncryptAllowTLS, TLS: a.DebugConfig.AutoEncryptTLS, DNSSAN: a.DebugConfig.AutoEncryptDNSSAN, IPSAN: a.DebugConfig.AutoEncryptIPSAN},

		TLS: struct {
			Defaults struct {
				CAFile               string `json:"ca_file,omitempty"`
				CAPath               string `json:"ca_path,omitempty"`
				CertFile             string `json:"cert_file,omitempty"`
				KeyFile              string `json:"key_file,omitempty"`
				TLSCipherSuites      []any  `json:"tls_cipher_suites,omitempty"`
				TLSMinVersion        string `json:"tls_min_version,omitempty"`
				VerifyIncoming       bool   `json:"verify_incoming,omitempty"`
				VerifyOutgoing       bool   `json:"verify_outgoing,omitempty"`
				VerifyServerHostname bool   `json:"verify_server_hostname,omitempty"`
			} `json:"defaults,omitempty"`
			Grpc struct {
				CAFile          string `json:"ca_file,omitempty"`
				CAPath          string `json:"ca_path,omitempty"`
				CertFile        string `json:"cert_file,omitempty"`
				KeyFile         string `json:"key_file,omitempty"`
				TLSCipherSuites []any  `json:"tls_cipher_suites,omitempty"`
				TLSMinVersion   string `json:"tls_min_version,omitempty"`
				VerifyIncoming  bool   `json:"verify_incoming,omitempty"`
				VerifyOutgoing  bool   `json:"verify_outgoing,omitempty"`
				UseAutoCert     bool   `json:"use_auto_cert,omitempty"`
			} `json:"grpc,omitempty"`
			HTTPS struct {
				CAFile          string `json:"ca_file,omitempty"`
				CAPath          string `json:"ca_path,omitempty"`
				CertFile        string `json:"cert_file,omitempty"`
				KeyFile         string `json:"key_file,omitempty"`
				TLSCipherSuites []any  `json:"tls_cipher_suites,omitempty"`
				TLSMinVersion   string `json:"tls_min_version,omitempty"`
				VerifyIncoming  bool   `json:"verify_incoming,omitempty"`
				VerifyOutgoing  bool   `json:"verify_outgoing,omitempty"`
			} `json:"https,omitempty"`
			InternalRPC struct {
				CAFile               string `json:"ca_file,omitempty"`
				CAPath               string `json:"ca_path,omitempty"`
				CertFile             string `json:"cert_file,omitempty"`
				KeyFile              string `json:"key_file,omitempty"`
				TLSCipherSuites      []any  `json:"tls_cipher_suites,omitempty"`
				TLSMinVersion        string `json:"tls_min_version,omitempty"`
				VerifyIncoming       bool   `json:"verify_incoming,omitempty"`
				VerifyOutgoing       bool   `json:"verify_outgoing,omitempty"`
				VerifyServerHostname bool   `json:"verify_server_hostname,omitempty"`
			} `json:"internal_rpc,omitempty"`
			ServerName string `json:"server_name,omitempty"`
			NodeName   string `json:"NodeName,omitempty"`
			ServerMode bool   `json:"ServerMode,omitempty"`
		}(struct {
			Defaults struct {
				CAFile               string
				CAPath               string
				CertFile             string
				KeyFile              string
				TLSCipherSuites      []any
				TLSMinVersion        string
				VerifyIncoming       bool
				VerifyOutgoing       bool
				VerifyServerHostname bool
			}
			Grpc struct {
				CAFile          string
				CAPath          string
				CertFile        string
				KeyFile         string
				TLSCipherSuites []any
				TLSMinVersion   string
				VerifyIncoming  bool
				VerifyOutgoing  bool
				UseAutoCert     bool
			}
			HTTPS struct {
				CAFile          string
				CAPath          string
				CertFile        string
				KeyFile         string
				TLSCipherSuites []any
				TLSMinVersion   string
				VerifyIncoming  bool
				VerifyOutgoing  bool
			}
			InternalRPC struct {
				CAFile               string
				CAPath               string
				CertFile             string
				KeyFile              string
				TLSCipherSuites      []any
				TLSMinVersion        string
				VerifyIncoming       bool
				VerifyOutgoing       bool
				VerifyServerHostname bool
			}
			ServerName string
			NodeName   string
			ServerMode bool
		}{
			ServerName: a.DebugConfig.TLS.ServerName,
			NodeName:   a.DebugConfig.TLS.NodeName,
			ServerMode: a.DebugConfig.TLS.ServerMode,
			Defaults: struct {
				CAFile               string
				CAPath               string
				CertFile             string
				KeyFile              string
				TLSCipherSuites      []any
				TLSMinVersion        string
				VerifyIncoming       bool
				VerifyOutgoing       bool
				VerifyServerHostname bool
			}{CAFile: a.DebugConfig.TLS.InternalRPC.CAFile, CAPath: a.DebugConfig.TLS.InternalRPC.CAPath, CertFile: a.DebugConfig.TLS.InternalRPC.CertFile, KeyFile: a.DebugConfig.TLS.InternalRPC.KeyFile, TLSCipherSuites: a.DebugConfig.TLS.InternalRPC.CipherSuites, TLSMinVersion: a.DebugConfig.TLS.InternalRPC.TLSMinVersion, VerifyIncoming: a.DebugConfig.TLS.InternalRPC.VerifyIncoming, VerifyOutgoing: a.DebugConfig.TLS.InternalRPC.VerifyOutgoing, VerifyServerHostname: a.DebugConfig.TLS.InternalRPC.VerifyServerHostname},

			Grpc: struct {
				CAFile          string
				CAPath          string
				CertFile        string
				KeyFile         string
				TLSCipherSuites []any
				TLSMinVersion   string
				VerifyIncoming  bool
				VerifyOutgoing  bool
				UseAutoCert     bool
			}{CAFile: a.DebugConfig.TLS.Grpc.CAFile, CAPath: a.DebugConfig.TLS.Grpc.CAPath, CertFile: a.DebugConfig.TLS.Grpc.CertFile, KeyFile: a.DebugConfig.TLS.Grpc.KeyFile, TLSCipherSuites: a.DebugConfig.TLS.Grpc.CipherSuites, TLSMinVersion: a.DebugConfig.TLS.Grpc.TLSMinVersion, VerifyIncoming: a.DebugConfig.TLS.Grpc.VerifyIncoming, VerifyOutgoing: a.DebugConfig.TLS.Grpc.VerifyOutgoing, UseAutoCert: a.DebugConfig.TLS.Grpc.UseAutoCert},

			HTTPS: struct {
				CAFile          string
				CAPath          string
				CertFile        string
				KeyFile         string
				TLSCipherSuites []any
				TLSMinVersion   string
				VerifyIncoming  bool
				VerifyOutgoing  bool
			}{CAFile: a.DebugConfig.TLS.HTTPS.CAFile, CAPath: a.DebugConfig.TLS.HTTPS.CAPath, CertFile: a.DebugConfig.TLS.HTTPS.CertFile, KeyFile: a.DebugConfig.TLS.HTTPS.KeyFile, TLSCipherSuites: a.DebugConfig.TLS.HTTPS.CipherSuites, TLSMinVersion: a.DebugConfig.TLS.HTTPS.TLSMinVersion, VerifyIncoming: a.DebugConfig.TLS.HTTPS.VerifyIncoming, VerifyOutgoing: a.DebugConfig.TLS.HTTPS.VerifyOutgoing},

			InternalRPC: struct {
				CAFile               string
				CAPath               string
				CertFile             string
				KeyFile              string
				TLSCipherSuites      []any
				TLSMinVersion        string
				VerifyIncoming       bool
				VerifyOutgoing       bool
				VerifyServerHostname bool
			}{CAFile: a.DebugConfig.TLS.InternalRPC.CAFile, CAPath: a.DebugConfig.TLS.InternalRPC.CAPath, CertFile: a.DebugConfig.TLS.InternalRPC.CertFile, KeyFile: a.DebugConfig.TLS.InternalRPC.KeyFile, TLSCipherSuites: a.DebugConfig.TLS.InternalRPC.CipherSuites, TLSMinVersion: a.DebugConfig.TLS.InternalRPC.TLSMinVersion, VerifyIncoming: a.DebugConfig.TLS.InternalRPC.VerifyIncoming, VerifyOutgoing: a.DebugConfig.TLS.InternalRPC.VerifyOutgoing, VerifyServerHostname: a.DebugConfig.TLS.InternalRPC.VerifyServerHostname},
		}),

		DNS: struct {
			AllowStale         bool              `json:"allow_stale,omitempty"`
			ARecordLimit       int               `json:"a_record_limit,omitempty"`
			DisableCompression bool              `json:"disable_compression,omitempty"`
			EnableTruncate     bool              `json:"enable_truncate,omitempty"`
			MaxStale           string            `json:"max_stale,omitempty"`
			NodeTTL            string            `json:"node_ttl,omitempty"`
			OnlyPassing        bool              `json:"only_passing,omitempty"`
			RecursorStrategy   string            `json:"recursor_strategy,omitempty"`
			RecursorTimeout    string            `json:"recursor_timeout,omitempty"`
			ServiceTTL         map[string]string `json:"service_ttl,omitempty"`
			UDPAnswerLimit     int               `json:"udp_answer_limit,omitempty"`
			NodeMetaTXT        bool              `json:"enable_additional_node_meta_txt,omitempty"`
			SOA                struct {
				Refresh int `json:"refresh,omitempty"`
				Retry   int `json:"retry,omitempty"`
				Expire  int `json:"expire,omitempty"`
				Minttl  int `json:"min_ttl,omitempty"`
			} `json:"soa,omitempty"`
			UseCache    bool   `json:"use_cache"`
			CacheMaxAge string `json:"cache_max_age"`
		}{
			AllowStale:         a.DebugConfig.DNSAllowStale,
			ARecordLimit:       a.DebugConfig.DNSARecordLimit,
			DisableCompression: a.DebugConfig.DNSDisableCompression,
			MaxStale:           a.DebugConfig.DNSMaxStale,
			NodeTTL:            a.DebugConfig.DNSNodeTTL,
			OnlyPassing:        a.DebugConfig.DNSOnlyPassing,
			RecursorStrategy:   a.DebugConfig.DNSRecursorStrategy,
			RecursorTimeout:    a.DebugConfig.DNSRecursorTimeout,
			ServiceTTL:         a.DebugConfig.DNSServiceTTL,
			UDPAnswerLimit:     a.DebugConfig.DNSUDPAnswerLimit,
			NodeMetaTXT:        a.DebugConfig.DNSNodeMetaTXT,
			SOA: struct {
				Refresh int `json:"refresh,omitempty"`
				Retry   int `json:"retry,omitempty"`
				Expire  int `json:"expire,omitempty"`
				Minttl  int `json:"min_ttl,omitempty"`
			}{
				Refresh: a.DebugConfig.DNSSOA.Refresh,
				Retry:   a.DebugConfig.DNSSOA.Retry,
				Expire:  a.DebugConfig.DNSSOA.Expire,
				Minttl:  a.DebugConfig.DNSSOA.Minttl,
			},
			UseCache:    a.DebugConfig.DNSUseCache,
			CacheMaxAge: a.DebugConfig.DNSCacheMaxAge,
		},

		Cache: struct {
			EntryFetchMaxBurst int     `json:"entry_fetch_max_burst"`
			EntryFetchRate     float64 `json:"entry_fetch_rate"`
		}{EntryFetchMaxBurst: a.DebugConfig.Cache.EntryFetchMaxBurst, EntryFetchRate: a.DebugConfig.Cache.EntryFetchRate},

		EnableDebug: a.DebugConfig.EnableDebug,

		Telemetry: struct {
			CirconusAPIApp                     string   `json:"circonus_api_app,omitempty"`
			CirconusAPIToken                   string   `json:"circonus_api_token,omitempty"`
			CirconusAPIURL                     string   `json:"circonus_api_url,omitempty"`
			CirconusBrokerID                   string   `json:"circonus_broker_id,omitempty"`
			CirconusBrokerSelectTag            string   `json:"circonus_broker_select_tag,omitempty"`
			CirconusCheckDisplayName           string   `json:"circonus_check_display_name,omitempty"`
			CirconusCheckForceMetricActivation string   `json:"circonus_check_force_metric_activation,omitempty"`
			CirconusCheckID                    string   `json:"circonus_check_id,omitempty"`
			CirconusCheckInstanceID            string   `json:"circonus_check_instance_id,omitempty"`
			CirconusCheckSearchTag             string   `json:"circonus_check_search_tag,omitempty"`
			CirconusCheckTags                  string   `json:"circonus_check_tags,omitempty"`
			CirconusSubmissionInterval         string   `json:"circonus_submission_interval,omitempty"`
			CirconusSubmissionURL              string   `json:"circonus_submission_url,omitempty"`
			DisableHostname                    bool     `json:"disable_hostname,omitempty"`
			EnableHostMetrics                  bool     `json:"enable_host_metrics,omitempty"`
			DogstatsdAddr                      string   `json:"dogstatsd_addr,omitempty"`
			DogstatsdTags                      []string `json:"dogstatsd_tags,omitempty"`
			RetryFailedConfiguration           bool     `json:"retry_failed_connection,omitempty"`
			FilterDefault                      bool     `json:"filter_default,omitempty"`
			PrefixFilter                       []string `json:"prefix_filter,omitempty"`
			MetricsPrefix                      string   `json:"metrics_prefix,omitempty"`
			PrometheusRetentionTime            string   `json:"prometheus_retention_time,omitempty"`
			StatsdAddr                         string   `json:"statsd_address,omitempty"`
			StatsiteAddr                       string   `json:"statsite_address,omitempty"`
		}(struct {
			CirconusAPIApp                     string
			CirconusAPIToken                   string
			CirconusAPIURL                     string
			CirconusBrokerID                   string
			CirconusBrokerSelectTag            string
			CirconusCheckDisplayName           string
			CirconusCheckForceMetricActivation string
			CirconusCheckID                    string
			CirconusCheckInstanceID            string
			CirconusCheckSearchTag             string
			CirconusCheckTags                  string
			CirconusSubmissionInterval         string
			CirconusSubmissionURL              string
			DisableHostname                    bool
			EnableHostMetrics                  bool
			DogstatsdAddr                      string
			DogstatsdTags                      []string
			RetryFailedConfiguration           bool
			FilterDefault                      bool
			PrefixFilter                       []string
			MetricsPrefix                      string
			PrometheusRetentionTime            string
			StatsdAddr                         string
			StatsiteAddr                       string
		}{CirconusAPIApp: a.DebugConfig.Telemetry.CirconusAPIApp, CirconusAPIToken: a.DebugConfig.Telemetry.CirconusAPIURL, CirconusAPIURL: a.DebugConfig.Telemetry.CirconusAPIURL, CirconusBrokerID: a.DebugConfig.Telemetry.CirconusBrokerID, CirconusBrokerSelectTag: a.DebugConfig.Telemetry.CirconusBrokerSelectTag, CirconusCheckDisplayName: a.DebugConfig.Telemetry.CirconusCheckDisplayName, CirconusCheckForceMetricActivation: a.DebugConfig.Telemetry.CirconusCheckForceMetricActivation, CirconusCheckID: a.DebugConfig.Telemetry.CirconusCheckInstanceID, CirconusCheckInstanceID: a.DebugConfig.Telemetry.CirconusCheckInstanceID, CirconusCheckSearchTag: a.DebugConfig.Telemetry.CirconusCheckSearchTag, CirconusCheckTags: a.DebugConfig.Telemetry.CirconusCheckTags, CirconusSubmissionInterval: a.DebugConfig.Telemetry.CirconusSubmissionInterval, CirconusSubmissionURL: a.DebugConfig.Telemetry.CirconusSubmissionURL, DisableHostname: a.DebugConfig.Telemetry.DisableHostname, EnableHostMetrics: a.DebugConfig.Telemetry.EnableHostMetrics, DogstatsdAddr: a.DebugConfig.Telemetry.DogstatsdAddr, DogstatsdTags: a.DebugConfig.Telemetry.DogstatsdTags, RetryFailedConfiguration: a.DebugConfig.Telemetry.RetryFailedConfiguration, FilterDefault: a.DebugConfig.Telemetry.FilterDefault, PrefixFilter: append(a.DebugConfig.Telemetry.BlockedPrefixes, a.DebugConfig.Telemetry.AllowedPrefixes...), MetricsPrefix: a.DebugConfig.Telemetry.MetricsPrefix, PrometheusRetentionTime: a.DebugConfig.Telemetry.PrometheusOpts.Expiration, StatsdAddr: a.DebugConfig.Telemetry.StatsdAddr, StatsiteAddr: a.DebugConfig.Telemetry.StatsiteAddr}),
		XDS: struct {
			UpdateMaxPerSecond float64 `json:"update_max_per_second,omitempty"`
		}(struct{ UpdateMaxPerSecond float64 }{UpdateMaxPerSecond: float64(a.DebugConfig.XDSUpdateRateLimit)}),
	}
	return agentConfig, nil
}

func (a *Agent) LogLevel() string {
	var defaultLogLevel string
	check := CompareVersion(a.Config, "1.13.0")
	switch check {
	case 1:
		defaultLogLevel = "TRACE"
	case -1:
		defaultLogLevel = "DEBUG"
	case 0:
		defaultLogLevel = "TRACE"
	}
	return defaultLogLevel
}

func (a *Agent) wanFederatedStatus() (string, bool) {
	if a.DebugConfig.ConnectMeshGatewayWANFederationEnabled {
		return "Mesh Gateway(s)", true
	} else if len(a.DebugConfig.RetryJoinIntervalWAN) > 0 {
		return "Basic (WAN Gossip)", true
	} else {
		return "N/A", false
	}
}

// WanMemberCount
// Function to count WAN Members
func (a *Agent) WanMemberCount() int {
	count := 0
	for _, member := range a.Members {
		if member.Tags.Dc != a.Config.Datacenter {
			count++
		}
	}
	return count
}

// FederatedDatacenterCount
// Function to count non-local datacenters in federated configuration
func (a *Agent) FederatedDatacenterCount() int {
	uniqueDatacenters := make(map[string]struct{})

	for _, member := range a.Members {
		if member.Tags.Dc != a.Config.Datacenter {
			uniqueDatacenters[member.Tags.Dc] = struct{}{}
		}
	}
	return len(uniqueDatacenters)
}

func (a *Agent) MembersStandard() string {
	if !a.Config.Server {
		return "=> bundle is from non-server consul agent (client agent). membership info unavailable (/v1/agent/members?wan)."
	}
	result := make([]string, 0, len(a.Members))
	header := "Node\x1fAddress\x1fStatus\x1fType\x1fBuild\x1fProtocol\x1fDC"
	result = append(result, header)
	sort.Sort(ByMemberName(a.Members))
	for _, member := range a.Members {
		tags := member.Tags

		addr := net.TCPAddr{IP: net.ParseIP(member.Addr), Port: int(member.Port)}
		protocol := tags.Vsn
		build := tags.Build
		if build == "" {
			build = "< 0.3"
		} else if idx := strings.Index(build, ":"); idx != -1 {
			build = build[:idx]
		}
		nameIdx := strings.Index(member.Name, ".")
		name := member.Name[:nameIdx]

		var statusString string
		switch {
		case member.Status == 0:
			statusString = "None"
		case member.Status == 1:
			statusString = "Alive"
		case member.Status == 2:
			statusString = "Leaving"
		case member.Status == 3:
			statusString = "Left"
		case member.Status == 4:
			statusString = "Failed"
		}
		switch tags.Role {
		case "node":
			line := fmt.Sprintf("%s\x1f%s\x1f%s\x1fclient\x1f%s\x1f%s\x1f%s",
				name, addr.String(), statusString, build, protocol, tags.Dc)
			result = append(result, line)

		case "consul":
			line := fmt.Sprintf("%s\x1f%s\x1f%s\x1fserver\x1f%s\x1f%s\x1f%s",
				name, addr.String(), statusString, build, protocol, tags.Dc)
			result = append(result, line)

		default:
			line := fmt.Sprintf("%s\x1f%s\x1f%s\x1funknown\x1f\x1f\x1f",
				name, addr.String(), statusString)
			result = append(result, line)
		}
	}

	output := columnize.Format(result, &columnize.Config{Delim: string([]byte{0x1f}), Glue: " "})
	return output
}

func (a *Agent) convertToRaftServer(raftDebugString string) ([]byte, error) {
	var correctedRaftConfig []byte
	// Define a struct to match the structure of your JSON data
	var data []map[string]interface{}

	// Unmarshal the JSON into the data structure
	err := json.Unmarshal([]byte(raftDebugString), &data)
	if err != nil {
		fmt.Println("Error:", err)
		return []byte(""), err
	}

	// Iterate through the data and replace "Suffrage": "Voter" with "Voter": true
	// Suffrage can be Voter, Nonvoter, or Staging (Deprecated)
	for i := range data {
		if suffrage, ok := data[i]["Suffrage"]; ok {
			if suffrage == "Voter" {
				data[i]["Voter"] = true
			} else {
				data[i]["Voter"] = false
			}
		}
		// Remove the old "Suffrage" key
		delete(data[i], "Suffrage")

		// Set raftServer "Node" field to corresponding member node name
		for _, member := range a.Members {
			if nodeID, ok := data[i]["ID"]; ok {
				if nodeID == member.Tags.ID {
					// Strip domain info from node name
					nameIdx := strings.Index(member.Name, ".")
					name := member.Name[:nameIdx]
					data[i]["Node"] = name
				}
			}
		}

	}

	// Marshal the data back to a JSON string
	correctedRaftConfig, err = json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
		return []byte(""), err
	}
	return correctedRaftConfig, nil
}

func (b *Debug) RaftListPeers() (string, error) {
	thisNode := b.Agent.Config.NodeName
	var debugBundleRaftConfig []byte
	var err error
	if !b.Agent.Config.Server {
		output := "=> bundle is from non-server consul agent (client agent). raft configuration unavailable."
		return output, nil
	}
	if debugBundleRaftConfig, err = b.Agent.convertToRaftServer(b.Agent.ParseDebugRaftConfig()); err != nil {
		return "", err
	}
	var raftServers []RaftServer
	err = json.Unmarshal(debugBundleRaftConfig, &raftServers)
	if err != nil {
		return "", err
	}

	// Format it as a nice table.
	result := []string{"Node\x1fID\x1fAddress\x1fState\x1fVoter\x1fAppliedIndex\x1fCommitIndex"}
	// Determine leader for processing output table
	raftLeaderAddr := b.Agent.Stats.Consul.LeaderAddr
	for _, s := range raftServers {
		state := "follower"
		if s.Address == raftLeaderAddr {
			state = "leader"
		}
		if s.Node == thisNode {
			appliedIndex := b.Agent.Stats.Raft.AppliedIndex
			commitIndex := b.Agent.Stats.Raft.CommitIndex
			result = append(result, fmt.Sprintf("%s\x1f%s\x1f%s\x1f%s\x1f%v\x1f%s\x1f%s",
				s.Node, s.ID, s.Address, state, s.Voter, appliedIndex, commitIndex))
		} else {
			result = append(result, fmt.Sprintf("%s\x1f%s\x1f%s\x1f%s\x1f%v\x1f%s\x1f%s",
				s.Node, s.ID, s.Address, state, s.Voter, "-", "-"))
		}
	}
	output := columnize.Format(result, &columnize.Config{Delim: string([]byte{0x1f}), Glue: " "})
	return output, nil
}

func (a *Agent) Summary() string {
	federationType, isFederated := a.wanFederatedStatus()
	var wanMemberCount, federatedDCCount int
	if isFederated {
		wanMemberCount = a.WanMemberCount()
		federatedDCCount = a.FederatedDatacenterCount()
	}
	title := "Agent Configuration Summary:"
	ul := strings.Repeat("-", len(title))
	return fmt.Sprintf("%s\n%s\nVersion: %s\nServer: %v\nRaft State: %s\nWAN Federation Status: %v\nWAN Federation Method: %s\nWAN Member Count: %d\nWAN Datacenter Count: %d\nDatacenter: %s\nPrimary DC: %s\nNodeName: %s\nSupported Envoy Versions: %v\n",
		title,
		ul,
		a.Config.Version,
		a.Config.Server,
		a.Stats.Raft.State,
		isFederated,
		federationType,
		wanMemberCount,
		federatedDCCount,
		a.Config.Datacenter,
		a.Config.PrimaryDatacenter,
		a.Config.NodeName,
		a.XDS.SupportedProxies.Envoy)
}
