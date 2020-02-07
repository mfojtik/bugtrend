package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/mfojtik/bugtrend/pkg/analyze"
	"github.com/mfojtik/bugtrend/pkg/bugzilla"
	"github.com/mfojtik/bugtrend/pkg/report"
)

func writeBurnDownReport(release string, bugs []bugzilla.Bug) error {
	burnDownReport := report.NewBurnDown(bugs)
	burnDownBytes, err := burnDownReport.ToJson()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fmt.Sprintf("reports/%s/burndown_%d.json", release, burnDownReport.Timestamp.Unix()), burnDownBytes, os.ModePerm)
}

func writeBurnDownSummary(release string) error {
	summaryFile := ""
	err := filepath.Walk(path.Join("reports", release), func(p string, info os.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".json") {
			return nil
		}
		if info.Name() == "burndown.json" {
			return nil
		}
		content, err := ioutil.ReadFile(path.Join("reports", release, info.Name()))
		if len(summaryFile) == 0 {
			summaryFile = string(content)
			return nil
		}
		summaryFile = strings.Join([]string{summaryFile, string(content)}, ",\n")
		return nil
	})
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join("reports", release, "burndown.json"), []byte("["+summaryFile+"]"), os.ModePerm)
}

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

	client := bugzilla.NewClient(os.Getenv("BUGZILLA_API_KEY"))
	values := url.Values{
		"product":        []string{"OpenShift Container Platform"},
		"component":      []string{"kube-apiserver", "kube-controller-manager", "kube-scheduler", "oc", "service-ca", "Auth", "openshift-apiserver", "openshift-controller-manager", "Etcd", "Etcd Operator"},
		"version":        []string{"4.1.0", "4.2.0", "4.3.0", "4.4.0", "4.5.0", "4.6.0", "4.7.0", "4.9.0"},
		"target_release": []string{"---", release},
		"severity":       []string{"unspecified", "urgent", "high", "medium"},
		"priority":       []string{"unspecified", "urgent", "high", "medium"},
		"limit":          []string{"0"},
	}

	if err := os.MkdirAll(fmt.Sprintf("reports/%s", release), os.ModePerm); err != nil {
		log.Fatalf("Unable to make reports directory: %v", err)
	}

	go func() {
		if err := runDashboardHttpServer(release); err != nil {
			log.Fatal(err)
		}
	}()

	reportsDir := path.Join("reports", release)

	dashboardConfig := analyze.DashboardConfig{
		BurndownSeriesFile: path.Join(reportsDir, "burndown.json"),
		OutputFile:         path.Join(reportsDir, "index.html"),
		Release:            release,
	}

	// main loop
	for {
		result, err := client.Search(values)
		if err != nil {
			log.Printf("WARNING: Failed to query bugzilla: %v", err)
			time.Sleep(1 * time.Minute) // make retry faster when we fail to query BZ
			continue
		}

		if err := writeBurnDownReport(release, result.Bugs); err != nil {
			log.Printf("WARNING: Unable to write %s burndown report: %v", release, err)
		}

		if err := writeBurnDownSummary(release); err != nil {
			log.Printf("WARNING: Unable to write %s burndown summary: %v", release, err)
		}

		if err := analyze.WriteDashboard(dashboardConfig); err != nil {
			log.Printf("WARNING: Unable to write %s dashboard: %v", err)
		}

		log.Printf("Successfully processed %d bugs... ", len(result.Bugs))
		time.Sleep(4 * time.Hour)
	}
}
