package charts

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/chart"
)

func TestLoadChart(t *testing.T) {
	t.Run("valid chart", func(t *testing.T) {
		// Create a temporary directory and a chart in it
		tmpDir, err := ioutil.TempDir("", "test")
		assert.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		chartPath := filepath.Join(tmpDir, "mychart")
		err = os.Mkdir(chartPath, 0755)
		assert.NoError(t, err)

		err = ioutil.WriteFile(filepath.Join(chartPath, "Chart.yaml"), []byte("name: mychart\nversion: 1.0.0\n"), 0644)
		assert.NoError(t, err)

		err = ioutil.WriteFile(filepath.Join(chartPath, "templates", "mytemplate.yaml"), []byte("content"), 0644)
		assert.NoError(t, err)

		// Load the chart
		c, err := LoadChart(chartPath)
		assert.NoError(t, err)

		// Check the chart metadata and templates
		assert.Equal(t, "mychart", c.Metadata.Name)
		assert.Equal(t, "1.0.0", c.Metadata.Version)
		assert.Len(t, c.Templates, 1)
		assert.Equal(t, "mytemplate.yaml", c.Templates[0].Name)
		assert.Equal(t, "content", string(c.Templates[0].Data))
	})

	t.Run("invalid chart", func(t *testing.T) {
		// Create a temporary file and write some invalid data to it
		tmpFile, err := ioutil.TempFile("", "test")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.Write([]byte("invalid chart data"))
		assert.NoError(t, err)

		// Load the chart and check that it returns an error
		_, err = LoadChart(tmpFile.Name())
		assert.Error(t, err)
	})
}

func TestLintChart(t *testing.T) {
	t.Run("valid chart", func(t *testing.T) {
		// Create a valid chart
		c := &chart.Chart{
			Metadata: &chart.Metadata{
				Name:    "mychart",
				Version: "1.0.0",
			},
			Templates: []*chart.File{
				{Name: "mytemplate.yaml", Data: []byte("content")},
			},
		}

		// Lint the chart and check that it doesn't return an error
		err := lintChart(c)
		assert.NoError(t, err)
	})

	t.Run("invalid chart", func(t *testing.T) {
		// Create an invalid chart with missing metadata
		c := &chart.Chart{
			Templates: []*chart.File{
				{Name: "mytemplate.yaml", Data: []byte("content")},
			},
		}

		// Lint the chart and check that it returns an error
		err := lintChart(c)
		assert.Error(t, err)
	})
}
