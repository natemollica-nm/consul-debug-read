package types

type CPU struct {
	CacheSize  int      `json:"cacheSize"`
	CoreID     string   `json:"coreId"`
	Cores      int      `json:"cores"`
	CPU        int      `json:"cpu"`
	Family     string   `json:"family"`
	Flags      []string `json:"flags"`
	Mhz        float64  `json:"mhz"`
	Microcode  string   `json:"microcode"`
	Model      string   `json:"model"`
	ModelName  string   `json:"modelName"`
	PhysicalID string   `json:"physicalId"`
	Stepping   int      `json:"stepping"`
	VendorID   string   `json:"vendorId"`
}

type Disk struct {
	Free              float64 `json:"free"`
	Fstype            string  `json:"fstype"`
	InodesFree        int     `json:"inodesFree"`
	InodesTotal       int     `json:"inodesTotal"`
	InodesUsed        int     `json:"inodesUsed"`
	InodesUsedPercent float64 `json:"inodesUsedPercent"`
	Path              string  `json:"path"`
	Total             float64 `json:"total"`
	Used              float64 `json:"used"`
	UsedPercent       float64 `json:"usedPercent"`
}

type HostInfo struct {
	BootTime             int    `json:"bootTime"`
	HostID               string `json:"hostId"`
	Hostname             string `json:"hostname"`
	KernelArch           string `json:"kernelArch"`
	KernelVersion        string `json:"kernelVersion"`
	Os                   string `json:"os"`
	Platform             string `json:"platform"`
	PlatformFamily       string `json:"platformFamily"`
	PlatformVersion      string `json:"platformVersion"`
	Procs                int    `json:"procs"`
	Uptime               int    `json:"uptime"`
	VirtualizationRole   string `json:"virtualizationRole"`
	VirtualizationSystem string `json:"virtualizationSystem"`
}

type Memory struct {
	Active         int     `json:"active"`
	Available      float64 `json:"available"`
	Buffers        int     `json:"buffers"`
	Cached         float64 `json:"cached"`
	CommitLimit    float64 `json:"commitLimit"`
	CommittedAS    float64 `json:"committedAS"`
	Dirty          int     `json:"dirty"`
	Free           float64 `json:"free"`
	HighFree       int     `json:"highFree"`
	HighTotal      int     `json:"highTotal"`
	HugePageSize   int     `json:"hugePageSize"`
	HugePagesFree  int     `json:"hugePagesFree"`
	HugePagesRsvd  int     `json:"hugePagesRsvd"`
	HugePagesSurp  int     `json:"hugePagesSurp"`
	HugePagesTotal int     `json:"hugePagesTotal"`
	Inactive       float64 `json:"inactive"`
	Laundry        int     `json:"laundry"`
	LowFree        int     `json:"lowFree"`
	LowTotal       int     `json:"lowTotal"`
	Mapped         int     `json:"mapped"`
	PageTables     int     `json:"pageTables"`
	Shared         int     `json:"shared"`
	Slab           int     `json:"slab"`
	Sreclaimable   int     `json:"sreclaimable"`
	Sunreclaim     int     `json:"sunreclaim"`
	SwapCached     int     `json:"swapCached"`
	SwapFree       int     `json:"swapFree"`
	SwapTotal      int     `json:"swapTotal"`
	Total          float64 `json:"total"`
	Used           float64 `json:"used"`
	UsedPercent    float64 `json:"usedPercent"`
	VmallocChunk   int     `json:"vmallocChunk"`
	VmallocTotal   float64 `json:"vmallocTotal"`
	VmallocUsed    int     `json:"vmallocUsed"`
	Wired          int     `json:"wired"`
	WriteBack      int     `json:"writeBack"`
	WriteBackTmp   int     `json:"writeBackTmp"`
}

type Host struct {
	CPU            []CPU    `json:"CPU"`
	CollectionTime float64  `json:"CollectionTime"`
	Disk           Disk     `json:"Disk"`
	Errors         any      `json:"Errors"`
	Host           HostInfo `json:"Host"`
	Memory         Memory   `json:"Memory"`
}
