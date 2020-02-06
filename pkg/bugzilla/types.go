package bugzilla

import "time"

type SearchResult struct {
	Bugs []Bug `json:"bugs"`
}

type Bug struct {
	Priority         string    `json:"priority"`
	CfLastClosed     time.Time `json:"cf_last_closed"`
	AssignedToDetail struct {
		Email string `json:"email"`
	} `json:"assigned_to_detail"`
	Blocks         []int     `json:"blocks"`
	Creator        string    `json:"creator"`
	LastChangeTime time.Time `json:"last_change_time"`
	IsCcAccessible bool      `json:"is_cc_accessible"`
	Keywords       []string  `json:"keywords"`
	CreatorDetail  struct {
		Email string `json:"email"`
	} `json:"creator_detail"`
	CfQaWhiteboard          string    `json:"cf_qa_whiteboard"`
	Cc                      []string  `json:"cc"`
	URL                     string    `json:"url"`
	AssignedTo              string    `json:"assigned_to"`
	CfCustFacing            string    `json:"cf_cust_facing"`
	ID                      int       `json:"id"`
	Whiteboard              string    `json:"whiteboard"`
	CreationTime            time.Time `json:"creation_time"`
	QaContact               string    `json:"qa_contact"`
	DependsOn               []int     `json:"depends_on"`
	DupeOf                  int       `json:"dupe_of"`
	DocsContact             string    `json:"docs_contact"`
	CfTargetUpstreamVersion string    `json:"cf_target_upstream_version"`
	EstimatedTime           int       `json:"estimated_time"`
	RemainingTime           int       `json:"remaining_time"`
	CfPmScore               string    `json:"cf_pm_score"`
	Resolution              string    `json:"resolution"`
	Classification          string    `json:"classification"`
	CfDevelWhiteboard       string    `json:"cf_devel_whiteboard"`
	CfVerified              []string  `json:"cf_verified"`
	Alias                   []string  `json:"alias"`
	CfInternalWhiteboard    string    `json:"cf_internal_whiteboard"`
	CfDocType               string    `json:"cf_doc_type"`
	OpSys                   string    `json:"op_sys"`
	CfPgmInternal           string    `json:"cf_pgm_internal"`
	TargetRelease           []string  `json:"target_release"`
	Status                  string    `json:"status"`
	CfCloneOf               int       `json:"cf_clone_of"`
	Summary                 string    `json:"summary"`
	IsOpen                  bool      `json:"is_open"`
	Platform                string    `json:"platform"`
	Severity                string    `json:"severity"`
	CfEnvironment           string    `json:"cf_environment"`
	Version                 []string  `json:"version"`
	Component               []string  `json:"component"`
	ActualTime              int       `json:"actual_time"`
	CfFixedIn               string    `json:"cf_fixed_in"`
	IsCreatorAccessible     bool      `json:"is_creator_accessible"`
	IsConfirmed             bool      `json:"is_confirmed"`
	TargetMilestone         string    `json:"target_milestone"`
	Product                 string    `json:"product"`
	CfReleaseNotes          string    `json:"cf_release_notes"`
	CfBuildID               string    `json:"cf_build_id"`
}
