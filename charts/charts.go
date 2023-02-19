package charts

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
)

// LoadChart loads a Helm chart from the given path, which can be a YAML file, a directory, or a compressed directory.
func LoadChart(chartPath string) (*chart.Chart, error) {
	// Determine the type of the chart
	fileInfo, err := os.Stat(chartPath)
	if err != nil {
		return nil, err
	}

	var chart *chart.Chart

	if fileInfo.IsDir() {
		// Load the chart directory
		files, err := ioutil.ReadDir(chartPath)
		if err != nil {
			return nil, err
		}

		dirFiles := make([]*loader.BufferedFile, 0)

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			path := filepath.Join(chartPath, file.Name())

			content, err := ioutil.ReadFile(path)
			if err != nil {
				return nil, err
			}

			dirFiles = append(dirFiles, &loader.BufferedFile{Name: file.Name(), Data: content})
		}

		chart, err = loader.LoadFiles(dirFiles)
		if err != nil {
			return nil, err
		}
	} else {
		// Load the chart archive
		file, err := os.Open(chartPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		content, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		chart, err = loader.LoadArchive(bytes.NewReader(content))
		if err != nil {
			return nil, err
		}
	}

	// Validate the chart
	if err := chart.Validate(); err != nil {
		return nil, err
	}

	// Lint the chart
	if err := lintChart(chart); err != nil {
		return nil, err
	}

	return chart, nil
}

// lintChart lints the given chart and returns an error if there are any linting issues.
func lintChart(c *chart.Chart) error {
	if c.Metadata == nil {
		return errors.New("chart has no metadata")
	}

	if c.Metadata.Name == "" {
		return errors.New("chart has no name")
	}

	if c.Metadata.Version == "" {
		return errors.New("chart has no version")
	}

	if c.Templates == nil {
		return errors.New("chart has no templates")
	}

	if len(c.Templates) == 0 {
		return errors.New("chart has no templates")
	}

	return nil
}
