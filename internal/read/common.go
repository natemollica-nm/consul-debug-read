package read

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	DebugReadEnvVar             = "CONSUL_DEBUG_PATH"
	DefaultCmdConfigFileName    = "config.yaml"
	DefaultCmdConfigFileDirName = ".consul-debug-read"
)

var (
	UserHomeDir, _          = os.UserHomeDir()
	CurrentDir, _           = os.Getwd()
	DebugReadConfigDirPath  = fmt.Sprintf("%s/%s", UserHomeDir, DefaultCmdConfigFileDirName)
	DebugReadConfigFullPath = fmt.Sprintf("%s/%s", DebugReadConfigDirPath, DefaultCmdConfigFileName)
	timeReg                 = regexp.MustCompile("^ns$|^ms$|^seconds$|^hours$")
	bytesReg                = regexp.MustCompile("bytes")
	percentageReg           = regexp.MustCompile("percentage")
)

func ToRFC3339(ts string) (string, error) {
	// Parse the timestamp string into a time.Time value
	timestamp, err := time.Parse("2006-01-02 15:04:05 -0700 MST", ts)
	if err != nil {
		return "", fmt.Errorf("error parsing timestamp: %v\n", err)
	}

	// Convert the time.Time value to RFC3339 format
	rfc3339Str := timestamp.Format(time.RFC3339)
	return rfc3339Str, nil
}

func extractTarGz(srcFile, destDir string) (string, error) {
	var extractRootDir string
	directoryPrefixCount := make(map[string]int)

	// Open the source .tar.gz file
	srcFileReader, err := os.Open(srcFile)
	if err != nil {
		return "", fmt.Errorf("extract-tar-gz: failed to open %s: %v\n", srcFile, err)
	}
	defer srcFileReader.Close()
	cleanup := func(err error) error {
		_ = srcFileReader.Close()
		return err
	}
	// Create a gzip reader
	gzipReader, err := gzip.NewReader(srcFileReader)
	if err != nil {
		return "", err
	}
	defer gzipReader.Close()

	// Create a tar reader
	tarReader := tar.NewReader(gzipReader)

	// Iterate through the tar archive and extract files
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return "", err
		}

		// Calculate the file path for extraction
		destFilePath := fmt.Sprintf("%s/%s", destDir, header.Name)

		// Create directories as needed
		if header.FileInfo().IsDir() {
			if err := os.MkdirAll(destFilePath, 0755); err != nil {
				return "", fmt.Errorf("failed to create dir %s: %v\n", destFilePath, err)
			}
			continue
		}

		// Root Directory Prefix Determination
		// * Extract directory from filepath of header
		dir := ""
		if idx := strings.LastIndex(header.Name, "/"); idx >= 0 {
			dir = header.Name[:idx]
		} else {
			dir = "."
		}
		directoryPrefixCount[dir]++
		// Root Directory Prefix Determination:
		// 1. Iterate through dir prefix count map[string]int
		// 2. Dir with most counts will most likely be the root extract dir
		dirMaxCount := 0
		for dir, count := range directoryPrefixCount {
			if count > dirMaxCount {
				dirMaxCount = count
				extractRootDir = dir
			}
		}

		// extractRootFullPath = fmt.Sprintf("%s/%s", destDir, extractRootDir)
		//// Check if destination dir exists
		//if _, err := os.Stat(extractRootFullPath); err == nil {
		//	log.Printf("removing previous extract dir - %s\n", extractRootFullPath)
		//	err := os.RemoveAll(extractRootFullPath)
		//	if err != nil {
		//		return "", fmt.Errorf("unable to delete existing file: %v", err)
		//	}
		//}

		// Create and open the destination file
		destFile, err := os.Create(destFilePath)
		if err != nil {
			return "", fmt.Errorf("failed to create %s: %v\n", destFilePath, err)
		}

		// Copy file contents from the tar archive to the destination file
		if _, err := io.Copy(destFile, tarReader); err != nil {
			return "", err
		}
		if err := destFile.Close(); err != nil {
			return "", cleanup(err)
		}
	}
	if err := gzipReader.Close(); err != nil {
		return "", cleanup(err)
	}
	if err := srcFileReader.Close(); err != nil {
		return "", cleanup(err)
	}
	return extractRootDir, nil
}

