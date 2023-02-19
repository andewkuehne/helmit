package k8stestinit

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const kubeconfigDir = "/tmp"

// InitTestEnv initializes a temporary test environment for helmit.
// It will test the kubeconfig to ensure that it is valid and working,
// and if it is, it will write the kubeconfig to a temporary file for
// use in the `--test` flag of the helmit command.
func InitTestEnv(kubeconfigPath string) (errorMessage error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		errorMessage = fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		errorMessage = fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	_, err = clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		errorMessage = fmt.Errorf("failed to list namespaces: %w", err)
	}

	// Write kubeconfig to temp file.
	kubeconfigBytes, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		errorMessage = fmt.Errorf("failed to read kubeconfig file: %w", err)
	}

	now := time.Now().Format("20060102-150405")
	kubeconfigTempPath := filepath.Join(kubeconfigDir, fmt.Sprintf("kubeconfig-%s", now))
	err = ioutil.WriteFile(kubeconfigTempPath, kubeconfigBytes, 0600)
	if err != nil {
		errorMessage = fmt.Errorf("failed to write kubeconfig to temp file: %w", err)
	}

	fmt.Printf("Successfully wrote kubeconfig to %s\n", kubeconfigTempPath)

	return nil
}
