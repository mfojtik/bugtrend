package dashboard

import (
	"fmt"
	"os"

	"github.com/go-echarts/go-echarts/charts"
)

type DashboardConfig struct {
	BurndownSeriesFile  string
	ClosedSeriesFile    string
	BugTypesSeriesFile  string
	OutputFile          string
	Release             string
	ComponentSeriesFile string
	ComponentList       []string
}

func WriteDashboard(config DashboardConfig) error {
	p := charts.NewPage()
	p.PageTitle = fmt.Sprintf("Group-B: %s Bugzilla Charts", config.Release)

	burndown, err := getBugBurndownChart(config)
	if err != nil {
		return err
	}
	burndown.Title = "Burndown"
	burndown.Subtitle = fmt.Sprintf("%s target release or unspecified.", config.Release) + burndown.Subtitle
	burndown.LegendOpts.Bottom = "0"
	p.Add(burndown)

	closed, err := getClosedChart(config)
	if err != nil {
		return err
	}
	closed.Title = "Closed"
	closed.Subtitle = fmt.Sprintf("%s closed bugs resolutions.", config.Release) + closed.Subtitle
	closed.LegendOpts.Bottom = "0"
	p.Add(closed)

	types, err := getBugTypesChart(config)
	if err != nil {
		return err
	}
	types.Title = "Types"
	types.Subtitle = fmt.Sprintf("%s bugs counts by type.", config.Release)
	types.LegendOpts.Bottom = "0"
	p.Add(types)

	components, err := getComponentsChart(config)
	if err != nil {
		return err
	}
	components.Title = "Components"
	components.Subtitle = fmt.Sprintf("%s bugs counts by component.", config.Release)
	components.LegendOpts.Bottom = "0"
	p.Add(components)

	f, err := os.Create(config.OutputFile)
	if err != nil {
		return err
	}

	return p.Render(f)
}
