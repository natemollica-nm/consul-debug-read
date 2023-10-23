package lib

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func TimeStampDuration(timeStart, timeStop string) time.Duration {
	start, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", timeStart)
	stop, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", timeStop)
	return stop.Sub(start)
}

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

// GetExtractName
// Returns tar rendered root directory name
func GetExtractName(bundleName string) string {
	var capturedPart string

	// Define a regular expression pattern to capture the desired part of the string.
	pattern := `^\d+(.*?)(?:\.tar\.gz)$`

	// Compile the regular expression.
	re := regexp.MustCompile(pattern)

	// Find the submatch using FindStringSubmatch.
	submatches := re.FindStringSubmatch(bundleName)

	if submatches != nil && len(submatches) > 1 {
		// Extract and print the captured part of the string.
		capturedPart = submatches[1]
	}
	return capturedPart
}

// GetMostRecentFile returns the path of the most recently modified file in a directory.
func GetMostRecentFile(directory string) (string, error) {
	var mostRecentFile string
	var mostRecentTime time.Time
	var returnFile string
	var debugExtracts []os.DirEntry

	files, err := os.ReadDir(directory)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), "consul-debug-") {
			debugExtracts = append(debugExtracts, file)
		}
	}

	fmt.Printf("\nSelect a consul debug extract directory to use:\n")
	for i, file := range debugExtracts {
		if file.IsDir() && strings.HasPrefix(file.Name(), "consul-debug-") {
			info, _ := file.Info()
			modTime := info.ModTime()
			if modTime.After(mostRecentTime) {
				mostRecentTime = modTime
				mostRecentFile = filepath.Join(directory, file.Name())
			}
			fmt.Printf("%d: %s\n", i+1, file.Name())
		}
	}
	fmt.Printf("[*] Most recently extract: '%s'\n", mostRecentFile)

	fmt.Print("Enter the number of debug extract directory to use: ")
	var selected int
	if _, err := fmt.Scanf("%d", &selected); err != nil {
		return "", err
	}

	if selected < 0 || selected > len(files) {
		return "", fmt.Errorf("invalid selection")
	}

	switch {
	case selected == 0:
		returnFile = mostRecentFile
	case selected >= 1 && selected < len(files):
		selectedFile := debugExtracts[selected-1]
		returnFile = filepath.Join(directory, selectedFile.Name())
	}

	return returnFile, nil
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

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(srcFileReader)
	if err != nil {
		return "", err
	}
	defer gzipReader.Close()

	// Create a tar reader
	tarReader := tar.NewReader(gzipReader)

	destFileName := GetExtractName(filepath.Base(srcFile))
	destFilePath := fmt.Sprintf("%s/%s", destDir, destFileName)
	log.Printf("destination File Extract Path - %s\n", destFilePath)
	// Check if destination dir exists
	if _, err := os.Stat(destFilePath); err == nil {
		log.Printf("removing previous extract dir - %s\n", destFilePath)
		err := os.RemoveAll(destFilePath)
		if err != nil {
			return "", fmt.Errorf("unable to delete existing file: %v", err)
		}
	}

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

		// Create and open the destination file
		destFile, err := os.Create(destFilePath)
		if err != nil {
			return "", fmt.Errorf("failed to create %s: %v\n", destFilePath, err)
		}
		defer destFile.Close()

		// Copy file contents from the tar archive to the destination file
		if _, err := io.Copy(destFile, tarReader); err != nil {
			return "", err
		}
	}
	log.Printf("debug-extract: root extraction directory is %s\n", extractRootDir)
	return extractRootDir, nil
}

// SelectAndExtractTarGzFilesInDir allows the user to interactively select and extract a .tar.gz file from a directory.
func SelectAndExtractTarGzFilesInDir(sourceDir string) (string, error) {
	var selectedFile os.DirEntry
	var sourceFilePath string
	var extractRoot string

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

	log.Printf("[select-and-extract] extracting %s\n", sourceFilePath)
	extractRoot, err := extractTarGz(sourceFilePath, filepath.Dir(sourceFilePath))

	if err != nil {
		return "", fmt.Errorf("[select-and-extract] error extracting %s: %v\n", sourceFilePath, err)
	}
	extractedDebugPath := filepath.Join(sourceDir, extractRoot)

	log.Printf("[select-and-extract] extraction of %s completed successfully!\n", sourceFilePath)
	log.Printf("[select-and-extract] setting debug path to %s\n", extractedDebugPath)
	return extractedDebugPath, nil
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
		return fmt.Sprintf("%d bytes", bytes)
	}
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

func ExecuteJQ(data string, jqFilter string) (string, error) {
	cmd := exec.Command("jq", "--raw-output", jqFilter)

	// Create pipes for stdin, stdout, and stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}

	// Start the jq command
	if err := cmd.Start(); err != nil {
		return "", err
	}

	// Write data to stdin
	_, err = io.WriteString(stdin, data)
	if err != nil {
		return "", err
	}
	stdin.Close()

	// Read the result from stdout
	result, err := io.ReadAll(stdout)
	if err != nil {
		return "", err
	}

	// Read and display any errors from stderr
	errorOutput, err := io.ReadAll(stderr)
	if err != nil {
		return "", err
	}

	if len(errorOutput) > 0 {
		return "", fmt.Errorf("jq error: %s", string(errorOutput))
	}

	// Wait for the jq command to complete
	if err := cmd.Wait(); err != nil {
		return "", err
	}

	return string(result), nil
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

// WriteFileWithPerms will write payload as the contents of the outputFile and set permissions after writing the contents. This function is necessary since using os.WriteFile() alone will create the new file with the requested permissions prior to actually writing the file, so you can't set read-only permissions.
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
