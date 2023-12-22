package metrics

import (
	"bytes"
	"consul-debug-read/internal/read"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkDecodeMetrics(b *testing.B) {
	// Load your JSON data into a buffer or file
	repoRoot, err := filepath.Abs("../../../../")
	if err != nil {
		b.Fatalf("Error determining repository root path: %v", err)
	}
	testFile := fmt.Sprintf("%s/%s", repoRoot, "bundles/consul-debug-2023-12-04T22-53-46-0500/metrics.json")
	jsonData, err := os.ReadFile(testFile)
	if err != nil {
		b.Fatalf("Error reading test metrics.json file: %v", err)
	}

	for i := 0; i < b.N; i++ {
		// Create a new buffer with JSON data (you can use a file reader in your real scenario)
		buffer := bytes.NewBuffer(jsonData)

		// Create a JSON decoder
		decoder := json.NewDecoder(buffer)

		// Create a Debug instance to decode into
		debug := &read.Debug{}

		// Run the benchmarked function
		err := debug.DecodeMetrics(decoder)
		if err != nil {
			b.Fatalf("Error in DecodeMetrics: %v", err)
		}
	}
}
