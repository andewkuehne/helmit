package k8stesting

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestHelmChart(chartPath string) error {
	kubeconfigPath := "/tmp/helmit_kubeconfig"
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return fmt.Errorf("error building config from flags: %v", err)
	}

	_, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("error creating kubernetes client: %v", err)
	}

	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Printf(format, v)
	}); err != nil {
		return fmt.Errorf("error initializing action configuration: %v", err)
	}

	var chart *chart.Chart

	// detect chart type and load it accordingly
	if strings.HasSuffix(chartPath, ".tgz") || strings.HasSuffix(chartPath, ".tar.gz") {
		f, err := os.Open(chartPath)
		if err != nil {
			return fmt.Errorf("error opening chart file: %v", err)
		}
		defer f.Close()

		gzf, err := gzip.NewReader(f)
		if err != nil {
			return fmt.Errorf("error creating gzip reader: %v", err)
		}
		defer gzf.Close()

		tr := tar.NewReader(gzf)

		// Use a buffer to read the tar file
		var buf bytes.Buffer
		_, err = io.Copy(&buf, tr)
		if err != nil {
			return fmt.Errorf("error reading chart bytes: %v", err)
		}

		chart, err = loader.LoadArchive(bytes.NewReader(buf.Bytes()))
		if err != nil {
			return fmt.Errorf("error loading chart archive: %v", err)
		}
	} else if strings.HasSuffix(chartPath, ".yaml") || strings.HasSuffix(chartPath, ".yml") {
		chart, err = loader.LoadArchive(bytes.NewReader([]byte(chartPath)))
		if err != nil {
			return fmt.Errorf("error loading chart yaml: %v", err)
		}
	} else {
		chart, err = loader.Load(chartPath)
		if err != nil {
			return fmt.Errorf("error loading chart: %v", err)
		}
	}

	// Set up install action
	installAction := action.NewInstall(actionConfig)
	installAction.Namespace = "helmit-test-release"
	installAction.ReleaseName = fmt.Sprintf("helmit-test-release-%s", time.Now().Format("20060102150405"))
	installAction.CreateNamespace = true
	installAction.Timeout = 5 * time.Minute

	// Install chart
	_, err = installAction.Run(chart, nil)
	if err != nil {
		return fmt.Errorf("error installing chart: %v", err)
	}

	// Verify release
	listAction := action.NewList(actionConfig)
	listAction.Filter = installAction.ReleaseName
	releaseList, err := listAction.Run()
	if err != nil {
		return fmt.Errorf("error listing releases: %v", err)
	}

	if len(releaseList) != 1 {
		return fmt.Errorf("expected 1 release, found %d", len(releaseList))
	}

	if releaseList[0].Info.Status != release.StatusDeployed {
		return fmt.Errorf("release not deployed, status: %v", releaseList[0].Info.Status)
	}

	// Uninstall chart
	uninstallAction := action.NewUninstall(actionConfig)
	_, err = uninstallAction.Run(installAction.ReleaseName)
	if err != nil {
		return fmt.Errorf("error uninstalling chart: %v", err)
	}

	return nil
}
