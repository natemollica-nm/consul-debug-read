package types

import "time"

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
	AutoEncryptAllowTLS              bool   `json:"AutoEncryptAllowTLS"`
	AutoEncryptDNSSAN                []any  `json:"AutoEncryptDNSSAN"`
	AutoEncryptIPSAN                 []any  `json:"AutoEncryptIPSAN"`
	AutoEncryptTLS                   bool   `json:"AutoEncryptTLS"`
	AutoReloadConfig                 bool   `json:"AutoReloadConfig"`
	AutoReloadConfigCoalesceInterval string `json:"AutoReloadConfigCoalesceInterval"`
	AutopilotCleanupDeadServers      bool   `json:"AutopilotCleanupDeadServers"`
	AutopilotDisableUpgradeMigration bool   `json:"AutopilotDisableUpgradeMigration"`
	AutopilotLastContactThreshold    string `json:"AutopilotLastContactThreshold"`
	AutopilotMaxTrailingLogs         int    `json:"AutopilotMaxTrailingLogs"`
	AutopilotMinQuorum               int    `json:"AutopilotMinQuorum"`
	AutopilotRedundancyZoneTag       string `json:"AutopilotRedundancyZoneTag"`
	AutopilotServerStabilizationTime string `json:"AutopilotServerStabilizationTime"`
	AutopilotUpgradeVersionTag       string `json:"AutopilotUpgradeVersionTag"`
	BindAddr                         string `json:"BindAddr"`
	Bootstrap                        bool   `json:"Bootstrap"`
	BootstrapExpect                  int    `json:"BootstrapExpect"`
	BuildDate                        string `json:"BuildDate"`
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
	DNSServiceTTL struct {
		Name string `json:"*"`
	} `json:"DNSServiceTTL"`
	DNSUDPAnswerLimit                int    `json:"DNSUDPAnswerLimit"`
	DNSUseCache                      bool   `json:"DNSUseCache"`
	DataDir                          string `json:"DataDir"`
	Datacenter                       string `json:"Datacenter"`
	DefaultQueryTime                 string `json:"DefaultQueryTime"`
	DevMode                          bool   `json:"DevMode"`
	DisableAnonymousSignature        bool   `json:"DisableAnonymousSignature"`
	DisableCoordinates               bool   `json:"DisableCoordinates"`
	DisableHTTPUnprintableCharFilter bool   `json:"DisableHTTPUnprintableCharFilter"`
	DisableHostNodeID                bool   `json:"DisableHostNodeID"`
	DisableKeyringFile               bool   `json:"DisableKeyringFile"`
	DisableRemoteExec                bool   `json:"DisableRemoteExec"`
	DisableUpdateCheck               bool   `json:"DisableUpdateCheck"`
	DiscardCheckOutput               bool   `json:"DiscardCheckOutput"`
	DiscoveryMaxStale                string `json:"DiscoveryMaxStale"`
	EnableAgentTLSForChecks          bool   `json:"EnableAgentTLSForChecks"`
	EnableCentralServiceConfig       bool   `json:"EnableCentralServiceConfig"`
	EnableDebug                      bool   `json:"EnableDebug"`
	EnableLocalScriptChecks          bool   `json:"EnableLocalScriptChecks"`
	EnableRemoteScriptChecks         bool   `json:"EnableRemoteScriptChecks"`
	EncryptKey                       string `json:"EncryptKey"`
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
	GRPCAddrs               []any    `json:"GRPCAddrs"`
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
		AllowedPrefixes                    []any    `json:"AllowedPrefixes"`
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
	xDS         xDS         `json:"xDS"`
}