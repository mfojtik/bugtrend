package dashboard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/go-echarts/go-echarts/charts"

	"github.com/mfojtik/bugtrend/pkg/report"
)

func getClosedChart(config DashboardConfig) (*charts.Line, error) {
	reportFile, err := ioutil.ReadFile(config.ClosedSeriesFile)
	if err != nil {
		return nil, err
	}
	var series report.ClosedList
	if err := json.Unmarshal(reportFile, &series); err != nil {
		return nil, err
	}

	resolutions := []string{
		"CANTFIX",
		"CURRENTRELEASE",
		"DEFERRED",
		"DUPLICATE",
		"ERRATA",
		"INSUFFICIENT_DATA",
		"NOTABUG",
		"UPSTREAM",
		"WONTFIX",
	}

	perStateSeries := map[string][]float64{}
	var xValues []string
	for _, serie := range series {
		xValues = append(xValues, fmt.Sprintf("%02d/%02d/%d %02d:%02d", serie.Timestamp.Month(), serie.Timestamp.Day(), serie.Timestamp.Year(), serie.Timestamp.Hour(), serie.Timestamp.Minute()))
		for _, resolution := range resolutions {
			found := false
			for _, c := range serie.Counts {
				if c.Resolution == resolution {
					found = true
					perStateSeries[resolution] = append(perStateSeries[resolution], float64(c.Count))
					break
				}
			}
			if !found {
				perStateSeries[resolution] = append(perStateSeries[resolution], float64(0))
			}
		}
	}
	c := charts.NewLine()
	c.SetGlobalOptions(charts.TitleOpts{Title: fmt.Sprintf("%s", config.Release)})
	c.AddXAxis(xValues)
	for _, state := range resolutions {
		c.AddYAxis(state, perStateSeries[state])
	}
	return c, nil
}
