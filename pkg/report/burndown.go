package report

import (
	"encoding/json"
	"log"
	"time"

	"github.com/mfojtik/bugtrend/pkg/bugzilla"
)

type StatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type BurnDownReport struct {
	Timestamp time.Time      `json:"timestamp"`
	Total     int            `json:"total"`
	Counts    []*StatusCount `json:"counts"`
}

func (b *BurnDownReport) ToJson() ([]byte, error) {
	return json.Marshal(&b)
}

func NewBurnDown(bugs []bugzilla.Bug) *BurnDownReport {
	counts := []*StatusCount{}
	for i := range bugs {
		found := false
		log.Printf("%q\n", bugs[i].Status)
		for c := range counts {
			if counts[c].Status == bugs[i].Status {
				counts[c].Count++
				found = true
				break
			}
		}
		if !found {
			counts = append(counts, &StatusCount{
				Status: bugs[i].Status,
				Count:  1,
			})
		}
	}
	return &BurnDownReport{
		Timestamp: time.Now(),
		Total:     len(bugs),
		Counts:    counts,
	}
}
