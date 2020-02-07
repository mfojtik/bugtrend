package dashboard

import (
	"fmt"
	"os"

	"github.com/go-echarts/go-echarts/charts"
)

type DashboardConfig struct {
	BurndownSeriesFile string
	ClosedSeriesFile   string
	BugTypesSeriesFile string
	OutputFile         string
	Release            string
}

func WriteDashboard(config DashboardConfig) error {
	p := charts.NewPage()
	p.PageTitle = "Group-B: Bugzilla Charts"

	burndown, err := getBugBurndownChart(config)
	if err != nil {
		return err
	}
	burndown.Title = "Bug Burndown"
	burndown.Subtitle = fmt.Sprintf("Bugs with Target Release %s or Unspecified.", config.Release) + burndown.Subtitle
	burndown.LegendOpts.Bottom = "0"
	p.Add(burndown)

	closed, err := getClosedChart(config)
	if err != nil {
		return err
	}
	closed.Title = "Closed Breakout"
	closed.Subtitle = fmt.Sprintf("Closed %s bugs resolution breakout", config.Release)
	closed.LegendOpts.Bottom = "0"
	p.Add(closed)

	types, err := getBugTypesChart(config)
	if err != nil {
		return err
	}
	types.Title = "Bug By Type"
	types.Subtitle = fmt.Sprintf("%s bugs counts by type", config.Release)
	types.LegendOpts.Bottom = "0"
	p.Add(types)

	f, err := os.Create(config.OutputFile)
	if err != nil {
		return err
	}

	return p.Render(f)
}
