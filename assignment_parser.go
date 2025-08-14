package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	defaultGrading    = "Pass/fail"
	defaultAuto       = "Automatic"
	defaultManual     = "TA Approval"
	noApproval        = "No Approval"
	defaultIndividual = "Individually"
	defaultGroup      = "Group"
)

// AssignmentInfo contains fields present in the assignment.json files which accompany each lab.
type AssignmentInfo struct {
	Order         int
	Name          string
	ScriptFile    string
	Deadline      string // time.Time
	ShortDeadline string // only date
	Year          string // derived from deadline
	AutoApprove   bool
	IsGroupLab    bool
	ScoreLimit    int
	Reviewers     int
	Title         string
	HoursMin      int
	HoursMax      int
	Effort        string // replaces the HoursMin-HoursMax
	// defaults to $COURSE $NAME
	Subject string
	// defaults to defaultGrading
	Grading string
	// SubmissionType is interpreted based on the value of IsGroupLab ("Individually" or "Group")
	SubmissionType string
	ApproveType    string
	CourseOrg      string // this field is populated after parsing by the main script
}

func parseAssignment(filename string) (*AssignmentInfo, error) {
	jsonFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read json file: %w", err)
	}

	res := &AssignmentInfo{}
	err = json.Unmarshal(jsonFile, res)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	if res.Subject == "" {
		res.Subject = fmt.Sprintf("%s %s", strings.ToUpper(course()), name())
	}

	if res.ApproveType == "" {
		if res.AutoApprove {
			res.ApproveType = defaultAuto
		} else {
			res.ApproveType = defaultManual
		}
	}

	switch res.Grading {
	case "":
		res.Grading = defaultGrading
	case "No grading":
		res.ApproveType = noApproval
	}

	if res.IsGroupLab {
		res.SubmissionType = defaultGroup
	} else {
		res.SubmissionType = defaultIndividual
	}

	if res.Effort == "" {
		if res.HoursMin > 0 && res.HoursMax > 0 {
			res.Effort = fmt.Sprintf("%d-%d hours", res.HoursMin, res.HoursMax)
		}
	}

	// copied from qf.TimeLayout
	const layout = "2006-01-02T15:04:05"
	deadline, err := time.Parse(layout, res.Deadline)
	if err == nil {
		res.Deadline = fmt.Sprintf("%s %d, %d %02d:%02d",
			deadline.Month(), deadline.Day(), deadline.Year(), deadline.Hour(), deadline.Minute())
		res.ShortDeadline = fmt.Sprintf("%s %d", deadline.Month(), deadline.Day())
		res.Year = fmt.Sprintf("%d", deadline.Year())
		if res.Year != year() {
			fmt.Printf("Warning(%s): deadline (%s) does not match environment variable $YEAR=%s\n", res.Name, res.Deadline, year())
		}
	} else {
		fmt.Printf("Warning(%s): failed to parse deadline: %s\n", res.Name, err)
	}
	return res, nil
}
