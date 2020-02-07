package dashboard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/go-echarts/go-echarts/charts"

	"github.com/mfojtik/bugtrend/pkg/report"
)

func getBugBurndownChart(config DashboardConfig) (*charts.Line, error) {
	reportFile, err := ioutil.ReadFile(config.BurndownSeriesFile)
	if err != nil {
		return nil, err
	}
	var series report.BurndownList
	if err := json.Unmarshal(reportFile, &series); err != nil {
		return nil, err
	}

	states := []string{
		"NEW",
		"ASSIGNED",
		"POST",
		"MODIFIED",
		"ON_QA",
		"VERIFIED",
		"CLOSED",
	}

	perStateSeries := map[string][]float64{}
	var xValues []string
	for _, serie := range series {
		xValues = append(xValues, fmt.Sprintf("%02d/%02d/%d %02d:%02d", serie.Timestamp.Month(), serie.Timestamp.Day(), serie.Timestamp.Year(), serie.Timestamp.Hour(), serie.Timestamp.Minute()))
		for _, state := range states {
			found := false
			for _, c := range serie.Counts {
				if c.Status == state {
					found = true
					perStateSeries[state] = append(perStateSeries[state], float64(c.Count))
					break
				}
			}
			if !found {
				perStateSeries[state] = append(perStateSeries[state], float64(0))
			}
		}
	}
	c := charts.NewLine()
	c.SetGlobalOptions(charts.TitleOpts{Title: fmt.Sprintf("%s", config.Release)})
	c.AddXAxis(xValues)
	for _, state := range states {
		c.AddYAxis(state, perStateSeries[state])
	}
	c.Subtitle = fmt.Sprintf("Total: %d", series[len(series)-1].Total)
	return c, nil
}