func SelectAndExtractTarGzFilesInDir(sourceDir string) (string, error) {
	var selectedFile os.DirEntry
	var sourceFilePath, extractRoot, extractedDebugPath string

	// If debug path is not a bundle directly, parse for bundles and extract
	if !strings.HasSuffix(sourceDir, ".tar.gz") {
		var bundles []os.DirEntry
		files, err := os.ReadDir(sourceDir)
		if err != nil {
			return "", fmt.Errorf("[select-and-extract] failed to read debug-path directory %s\n%v\n", sourceDir, err)
		}
		// Filter files for .tar.gz bundles
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".tar.gz") {
				bundles = append(bundles, file)
			}
		}
		fmt.Println("select a .tar.gz file to extract:")
		for i, bundle := range bundles {
			fmt.Printf("%d: %s\n", i+1, bundle.Name())
		}
		fmt.Print("enter the number of the file to extract: ")
		var selected int
		if _, err := fmt.Scanf("%d", &selected); err != nil {
			return "", err
		}

		if selected < 1 || selected > len(bundles) {
			return "", fmt.Errorf("invalid selection: %v", err)
		}

		selectedFile = bundles[selected-1]
		sourceFilePath = filepath.Join(sourceDir, selectedFile.Name())
	} else {
		sourceFilePath = sourceDir
	}
	extractRoot, err := extractTarGz(sourceFilePath, filepath.Dir(sourceFilePath))
	if err != nil {
		return "", fmt.Errorf("[select-and-extract] error extracting %s: %v\n", sourceFilePath, err)
	}

	if strings.HasSuffix(sourceDir, ".tar.gz") {
		sourceFilePath, _ = filepath.Abs(sourceFilePath)
		extractedDebugPath = filepath.Join(filepath.Dir(sourceFilePath), extractRoot)
	} else {
		extractedDebugPath = filepath.Join(sourceDir, extractRoot)
	}
	return extractedDebugPath, nil
}

func ConvertToValidJSON(input string) string {
	// Replace { with {" and } with "},"
	input = strings.Replace(input, "{", `{`, -1)
	input = strings.Replace(input, "}", `},`, -1)

	// Replace Suffrage:Voter with "Suffrage":"Voter",
	re := regexp.MustCompile(`(Suffrage):(\w+)`)
	input = re.ReplaceAllString(input, `"$1":"$2",`)

	// Replace ID:c24d7789-af04-7bca-2649-42ebe6a227a3 with "ID":"c24d7789-af04-7bca-2649-42ebe6a227a3",
	re = regexp.MustCompile(`(\w+):(\w+-+\w+-\w+-\w+-\w+)`)
	input = re.ReplaceAllString(input, `"$1":"$2",`)

	re = regexp.MustCompile(`(\w+):(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5})`)
	input = re.ReplaceAllString(input, `"$1":"$2"`)

	// Remove the trailing comma after the last object
	input = strings.Replace(input, ",]", `]`, 1)

	return input
}

func ConvertSecondsReadable(seconds int) string {
	// Calculate days, hours, minutes, and seconds
	days := seconds / 86400
	seconds %= 86400
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	sec := seconds % 60

	// Format the uptime in a human-readable way
	formatted := fmt.Sprintf("%d days, %d hours, %d minutes, %d seconds", days, hours, minutes, sec)

	return formatted
}

func StructToHCL(data interface{}, indent string) string {
	hcl := ""
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if value.Kind() == reflect.Struct {
			hcl += fmt.Sprintf("%s%s {\n", indent, field.Name)
			hcl += StructToHCL(value.Interface(), indent+"  ")
			hcl += fmt.Sprintf("%s}\n", indent)
		} else {
			jsonTagName := field.Tag.Get("json")
			if jsonTagName == "" {
				jsonTagName = field.Name
			}
			hcl += fmt.Sprintf("%s%s = %v\n", indent, jsonTagName, value.Interface())
		}
	}

	return hcl
}

func WriteFileWithPerms(outputFile, payload string, mode os.FileMode) error {
	// os.WriteFile truncates existing files and overwrites them, but only if they are writable.
	// If the file exists it will already likely be read-only. Remove it first.
	if _, err := os.Stat(outputFile); err == nil {
		if err = os.RemoveAll(outputFile); err != nil {
			return fmt.Errorf("unable to delete existing file: %s", err)
		}
	}
	if err := os.WriteFile(outputFile, []byte(payload), os.ModePerm); err != nil {
		return fmt.Errorf("unable to write file: %s", err)
	}
	return os.Chmod(outputFile, mode)
}

