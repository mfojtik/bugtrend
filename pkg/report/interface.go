package report

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type ReportWriter interface {
	Write() ([]byte, error)
}

type TimeSerieWriterConfig struct {
	Name      string
	Release   string
	Timestamp time.Time
}

func NewTimeSerieWriter(name, release string, timestamp time.Time) *TimeSerieWriterConfig {
	return &TimeSerieWriterConfig{
		Name:      name,
		Release:   release,
		Timestamp: timestamp,
	}
}

func (c *TimeSerieWriterConfig) WriteTimeSerie(w ReportWriter) error {
	burnDownBytes, err := w.Write()
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(fmt.Sprintf("reports/%s/%s_%d.json", c.Release, c.Name, c.Timestamp.Unix()), burnDownBytes, os.ModePerm); err != nil {
		return err
	}
	return c.aggregatedTimeSeries()
}

func (c *TimeSerieWriterConfig) aggregatedTimeSeries() error {
	aggregatedFileContent := ""
	err := filepath.Walk(path.Join("reports", c.Release), func(p string, info os.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".json") {
			return nil
		}
		if info.Name() == fmt.Sprintf("%s.json", c.Name) {
			return nil
		}
		if !strings.HasPrefix(info.Name(), c.Name+"_") {
			return nil
		}
		content, err := ioutil.ReadFile(path.Join("reports", c.Release, info.Name()))
		if len(aggregatedFileContent) == 0 {
			aggregatedFileContent = string(content)
			return nil
		}
		aggregatedFileContent = strings.Join([]string{aggregatedFileContent, string(content)}, ",\n")
		return nil
	})
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join("reports", c.Release, c.Name+".json"), []byte("["+aggregatedFileContent+"]"), os.ModePerm)
}
