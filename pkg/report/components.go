package report

import (
	"encoding/json"
	"time"

	"github.com/mfojtik/bugtrend/pkg/bugzilla"
)

type Component struct {
	Timestamp time.Time          `json:"timestamp"`
	Total     int                `json:"total"`
	Counts    []*ComponentStatus `json:"counts"`
}

type ComponentStatus struct {
	ComponentName string `json:"componentName"`
	Count         int    `json:"count"`
}

type ComponentList []Component

func (c *Component) Write() ([]byte, error) {
	return json.Marshal(&c)
}

func NewComponent(bugs []bugzilla.Bug) *Component {
	status := []*ComponentStatus{}
	for i := range bugs {
		found := false
		for c := range status {
			if status[c].ComponentName == bugs[i].Component[0] {
				status[c].Count++
				found = true
				break
			}
		}
		if !found {
			status = append(status, &ComponentStatus{
				ComponentName: bugs[i].Component[0],
				Count:         1,
			})
		}
	}
	return &Component{
		Timestamp: time.Now(),
		Total:     len(bugs),
		Counts:    status,
	}
}