func (b *Debug) DecodeAgent(agentDecoder *json.Decoder) error {
	var agentConfig Agent
	err := agentDecoder.Decode(&agentConfig)
	if err != nil {
		log.Fatalf("error decoding agent: %v", err)
		return err
	}
	b.Agent = agentConfig
	return nil
}

func (b *Debug) DecodeMembers(memberDecoder *json.Decoder) error {
	var membersList []Member
	err := memberDecoder.Decode(&membersList)
	if err != nil {
		log.Fatalf("error decoding members: %v", err)
		return err
	}
	b.Members = membersList
	return nil
}

func (b *Debug) DecodeMetrics(metricsDecoder *json.Decoder) error {
	type MetricData struct {
		Timestamp string
		Metric    Metric
	}

	// Determine the number of CPU cores
	numCores := runtime.NumCPU()

	// Create a channel to receive metric data with a buffered size based on CPU cores
	bufferSize := numCores * 10000 // Tweak this as needed
	metricDataChan := make(chan MetricData, bufferSize)

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	for {
		var metric Metric
		err := metricsDecoder.Decode(&metric)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error decoding | file: metrics.json %v", err)
			return err
		}

		wg.Add(1)
		go func(metric Metric) {
			defer wg.Done()

			// Extract the timestamp from the metric
			timestamp := metric.Timestamp

			// Send metric data to the channel
			metricDataChan <- MetricData{Timestamp: timestamp, Metric: metric}
		}(metric)
	}

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(metricDataChan)
	}()

	// Create a map to organize metrics by timestamp
	metricsByTimestamp := make(map[string][]Metric)

	// Organize metrics by timestamp and build the map
	for metricData := range metricDataChan {
		timestamp := metricData.Timestamp
		metric := metricData.Metric

		// Append metric to the map
		metricsByTimestamp[timestamp] = append(metricsByTimestamp[timestamp], metric)
	}

	// Sort timestamps
	sortedTimestamps := make([]string, 0, len(metricsByTimestamp))
	for timestamp := range metricsByTimestamp {
		sortedTimestamps = append(sortedTimestamps, timestamp)
	}
	sort.Strings(sortedTimestamps)

	// Reconstruct the Metrics slice in timestamped order
	var sortedMetrics []Metric
	for _, timestamp := range sortedTimestamps {
		sortedMetrics = append(sortedMetrics, metricsByTimestamp[timestamp]...)
	}

	// Assign the sorted Metrics to the Debug struct
	b.Metrics.Metrics = sortedMetrics

	// Build the metrics map from the sorted metrics
	b.Metrics.BuildMetricsMap()

	return nil
}

func (b *Debug) DecodeMetricsIndex(indexDecoder *json.Decoder) error {
	var index Index
	err := indexDecoder.Decode(&index)
	if err != nil {
		log.Fatalf("error decoding metrics: %v", err)
		return err
	}
	b.Index = index
	return nil
}

func (b *Debug) DecodeHost(hostDecoder *json.Decoder) error {
	for {
		var hostObject Host
		err := hostDecoder.Decode(&hostObject)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error decoding host: %v", err)
			return err
		}
		b.Host = hostObject
	}
	return nil
}

func (b *Debug) DecodeJSON(debugPath, dataType string) error {
	configs := map[string]string{
		"agent":   "agent.json",
		"members": "members.json",
		"metrics": "metrics.json",
		"host":    "host.json",
		"index":   "index.json",
	}

	fileName, found := configs[dataType]
	if !found && dataType != "all" {
		return fmt.Errorf("unknown data type: %s", dataType)
	}

	if dataType == "all" {
		for dataType, fileName := range configs {
			if dataType == "all" {
				continue
			}
			if err := b.decodeFile(debugPath, fileName, dataType); err != nil {
				return err
			}
		}
		return nil
	}
	return b.decodeFile(debugPath, fileName, dataType)
}

func (b *Debug) decodeFile(debugPath, fileName, dataType string) error {
	filePath := fmt.Sprintf("%s/%s", debugPath, fileName)

	// Read the entire file into memory
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s. Ensure debug-path is set to a valid path\n", filePath)
		} else {
			return fmt.Errorf("error reading file: %s - %v\n", filePath, err)
		}
	}

	// Create a JSON decoder for the file data
	decoder := json.NewDecoder(bytes.NewReader(fileData))

	// Decode JSON based on the data type
	switch dataType {
	case "agent":
		return b.DecodeAgent(decoder)
	case "members":
		return b.DecodeMembers(decoder)
	case "metrics":
		return b.DecodeMetrics(decoder)
	case "host":
		return b.DecodeHost(decoder)
	case "index":
		return b.DecodeMetricsIndex(decoder)
	default:
		return fmt.Errorf("unknown data type: %s", dataType)
	}
}

