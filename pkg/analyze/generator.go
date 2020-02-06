package analyze

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-echarts/go-echarts/charts"

	"github.com/mfojtik/bugtrend/pkg/report"
)

type reportList []report.BurnDownReport

func WriteIndex(path string, summaryPath string, release string) {
	reportFile, err := ioutil.ReadFile(summaryPath)
	if err != nil {
		log.Printf("failed to open file: %v", err)
		return
	}

	var series reportList
	if err := json.Unmarshal(reportFile, &series); err != nil {
		log.Printf("failed to parse JSON: %v", err)
		return
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

	p := charts.NewPage()

	bar := charts.NewLine()
	bar.SetGlobalOptions(charts.TitleOpts{Title: fmt.Sprintf("%s bugs", release)})
	bar.AddXAxis(xValues).
		AddYAxis("NEW", perStateSeries["NEW"]).
		AddYAxis("ASSIGNED", perStateSeries["ASSIGNED"]).
		AddYAxis("POST", perStateSeries["POST"]).
		AddYAxis("MODIFIED", perStateSeries["MODIFIED"]).
		AddYAxis("ON_QA", perStateSeries["ON_QA"]).
		AddYAxis("VERIFIED", perStateSeries["VERIFIED"]).
		AddYAxis("CLOSED", perStateSeries["CLOSED"])

	p.Add(bar)

	f, err := os.Create(path)
	if err != nil {
		log.Printf("failed to create: %v", err)
		return
	}
	if err := p.Render(f); err != nil {
		log.Printf("render failed: %v", err)
	}
}
