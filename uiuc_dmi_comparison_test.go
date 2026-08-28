// Copyright 2026 Board of Trustees of the University of Illinois.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCompareUIUCAttributesProvostRowCreatesUnchangedCollegeAndDepartment(t *testing.T) {
	analysis := &uiucDMIAnalysis{
		Colleges: []uiucDMICollege{{
			ID:                "illinois.edu:LM",
			CampusID:          "illinois.edu",
			UnitID:            "LM",
			Name:              "Provost Academic Programs",
			ParentSourceLines: []int{311},
		}},
		Departments: []uiucDMIDepartment{{
			ID:                "illinois.edu:LM:290",
			CampusID:          "illinois.edu",
			UnitID:            "LM",
			DepartmentID:      "290",
			Name:              "Provost Courses",
			FullName:          "Provost Courses",
			ParentCollegeID:   "illinois.edu:LM",
			ParentCollegeName: "Provost Academic Programs",
			SourceLine:        311,
		}},
	}
	current := &currentUIUCAttributes{
		Colleges: []currentUIUCCollege{{Value: "Provost Academic Programs"}},
		Departments: []currentUIUCDepartment{{
			Group:             "Provost Academic Programs",
			Label:             "Provost Courses",
			RequiredCollege:   "Provost Academic Programs",
			currentIdentifier: "Provost Academic Programs\x00Provost Courses",
		}},
	}

	comparison := compareUIUCAttributes(analysis, current)
	if got := comparison.Colleges[0].Status; got != uiucStatusUnchanged {
		t.Fatalf("college status = %q, want %q", got, uiucStatusUnchanged)
	}
	if got := comparison.Departments[0].Status; got != uiucStatusUnchanged {
		t.Fatalf("department status = %q, want %q", got, uiucStatusUnchanged)
	}
	if len(comparison.CurrentOnlyColleges) != 0 || len(comparison.CurrentOnlyDepartments) != 0 {
		t.Fatalf("expected every current item to be matched: %+v", comparison)
	}
}

func TestCompareUIUCAttributesInfersCollegeRenameFromDepartmentOverlap(t *testing.T) {
	analysis := &uiucDMIAnalysis{
		Colleges: []uiucDMICollege{{
			ID:                "illinois.edu:KL",
			Name:              "New College Name",
			DirectSourceLines: []int{10},
		}},
		Departments: []uiucDMIDepartment{
			{ID: "illinois.edu:KL:001", Name: "Department A", ParentCollegeID: "illinois.edu:KL", ParentCollegeName: "New College Name"},
			{ID: "illinois.edu:KL:002", Name: "Department B", ParentCollegeID: "illinois.edu:KL", ParentCollegeName: "New College Name"},
		},
	}
	current := &currentUIUCAttributes{
		Colleges: []currentUIUCCollege{{Value: "Old College Name"}},
		Departments: []currentUIUCDepartment{
			{Group: "Old College Name", Label: "Department A", RequiredCollege: "Old College Name", currentIdentifier: "Old College Name\x00Department A"},
			{Group: "Old College Name", Label: "Department B", RequiredCollege: "Old College Name", currentIdentifier: "Old College Name\x00Department B"},
		},
	}

	comparison := compareUIUCAttributes(analysis, current)
	college := comparison.Colleges[0]
	if college.Status != uiucStatusChangedInferred {
		t.Fatalf("college status = %q, want %q", college.Status, uiucStatusChangedInferred)
	}
	if college.CurrentValue != "Old College Name" {
		t.Fatalf("current college = %q, want Old College Name", college.CurrentValue)
	}
	if !strings.Contains(college.Evidence, "2 exact department-name match(es)") {
		t.Fatalf("evidence does not explain overlap: %q", college.Evidence)
	}
	for _, department := range comparison.Departments {
		if department.Status != uiucStatusUnchanged {
			t.Fatalf("department %q status = %q, want %q", department.Source.Name, department.Status, uiucStatusUnchanged)
		}
	}
}

func TestCompareUIUCAttributesDoesNotClaimUnsupportedRename(t *testing.T) {
	analysis := &uiucDMIAnalysis{Colleges: []uiucDMICollege{{
		ID:                "illinois.edu:ZZ",
		Name:              "Unmatched Live College",
		DirectSourceLines: []int{10},
	}}}
	current := &currentUIUCAttributes{
		Colleges: []currentUIUCCollege{{Value: "Different Current College"}},
		Departments: []currentUIUCDepartment{{
			Group:             "Different Current College",
			Label:             "Unrelated Department",
			RequiredCollege:   "Different Current College",
			currentIdentifier: "Different Current College\x00Unrelated Department",
		}},
	}

	comparison := compareUIUCAttributes(analysis, current)
	if got := comparison.Colleges[0].Status; got != uiucStatusNew {
		t.Fatalf("college status = %q, want conservative %q", got, uiucStatusNew)
	}
	if len(comparison.CurrentOnlyColleges) != 1 {
		t.Fatalf("current-only colleges = %d, want 1", len(comparison.CurrentOnlyColleges))
	}
}

func TestReportPrintsHumanNameBeforeMetadata(t *testing.T) {
	analysis := &uiucDMIAnalysis{}
	current := &currentUIUCAttributes{Path: "contentdb.content_items.json"}
	comparison := &uiucAttributeComparison{
		Colleges: []uiucCollegeComparison{{
			Source: uiucDMICollege{
				ID:                "illinois.edu:LM",
				Name:              "Provost Academic Programs",
				ParentSourceLines: []int{311},
			},
			Status:            uiucStatusUnchanged,
			CurrentValue:      "Provost Academic Programs",
			CurrentDisplay:    "Provost Academic Programs",
			Evidence:          "exact source name match",
			SourceExplanation: "parent college required by line 311",
		}},
	}

	var output bytes.Buffer
	if err := printUIUCAttributeComparison(&output, analysis, current, comparison); err != nil {
		t.Fatalf("print report: %v", err)
	}
	want := "\nProvost Academic Programs\n  COMPARISON\n    status:   UNCHANGED\n"
	if !strings.Contains(output.String(), want) {
		t.Fatalf("report does not lead with the human name; missing %q in:\n%s", want, output.String())
	}
}
