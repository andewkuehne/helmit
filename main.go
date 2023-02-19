package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	charts "github.com/andrewkuehne/helmit/charts"
	k8stesting "github.com/andrewkuehne/helmit/k8s/k8stesting"
	k8stestinit "github.com/andrewkuehne/helmit/k8s/k8stestinit"
)

const usage = `
Usage: helmit [flags]
Flags:
  -c, --chart string          Path to the Helm chart
  -h, --help                  Output usage information
      --inittestenv string    Initialize a test environment with the given kubeconfig file
  -t, --test                  Test the Helm chart after loading`

func main() {
	// Command line flags
	chartPath := flag.String("chart", "", "Path to the Helm chart")
	helpFlag := flag.Bool("help", false, "Output usage information")
	initTestEnv := flag.String("inittestenv", "", "Initialize a test environment with the given kubeconfig file")
	testFlag := flag.Bool("test", false, "Test the Helm chart after loading")

	flag.Usage = func() {
		fmt.Println(usage)
	}

	flag.Parse()

	if *helpFlag || flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(0)
	}

	if *chartPath == "" {
		log.Fatal("Error: chart path must be specified")
	}

	if *initTestEnv != "" {
		err := k8stestinit.InitTestEnv(*initTestEnv)
		if err != nil {
			log.Fatalf("Failed to initialize test environment: %v", err)
		}
		fmt.Println("Test environment initialized.")
		os.Exit(0)
	}

	if *testFlag {
		err := k8stesting.TestHelmChart(*chartPath)
		if err != nil {
			log.Fatalf("Failed to test chart: %v", err)
		}
		fmt.Println("Chart tested successfully.")
		os.Exit(0)
	}

	// Load the Helm chart
	chart, err := charts.LoadChart(*chartPath)
	if err != nil {
		log.Fatalf("Failed to load chart: %v", err)
	}

	// Print the chart details
	fmt.Printf("Chart Name: %s\n", chart.Metadata.Name)
	fmt.Printf("Chart Description: %s\n", chart.Metadata.Description)
	fmt.Printf("Chart Version: %s\n", chart.Metadata.Version)
}
