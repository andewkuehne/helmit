package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/andrewkuehne/helmit/charts"
)

func TestMain(m *testing.M) {
	// Keep track of the original command-line arguments
	oldArgs := os.Args

	// Create a buffer to capture the output of the program
	var buf bytes.Buffer

	// Override the command-line arguments to test different cases
	testCases := []struct {
		args   []string
		output string
	}{
		{
			args:   []string{"helmit", "-h"},
			output: usage + "\n",
		},
		{
			args:   []string{"helmit", "--help"},
			output: usage + "\n",
		},
		{
			args:   []string{"helmit", "-chart", "non-existent-chart"},
			output: "Error: open non-existent-chart: no such file or directory\n",
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		// Override the command-line arguments with the test case
		os.Args = tc.args

		// Override the standard output to capture the program's output
		os.Stdout = (*os.File)(unsafe.Pointer(&buf))

		// Run the program
		main()

		// Check if the output matches the expected value
		if buf.String() != tc.output {
			fmt.Printf("Unexpected output. Expected: %q, got: %q", tc.output, buf.String())
		}

		// Reset the standard output and buffer
		os.Stdout = os.NewFile(1, os.DevNull)
		buf.Reset()
	}

	// Restore the original command-line arguments
	os.Args = oldArgs

	// Run the LoadChart test
	if err := testLoadChart(); err != nil {
		fmt.Printf("TestLoadChart failed: %v", err)
	}

	// Run the tests
	exitCode := m.Run()

	// Exit with the test's exit code
	os.Exit(exitCode)
}

func testLoadChart() error {
	// Create a temporary directory to hold the chart files
	chartDir, err := ioutil.TempDir("", "test-chart")
	if err != nil {
		return fmt.Errorf("failed to create temp chart directory: %v", err)
	}
	defer os.RemoveAll(chartDir)

	// Create a simple Helm chart
	chartName := "test-chart"
	chartVersion := "0.1.0"
	chartPath := filepath.Join(chartDir, chartName)
	chartYaml := fmt.Sprintf("apiVersion: v1\nname: %s\nversion: %s\ndescription: A Helm chart for testing\n", chartName, chartVersion)
	chartFile := filepath.Join(chartPath, "Chart.yaml")
	if err := os.Mkdir(chartPath, 0777); err != nil {
		return fmt.Errorf("failed to create chart directory: %v", err)
	}
	if err := ioutil.WriteFile(chartFile, []byte(chartYaml), 0666); err != nil {
		return fmt.Errorf("failed to create chart file: %v", err)
	}

	// Load the chart and check the metadata
	chart, err := charts.LoadChart(chartPath)
	if err != nil {
		return fmt.Errorf("failed to load chart: %v", err)
	}
	if chart.Metadata.Name != chartName {
		return fmt.Errorf("expected chart name to be %q, got %q", chartName, chart.Metadata.Name)
	}
	if chart.Metadata.Version != chartVersion {
		return fmt.Errorf("expected chart version to be %q, got %q", chartVersion, chart.Metadata.Version)
	}
	return nil
}
