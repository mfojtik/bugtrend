package report

import (
	"encoding/json"
	"time"

	"github.com/mfojtik/bugtrend/pkg/bugzilla"
)

type Closed struct {
	Timestamp time.Time           `json:"timestamp"`
	Total     int                 `json:"total"`
	Counts    []*ResolutionStatus `json:"counts"`
}

type ResolutionStatus struct {
	Resolution string `json:"resolution"`
	Count      int    `json:"count"`
}

type ClosedList []Closed

func (b *Closed) Write() ([]byte, error) {
	return json.Marshal(&b)
}

func NewClosed(bugs []bugzilla.Bug) *Closed {
	counts := []*ResolutionStatus{}
	closedCount := 0
	for i := range bugs {
		if bugs[i].Status != "CLOSED" {
			continue
		}
		closedCount++
		found := false
		for c := range counts {
			if counts[c].Resolution == bugs[i].Resolution {
				counts[c].Count++
				found = true
				break
			}
		}
		if !found {
			counts = append(counts, &ResolutionStatus{
				Resolution: bugs[i].Resolution,
				Count:      1,
			})
		}
	}
	return &Closed{
		Timestamp: time.Now(),
		Total:     closedCount,
		Counts:    counts,
	}
}