type Map map[string][]map[string]interface{}

type ByValue []string

func (m ByValue) Len() int      { return len(m) }
func (m ByValue) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m ByValue) Less(i, j int) bool {
	columns_i := strings.Split(m[i], "\x1f")
	columns_j := strings.Split(m[j], "\x1f")
	var value_i, value_j float64
	if len(columns_i) >= 2 && len(columns_i) <= 4 {
		value_i, _ = strconv.ParseFloat(strings.TrimRight(columns_i[1], "%"), 64)
		value_j, _ = strconv.ParseFloat(strings.TrimRight(columns_j[1], "%"), 64)
	} else {
		value_i, _ = strconv.ParseFloat(strings.TrimRight(columns_i[4], "%"), 64)
		value_j, _ = strconv.ParseFloat(strings.TrimRight(columns_j[4], "%"), 64)
	}

	// using '>' vice '<' to sort from highest -> lowest
	return value_i > value_j
}

// nonNegativeDifference calculates the non-negative difference between two float64 values.
func nonNegativeDifference(a, b float64) float64 {
	diff := a - b
	if diff >= 0 {
		return diff
	}
	return 0 // Return the absolute value of the difference if < 0
}

// CalculateGCRate calculates the rate of Garbage Collection (GC) in nanoseconds per minute.
func CalculateGCRate(currentValueData, previousValueData []map[string]interface{}) (string, error) {
	var rate string

	currentValue, ok := currentValueData[0]["value"].(float64)
	if !ok {
		return "", fmt.Errorf("invalid 'value' field in data")
	}
	previousValue, ok := previousValueData[0]["value"].(float64)
	if !ok {
		return "", fmt.Errorf("invalid 'value' field in data")
	}
	// Calculate the non-negative difference in GC pause times
	diff := nonNegativeDifference(currentValue, previousValue)
	timeCurrent, err := time.Parse("2006-01-02 15:04:05 -0700 MST", fmt.Sprintf("%s", currentValueData[0]["timestamp"]))
	if err != nil {
		return "", err
	}
	timePrevious, err := time.Parse("2006-01-02 15:04:05 -0700 MST", fmt.Sprintf("%s", previousValueData[0]["timestamp"]))
	if err != nil {
		return "", err
	}
	// consul debug caputures default to 5m/30s capture intervals (>= v1.16.x)
	//
	timeDiff := timeCurrent.Sub(timePrevious).Seconds()
	if diff >= 0 && timeDiff > 0 {
		rate, err = ConvertToReadableTime(diff/(timeDiff/60), "ns") // convert to ns/min to most-readable-time/minute
		if err != nil {
			return "", err
		}
		rate = fmt.Sprintf("%s/min", rate)
	}
	if rate == "" {
		rate = "-"
	}
	return rate, nil
}

// ByteConverter
// Struct used to implement the ConvertToReadableBytes interface function for int and float64
// byte conversion.
type ByteConverter struct{}

func (bc ByteConverter) ConvertToReadableBytes(value interface{}) string {
	switch v := value.(type) {
	case int:
		return ConvertIntBytes(v)
	case float64:
		return ConvertFloatBytes(v)
	default:
		return "Unsupported type"
	}
}

func ConvertIntBytes(bytes int) string {
	const (
		kb int64 = 1024
		mb       = 1024 * kb
		gb       = 1024 * mb
		tb       = 1024 * gb
	)

	switch {
	case int64(bytes) >= tb:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(tb))
	case int64(bytes) >= gb:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(gb))
	case int64(bytes) >= mb:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(mb))
	case int64(bytes) >= kb:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

func ConvertFloatBytes(bytes float64) string {
	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
		tb = 1024 * gb
	)

	switch {
	case bytes >= tb:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(tb))
	case bytes >= gb:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%.4f bytes", bytes)
	}
}

// TimeConverter is the interface for converting time units.
type TimeConverter interface {
	Convert(timeValue interface{}) (string, error)
}

func ConvertToReadableTime(value interface{}, units string) (string, error) {
	var converter TimeConverter
	switch units {
	case "ns":
		converter = NanosecondsConverter{}
	case "ms":
		converter = MillisecondsConverter{}
	case "seconds":
		converter = SecondsConverter{}
	case "hours":
		converter = HoursConverter{}
	}
	v, err := converter.Convert(value)
	if err != nil {
		return "", err
	}
	return v, nil
}

