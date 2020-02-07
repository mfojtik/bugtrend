package dashboard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/go-echarts/go-echarts/charts"

	"github.com/mfojtik/bugtrend/pkg/report"
)

func getBugTypesChart(config DashboardConfig) (*charts.Line, error) {
	reportFile, err := ioutil.ReadFile(config.BugTypesSeriesFile)
	if err != nil {
		return nil, err
	}
	var series report.BugTypeList
	if err := json.Unmarshal(reportFile, &series); err != nil {
		return nil, err
	}

	resolutions := []string{
		"Flakes",
		"AutomationBlocker",
		"TestBlocker",
		"Upgrade",
		"Other",
	}

	perStateSeries := map[string][]float64{}
	var xValues []string
	for _, serie := range series {
		xValues = append(xValues, fmt.Sprintf("%02d/%02d/%d %02d:%02d", serie.Timestamp.Month(), serie.Timestamp.Day(), serie.Timestamp.Year(), serie.Timestamp.Hour(), serie.Timestamp.Minute()))
		perStateSeries["Flakes"] = append(perStateSeries["Flakes"], float64(serie.Counts.Flakes))
		perStateSeries["AutomationBlocker"] = append(perStateSeries["AutomationBlocker"], float64(serie.Counts.AutomationBlocker))
		perStateSeries["TestBlocker"] = append(perStateSeries["TestBlocker"], float64(serie.Counts.TestBlocker))
		perStateSeries["Upgrade"] = append(perStateSeries["Upgrade"], float64(serie.Counts.Upgrade))
		perStateSeries["Other"] = append(perStateSeries["Other"], float64(serie.Counts.Other))
	}

	c := charts.NewLine()
	c.SetGlobalOptions(charts.TitleOpts{Title: fmt.Sprintf("%s", config.Release)})
	c.AddXAxis(xValues)
	for _, state := range resolutions {
		c.AddYAxis(state, perStateSeries[state])
	}
	return c, nil
}
