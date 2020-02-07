package report

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/mfojtik/bugtrend/pkg/bugzilla"
)

type BugTypes struct {
	Timestamp time.Time      `json:"timestamp"`
	Total     int            `json:"total"`
	Counts    *BugTypeStatus `json:"counts"`
}

type BugTypeStatus struct {
	Flakes            int `json:"flakes"`
	Upgrade           int `json:"upgrade"`
	TestBlocker       int `json:"testBlocker"`
	AutomationBlocker int `json:"automationBlocker"`
	Other             int `json:"others"`
}

type BugTypeList []BugTypes

func (b *BugTypes) Write() ([]byte, error) {
	return json.Marshal(&b)
}

func isAutomationBlocker(bug bugzilla.Bug) bool {
	for _, k := range bug.Keywords {
		if k == "AutomationBlocker" {
			return true
		}
	}
	return false
}

func isTestBlocker(bug bugzilla.Bug) bool {
	for _, k := range bug.Keywords {
		if k == "TestBlocker" {
			return true
		}
	}
	return false
}

func isUpgrade(bug bugzilla.Bug) bool {
	switch {
	case strings.Contains(strings.ToLower(bug.Summary), "upgrad"): // upgrading, upgrade
		return true
	default:
		return false
	}
}

// isFlakeBug contains all keywords used to detect flake bugzilla's
func isFlakeBug(bug bugzilla.Bug) bool {
	switch {
	case strings.Contains(bug.Summary, "[ci]"):
		return true
	case strings.Contains(bug.Whiteboard, "buildcop"):
		return true
	case strings.Contains(bug.CfInternalWhiteboard, "buildcop"):
		return true
	case strings.Contains(bug.Summary, "buildcop"):
		return true
	case strings.Contains(strings.ToLower(bug.Summary), "conformance"):
		return true
	case strings.Contains(bug.Summary, "test/extended"):
		return true
	case strings.Contains(bug.Summary, "sig-"):
		return true
	case strings.Contains(bug.Summary, "e2e"):
		return true
	default:
		return false
	}
}

func NewBugTypes(bugs []bugzilla.Bug) *BugTypes {
	counts := &BugTypeStatus{}
	closedCount := 0
	for i := range bugs {
		switch {
		case isAutomationBlocker(bugs[i]):
			counts.AutomationBlocker++
		case isTestBlocker(bugs[i]):
			counts.TestBlocker++
		case isUpgrade(bugs[i]):
			counts.Upgrade++
		case isFlakeBug(bugs[i]):
			counts.Flakes++
		default:
			counts.Other++
		}
	}
	return &BugTypes{
		Timestamp: time.Now(),
		Total:     closedCount,
		Counts:    counts,
	}
}