// NanosecondsConverter implements TimeConverter for nanoseconds.
type NanosecondsConverter struct{}

func (n NanosecondsConverter) Convert(timeValue interface{}) (string, error) {
	const (
		nsInMs     = 1e6
		nsInSecond = 1e9
		nsInHour   = 3.6e12
	)

	switch v := timeValue.(type) {
	case int:
		switch {
		case int64(v) >= int64(nsInHour):
			return fmt.Sprintf("%.2fh", float64(v)/float64(nsInHour)), nil
		case int64(v) >= int64(nsInSecond):
			return fmt.Sprintf("%.2fs", float64(v)/float64(nsInSecond)), nil
		case int64(v) >= int64(nsInMs):
			return fmt.Sprintf("%.2fms", float64(v)/float64(nsInMs)), nil
		default:
			return fmt.Sprintf("%dns", v), nil
		}
	case float64:
		switch {
		case v >= nsInHour:
			return fmt.Sprintf("%.2fh", v/float64(nsInHour)), nil
		case v >= nsInSecond:
			return fmt.Sprintf("%.2fs", v/float64(nsInSecond)), nil
		case v >= nsInMs:
			return fmt.Sprintf("%.2fms", v/float64(nsInMs)), nil
		default:
			return fmt.Sprintf("%.4fns", v), nil
		}
	default:
		return "", fmt.Errorf("unsupported type: %v", reflect.TypeOf(timeValue))
	}
}

// MillisecondsConverter implements TimeConverter for milliseconds.
type MillisecondsConverter struct{}

func (m MillisecondsConverter) Convert(timeValue interface{}) (string, error) {
	const (
		msInSecond = 1e3
		msInHour   = 3.6e6
	)

	switch v := timeValue.(type) {
	case int:
		switch {
		case v >= msInHour:
			return fmt.Sprintf("%.2fh", float64(v)/float64(msInHour)), nil
		case v >= msInSecond:
			return fmt.Sprintf("%.2fs", float64(v)/float64(msInSecond)), nil
		default:
			return fmt.Sprintf("%.4fms", float64(v)), nil
		}
	case float64:
		switch {
		case v >= msInHour:
			return fmt.Sprintf("%.2fh", v/float64(msInHour)), nil
		case v >= msInSecond:
			return fmt.Sprintf("%.2fs", v/float64(msInSecond)), nil
		default:
			return fmt.Sprintf("%.4fms", v), nil
		}
	default:
		return "", fmt.Errorf("unsupported type: %v", reflect.TypeOf(timeValue))
	}
}

// SecondsConverter implements TimeConverter for seconds.
type SecondsConverter struct{}

func (s SecondsConverter) Convert(timeValue interface{}) (string, error) {
	const (
		secondsInHour = 3600
	)

	switch v := timeValue.(type) {
	case int:
		switch {
		case v >= secondsInHour:
			return fmt.Sprintf("%.2fh", float64(v)/float64(secondsInHour)), nil

		default:
			return fmt.Sprintf("%.2fs", float64(v)), nil
		}
	case float64:
		switch {
		case v >= secondsInHour:
			return fmt.Sprintf("%.2fh", v/float64(secondsInHour)), nil
		default:
			return fmt.Sprintf("%.2fs", v), nil
		}
	default:
		return "", fmt.Errorf("unsupported type: %v", reflect.TypeOf(timeValue))
	}
}

// HoursConverter implements TimeConverter for hours.
type HoursConverter struct{}

func (h HoursConverter) Convert(timeValue interface{}) (string, error) {
	switch v := timeValue.(type) {
	case int:
		return fmt.Sprintf("%.2fh", float64(v)), nil
	case float64:
		return fmt.Sprintf("%.2fh", v), nil
	default:
		return "", fmt.Errorf("unsupported type: %v", reflect.TypeOf(timeValue))
	}
}

func validateName(name string, info string) bool {
	// This metric name is dynamic and can be anything that the customer uses for service names
	reg := regexp.MustCompile(`^consul\.proxy\..+$`)
	if reg.MatchString(name) {
		fmt.Printf("built-in mesh proxy prefix used: %s\n", name)
		return true
	}
	// list of metrics contains the name somewhere, return with no error
	if strings.Contains(info, name) {
		return true
	}
	return false
}
