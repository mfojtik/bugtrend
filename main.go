package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/mfojtik/bugtrend/pkg/dashboard"
)

func runDashboardHttpServer(release string) error {
	releaseDir := path.Join("reports", release)
	http.Handle("/", http.FileServer(http.Dir(releaseDir)))
	log.Printf("Serving %q on HTTP: localhost:%s\n", releaseDir, "8080")
	return http.ListenAndServe(":8080", nil)
}

func main() {
	release := os.Getenv("OPENSHIFT_TARGET_RELEASE")
	if len(release) == 0 {
		panic("OPENSHIFT_TARGET_RELEASE environment variable must be set (ex. '4.4.0')")
	}
	apiKey := os.Getenv("BUGZILLA_API_KEY")
	if len(apiKey) == 0 {
		panic("BUGZILLA_API_KEY environment variable must be set (https://bugzilla.redhat.com/userprefs.cgi?tab=apikey)")
	}

	/*
		client := bugzilla.NewClient(os.Getenv("BUGZILLA_API_KEY"))

		// this is the main query that returns all group-b component bugs.
		values := url.Values{
			"product":        []string{"OpenShift Container Platform"},
			"component":      []string{"kube-apiserver", "kube-controller-manager", "kube-scheduler", "oc", "service-ca", "Auth", "openshift-apiserver", "openshift-controller-manager", "Etcd", "Etcd Operator"},
			"version":        []string{"4.1.0", "4.2.0", "4.3.0", "4.4.0", "4.5.0", "4.6.0", "4.7.0", "4.9.0"},
			"target_release": []string{"---", release},
			"severity":       []string{"unspecified", "urgent", "high", "medium"},
			"priority":       []string{"unspecified", "urgent", "high", "medium"},
			"limit":          []string{"0"},
		}
	*/

	if err := os.MkdirAll(fmt.Sprintf("reports/%s", release), os.ModePerm); err != nil {
		log.Fatalf("Unable to make reports directory: %v", err)
	}

	go func() {
		if err := runDashboardHttpServer(release); err != nil {
			log.Fatal(err)
		}
	}()

	reportsDir := path.Join("reports", release)

	dashboardConfig := dashboard.DashboardConfig{
		BurndownSeriesFile: path.Join(reportsDir, "burndown.json"),
		ClosedSeriesFile:   path.Join(reportsDir, "closed.json"),
		BugTypesSeriesFile: path.Join(reportsDir, "types.json"),
		OutputFile:         path.Join(reportsDir, "index.html"),
		Release:            release,
	}

	// main loop
	for {
		/*
			result, err := client.Search(values)
			if err != nil {
				log.Printf("WARNING: Failed to query bugzilla: %v", err)
				time.Sleep(1 * time.Minute) // make retry faster when we fail to query BZ
				continue
			}

			bugBurndownReport := report.NewBurnDown(result.Bugs)
			bugBurndown := report.NewTimeSerieWriter("burndown", release, time.Now())
			if err := bugBurndown.WriteTimeSerie(bugBurndownReport); err != nil {
				log.Printf("WARNING: Unable to write %s burndown report: %v", release, err)
			}

			closedReport := report.NewClosed(result.Bugs)
			closed := report.NewTimeSerieWriter("closed", release, time.Now())
			if err := closed.WriteTimeSerie(closedReport); err != nil {
				log.Printf("WARNING: Unable to write %s closed report: %v", release, err)
			}

			typesReport := report.NewBugTypes(result.Bugs)
			types := report.NewTimeSerieWriter("types", release, time.Now())
			if err := types.WriteTimeSerie(typesReport); err != nil {
				log.Printf("WARNING: Unable to write %s types report: %v", release, err)
			}

			log.Printf("Successfully processed %d bugs... ", len(result.Bugs))
		*/

		if err := dashboard.WriteDashboard(dashboardConfig); err != nil {
			log.Printf("WARNING: Unable to write %s dashboard: %v", err)
		}
		time.Sleep(4 * time.Hour)
	}
}
