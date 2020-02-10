package dashboard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/go-echarts/go-echarts/charts"

	"github.com/mfojtik/bugtrend/pkg/report"
)

func getComponentsChart(config DashboardConfig) (*charts.Line, error) {
	reportFile, err := ioutil.ReadFile(config.ComponentSeriesFile)
	if err != nil {
		return nil, err
	}
	var series report.ComponentList
	if err := json.Unmarshal(reportFile, &series); err != nil {
		return nil, err
	}

	perComponentSeries := map[string][]float64{}
	var xValues []string
	for _, serie := range series {
		xValues = append(xValues, fmt.Sprintf("%02d/%02d/%d %02d:%02d", serie.Timestamp.Month(), serie.Timestamp.Day(), serie.Timestamp.Year(), serie.Timestamp.Hour(), serie.Timestamp.Minute()))
		for _, componentName := range config.ComponentList {
			found := false
			for _, c := range serie.Counts {
				if c.ComponentName == componentName {
					found = true
					perComponentSeries[componentName] = append(perComponentSeries[componentName], float64(c.Count))
					break
				}
			}
			if !found {
				perComponentSeries[componentName] = append(perComponentSeries[componentName], float64(0))
			}
		}
	}
	c := charts.NewLine()
	c.SetGlobalOptions(charts.TitleOpts{Title: fmt.Sprintf("%s", config.Release)})
	c.AddXAxis(xValues)
	for _, component := range config.ComponentList {
		c.AddYAxis(component, perComponentSeries[component])
	}
	c.Subtitle = fmt.Sprintf(" Total: %d", series[len(series)-1].Total)
	return c, nil
}
