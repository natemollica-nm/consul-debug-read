package lib

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func ClearScreen() error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func TimeStampDuration(timeStart, timeStop string) time.Duration {
	start, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", timeStart)
	stop, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", timeStop)
	return stop.Sub(start)
}

// GetExtractName: returns tar rendered root file name
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

func extractTarGz(srcFile, destDir string) error {
	// Open the source .tar.gz file
	srcFileReader, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer srcFileReader.Close()

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(srcFileReader)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// Create a tar reader
	tarReader := tar.NewReader(gzipReader)

	destFileName := GetExtractName(filepath.Base(srcFile))
	destFilePath := fmt.Sprintf("%s/%s", destDir, destFileName)
	fmt.Printf("Destination File Name: %s\n", destFileName)
	fmt.Printf("Destination File Path: %s\n", destFilePath)
	// Check if destination dir exists
	if _, err := os.Stat(destFilePath); err == nil {
		fmt.Printf("Removing previous extract dir: %s\n", destFilePath)
		err := os.RemoveAll(destFilePath)
		if err != nil {
			return fmt.Errorf("unable to delete existing file: %v", err)
		}
	}

	// Iterate through the tar archive and extract files
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}

		// Calculate the file path for extraction
		destFilePath := fmt.Sprintf("%s/%s", destDir, header.Name)

		// Create directories as needed
		if header.FileInfo().IsDir() {
			if err := os.MkdirAll(destFilePath, 0755); err != nil {
				return err
			}
			continue
		}

		// Create and open the destination file
		destFile, err := os.Create(destFilePath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		// Copy file contents from the tar archive to the destination file
		if _, err := io.Copy(destFile, tarReader); err != nil {
			return err
		}
	}

	return nil
}

// SelectAndExtractTarGzFilesInDir allows the user to interactively select and extract a .tar.gz file from a directory.
func SelectAndExtractTarGzFilesInDir(sourceDir string) (string, error) {
	var selectedFile os.DirEntry
	var sourceFilePath string
	var extractedBundleDir string

	// If debug path is not a bundle directly, parse for bundles and extract
	if !strings.HasSuffix(sourceDir, ".tar.gz") {
		var bundles []os.DirEntry
		files, err := os.ReadDir(sourceDir)
		if err != nil {
			return "", err
		}
		// Filter files for .tar.gz bundles
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".tar.gz") {
				bundles = append(bundles, file)
			}
		}

		fmt.Println("Select a .tar.gz file to extract:")
		for i, bundle := range bundles {
			fmt.Printf("%d: %s\n", i+1, bundle.Name())
		}
		fmt.Print("Enter the number of the file to extract: ")
		var selected int
		if _, err := fmt.Scanf("%d", &selected); err != nil {
			return "", err
		}

		if selected < 1 || selected > len(bundles) {
			return "", fmt.Errorf("invalid selection: %v", err)
		}

		selectedFile = bundles[selected-1]
		sourceFilePath = filepath.Join(sourceDir, selectedFile.Name())
		extractedBundleDir = GetExtractName(filepath.Base(selectedFile.Name()))
	} else {
		sourceFilePath = sourceDir
		extractedBundleDir = GetExtractName(filepath.Base(sourceFilePath))
	}

	fmt.Printf("Extracting: %s\n", sourceFilePath)
	if err := extractTarGz(sourceFilePath, filepath.Dir(sourceFilePath)); err != nil {
		fmt.Printf("Error extracting %s: %v\n", sourceFilePath, err)
		return "", err
	}

	fmt.Printf("Extraction of %s completed successfully.\n", sourceFilePath)
	return extractedBundleDir, nil
}

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
