package analyze

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-echarts/go-echarts/charts"

	"github.com/mfojtik/bugtrend/pkg/report"
)

type reportList []report.BurnDownReport

func bugBurnDownChart(sourceJsonPath, release string) (*charts.Line, error) {
	reportFile, err := ioutil.ReadFile(sourceJsonPath)
	if err != nil {
		return nil, err
	}

	var series reportList
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
	xValues := []string{}

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
	c.SetGlobalOptions(charts.TitleOpts{Title: fmt.Sprintf("%s", release)})
	c.AddXAxis(xValues)
	for _, state := range states {
		c.AddYAxis(state, perStateSeries[state])
	}
	return c, nil
}

type DashboardConfig struct {
	BurndownSeriesFile string
	OutputFile         string
	Release            string
}

func WriteDashboard(config DashboardConfig) error {
	p := charts.NewPage()

	if burndown, err := bugBurnDownChart(config.BurndownSeriesFile, config.Release); err != nil {
		return err
	} else {
		p.Add(burndown)
	}

	f, err := os.Create(config.OutputFile)
	if err != nil {
		return err
	}
	return p.Render(f)
}
