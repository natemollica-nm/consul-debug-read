package lib

import (
	"archive/tar"
	"compress/gzip"
	"consul-debug-read/cmd/config"
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

func ClearScreen() {
	clearScreen := exec.Command("clear")
	clearScreen.Stdout = os.Stdout
	_ = clearScreen.Run()
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
	if config.Verbose {
		log.Printf("[extract-tar-gz]: root extraction directory is %s\n", extractRootDir)
	}
	return extractRootDir, nil
}

// SelectAndExtractTarGzFilesInDir allows the user to interactively select and extract a .tar.gz file from a directory.
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
	if config.Verbose {
		log.Printf("[select-and-extract] extracting %s\n", sourceFilePath)
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
	if config.Verbose {
		log.Printf("[select-and-extract] extraction of %s completed successfully!\n", sourceFilePath)
		log.Printf("[select-and-extract] setting debug path to %s\n", extractedDebugPath)
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

func Dots(msg string, ch <-chan bool) {
	dots := []string{".", "..", "...", "...."}
	i := 0
	for {
		select {
		case <-ch:
			fmt.Print("\r")                                         // Carriage return to the beginning of the line
			fmt.Print(fmt.Sprintf(strings.Repeat(" ", len(msg)+4))) // Overwrite the line with spaces
			fmt.Print("\r")                                         // Carriage return again to the beginning of the line
			return
		default:
			fmt.Printf("%s", msg)
			fmt.Print(dots[i%len(dots)], "\r")
			i++
			time.Sleep(300 * time.Millisecond)
		}
	}
}
