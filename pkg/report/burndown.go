package report

import (
	"encoding/json"
	"time"

	"github.com/mfojtik/bugtrend/pkg/bugzilla"
)

type BurndownStatus struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type Burndown struct {
	Timestamp time.Time         `json:"timestamp"`
	Total     int               `json:"total"`
	Counts    []*BurndownStatus `json:"counts"`
}

type BurndownList []Burndown

func (b *Burndown) Write() ([]byte, error) {
	return json.Marshal(&b)
}

func NewBurnDown(bugs []bugzilla.Bug) *Burndown {
	counts := []*BurndownStatus{}
	for i := range bugs {
		found := false
		for c := range counts {
			if counts[c].Status == bugs[i].Status {
				counts[c].Count++
				found = true
				break
			}
		}
		if !found {
			counts = append(counts, &BurndownStatus{
				Status: bugs[i].Status,
				Count:  1,
			})
		}
	}
	return &Burndown{
		Timestamp: time.Now(),
		Total:     len(bugs),
		Counts:    counts,
	}
}
