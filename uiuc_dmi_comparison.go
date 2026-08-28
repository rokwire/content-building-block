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
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

const (
	defaultCurrentUIUCAttributesPath = "contentdb.content_items.json"
	currentUIUCAttributesPathEnv     = "UIUC_CURRENT_ATTRIBUTES_FILE"

	uiucStatusUnchanged       = "UNCHANGED"
	uiucStatusChangedInferred = "CHANGED (INFERRED)"
	uiucStatusNew             = "NEW"
	uiucStatusReview          = "REVIEW"
)

type currentUIUCCollege struct {
	Value string
	Label string
}

func (college currentUIUCCollege) displayName() string {
	if college.Label != "" {
		return college.Label
	}
	return college.Value
}

func (college currentUIUCCollege) terms() []string {
	return uniqueNonEmptyStrings(college.Value, college.Label)
}

type currentUIUCDepartment struct {
	Group             string
	Label             string
	RequiredCollege   string
	currentIdentifier string
}

type currentUIUCAttributes struct {
	Path        string
	Colleges    []currentUIUCCollege
	Departments []currentUIUCDepartment
}

type uiucAttributeComparison struct {
	Colleges               []uiucCollegeComparison
	Departments            []uiucDepartmentComparison
	CurrentOnlyColleges    []currentUIUCCollege
	CurrentOnlyDepartments []currentUIUCDepartment
}

type uiucCollegeComparison struct {
	Source            uiucDMICollege
	Status            string
	CurrentValue      string
	CurrentDisplay    string
	Evidence          string
	SourceExplanation string
}

type uiucDepartmentComparison struct {
	Source            uiucDMIDepartment
	Status            string
	CurrentLabel      string
	CurrentGroup      string
	Evidence          string
	SourceExplanation string
}

type rawCurrentContentItem struct {
	Category string `json:"category"`
	Data     struct {
		Attributes []struct {
			ID     string            `json:"id"`
			Values []json.RawMessage `json:"values"`
		} `json:"attributes"`
	} `json:"data"`
}

func currentUIUCAttributesFile() string {
	if path := strings.TrimSpace(os.Getenv(currentUIUCAttributesPathEnv)); path != "" {
		return path
	}
	return defaultCurrentUIUCAttributesPath
}

func loadCurrentUIUCAttributes(path string) (*currentUIUCAttributes, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open current attributes snapshot %q: %w", path, err)
	}
	defer file.Close()

	var items []rawCurrentContentItem
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&items); err != nil {
		return nil, fmt.Errorf("decode current attributes snapshot %q: %w", path, err)
	}

	result := &currentUIUCAttributes{Path: path}
	foundAttributes := false
	for _, item := range items {
		if item.Category != "attributes" {
			continue
		}
		foundAttributes = true
		for _, attribute := range item.Data.Attributes {
			switch attribute.ID {
			case "college":
				for index, value := range attribute.Values {
					college, err := decodeCurrentUIUCCollege(value)
					if err != nil {
						return nil, fmt.Errorf("decode college value %d: %w", index, err)
					}
					result.Colleges = append(result.Colleges, college)
				}
			case "department":
				for index, value := range attribute.Values {
					department, err := decodeCurrentUIUCDepartment(value)
					if err != nil {
						return nil, fmt.Errorf("decode department value %d: %w", index, err)
					}
					result.Departments = append(result.Departments, department)
				}
			}
		}
	}
	if !foundAttributes {
		return nil, fmt.Errorf("current attributes snapshot %q has no content item with category=attributes", path)
	}
	if len(result.Colleges) == 0 || len(result.Departments) == 0 {
		return nil, fmt.Errorf(
			"current attributes snapshot %q must contain non-empty college and department attributes",
			path,
		)
	}

	sort.Slice(result.Colleges, func(i, j int) bool {
		return result.Colleges[i].displayName() < result.Colleges[j].displayName()
	})
	sort.Slice(result.Departments, func(i, j int) bool {
		if result.Departments[i].Group != result.Departments[j].Group {
			return result.Departments[i].Group < result.Departments[j].Group
		}
		return result.Departments[i].Label < result.Departments[j].Label
	})
	return result, nil
}

func decodeCurrentUIUCCollege(raw json.RawMessage) (currentUIUCCollege, error) {
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		text = strings.TrimSpace(text)
		if text == "" {
			return currentUIUCCollege{}, fmt.Errorf("empty string")
		}
		return currentUIUCCollege{Value: text}, nil
	}

	var object struct {
		Label string `json:"label"`
		Value string `json:"value"`
	}
	if err := json.Unmarshal(raw, &object); err != nil {
		return currentUIUCCollege{}, err
	}
	object.Label = strings.TrimSpace(object.Label)
	object.Value = strings.TrimSpace(object.Value)
	if object.Label == "" || object.Value == "" {
		return currentUIUCCollege{}, fmt.Errorf("object must have non-empty label and value")
	}
	return currentUIUCCollege{Value: object.Value, Label: object.Label}, nil
}

func decodeCurrentUIUCDepartment(raw json.RawMessage) (currentUIUCDepartment, error) {
	var object struct {
		Group        string `json:"group"`
		Label        string `json:"label"`
		Requirements struct {
			College string `json:"college"`
		} `json:"requirements"`
	}
	if err := json.Unmarshal(raw, &object); err != nil {
		return currentUIUCDepartment{}, err
	}
	department := currentUIUCDepartment{
		Group:           strings.TrimSpace(object.Group),
		Label:           strings.TrimSpace(object.Label),
		RequiredCollege: strings.TrimSpace(object.Requirements.College),
	}
	if department.Group == "" || department.Label == "" || department.RequiredCollege == "" {
		return currentUIUCDepartment{}, fmt.Errorf("object must have non-empty group, label and requirements.college")
	}
	department.currentIdentifier = department.Group + "\x00" + department.Label
	return department, nil
}

func compareUIUCAttributes(analysis *uiucDMIAnalysis, current *currentUIUCAttributes) *uiucAttributeComparison {
	comparison := &uiucAttributeComparison{}
	matchedColleges := make(map[string]bool)
	matchedDepartments := make(map[string]bool)

	collegeByTerm := make(map[string]currentUIUCCollege)
	for _, college := range current.Colleges {
		for _, term := range college.terms() {
			collegeByTerm[term] = college
		}
	}

	currentDepartmentsByGroup := make(map[string][]currentUIUCDepartment)
	currentDepartmentsByLabel := make(map[string][]currentUIUCDepartment)
	for _, department := range current.Departments {
		currentDepartmentsByGroup[department.Group] = append(currentDepartmentsByGroup[department.Group], department)
		currentDepartmentsByLabel[department.Label] = append(currentDepartmentsByLabel[department.Label], department)
	}

	sourceDepartmentsByCollege := make(map[string][]uiucDMIDepartment)
	for _, department := range analysis.Departments {
		sourceDepartmentsByCollege[department.ParentCollegeID] = append(
			sourceDepartmentsByCollege[department.ParentCollegeID], department,
		)
	}

	collegeComparisonByID := make(map[string]uiucCollegeComparison)
	for _, college := range analysis.Colleges {
		item := compareUIUCCollege(
			college,
			sourceDepartmentsByCollege[college.ID],
			collegeByTerm,
			currentDepartmentsByGroup,
		)
		comparison.Colleges = append(comparison.Colleges, item)
		collegeComparisonByID[college.ID] = item
		if item.CurrentValue != "" {
			matchedColleges[item.CurrentValue] = true
		}
	}

	for _, department := range analysis.Departments {
		item := compareUIUCDepartment(
			department,
			collegeComparisonByID[department.ParentCollegeID],
			currentDepartmentsByLabel,
		)
		comparison.Departments = append(comparison.Departments, item)
		if item.CurrentLabel != "" && item.CurrentGroup != "" {
			matchedDepartments[item.CurrentGroup+"\x00"+item.CurrentLabel] = true
		}
	}

	for _, college := range current.Colleges {
		if !matchedColleges[college.Value] {
			comparison.CurrentOnlyColleges = append(comparison.CurrentOnlyColleges, college)
		}
	}
	for _, department := range current.Departments {
		if !matchedDepartments[department.currentIdentifier] {
			comparison.CurrentOnlyDepartments = append(comparison.CurrentOnlyDepartments, department)
		}
	}
	return comparison
}

func compareUIUCCollege(
	source uiucDMICollege,
	sourceDepartments []uiucDMIDepartment,
	currentByTerm map[string]currentUIUCCollege,
	currentDepartmentsByGroup map[string][]currentUIUCDepartment,
) uiucCollegeComparison {
	result := uiucCollegeComparison{
		Source:            source,
		SourceExplanation: collegeSourceExplanation(source),
	}
	if current, ok := currentByTerm[source.Name]; ok {
		result.Status = uiucStatusUnchanged
		result.CurrentValue = current.Value
		result.CurrentDisplay = current.displayName()
		result.Evidence = fmt.Sprintf("exact source name match in current college values: %q", source.Name)
		if current.Label != "" && current.Label != current.Value {
			result.Evidence += fmt.Sprintf(" (current label=%q, value=%q)", current.Label, current.Value)
		}
		return result
	}

	sourceNames := make(map[string]bool)
	for _, department := range sourceDepartments {
		sourceNames[department.Name] = true
		if department.FullName != "" {
			sourceNames[department.FullName] = true
		}
	}
	type candidate struct {
		Group string
		Score int
	}
	candidates := make([]candidate, 0, len(currentDepartmentsByGroup))
	for group, departments := range currentDepartmentsByGroup {
		score := 0
		for _, department := range departments {
			if sourceNames[department.Label] {
				score++
			}
		}
		if score > 0 {
			candidates = append(candidates, candidate{Group: group, Score: score})
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Score != candidates[j].Score {
			return candidates[i].Score > candidates[j].Score
		}
		return candidates[i].Group < candidates[j].Group
	})

	if len(candidates) == 0 {
		result.Status = uiucStatusNew
		result.Evidence = "no exact college-name match and no exact department-name overlap with a current college"
		return result
	}
	if len(candidates) > 1 && candidates[0].Score == candidates[1].Score {
		result.Status = uiucStatusReview
		result.Evidence = fmt.Sprintf(
			"no exact college-name match; department overlap is tied between %q and %q (%d match(es) each)",
			candidates[0].Group, candidates[1].Group, candidates[0].Score,
		)
		return result
	}

	best := candidates[0]
	result.Status = uiucStatusChangedInferred
	result.CurrentValue = best.Group
	result.CurrentDisplay = best.Group
	if current, ok := currentByTerm[best.Group]; ok {
		result.CurrentValue = current.Value
		result.CurrentDisplay = current.displayName()
	}
	result.Evidence = fmt.Sprintf(
		"probable rename from current %q: %d exact department-name match(es) and no competing group with the same score",
		best.Group, best.Score,
	)
	return result
}

func compareUIUCDepartment(
	source uiucDMIDepartment,
	parent uiucCollegeComparison,
	currentByLabel map[string][]currentUIUCDepartment,
) uiucDepartmentComparison {
	result := uiucDepartmentComparison{
		Source: source,
		SourceExplanation: fmt.Sprintf(
			"department because Campus=1, Active=Y, Dean=F and department_id=%s",
			source.DepartmentID,
		),
	}

	if matches := currentByLabel[source.Name]; len(matches) > 0 {
		return resolveDepartmentMatch(result, matches, source.Name, false, parent)
	}
	if source.FullName != "" && source.FullName != source.Name {
		if matches := currentByLabel[source.FullName]; len(matches) > 0 {
			return resolveDepartmentMatch(result, matches, source.FullName, true, parent)
		}
	}

	result.Status = uiucStatusNew
	result.Evidence = "neither the source short name nor full name exists in current department labels"
	return result
}

func resolveDepartmentMatch(
	result uiucDepartmentComparison,
	matches []currentUIUCDepartment,
	matchedName string,
	matchedFullName bool,
	parent uiucCollegeComparison,
) uiucDepartmentComparison {
	if parent.CurrentValue != "" {
		for _, match := range matches {
			if match.Group == parent.CurrentValue || match.RequiredCollege == parent.CurrentValue {
				result.CurrentLabel = match.Label
				result.CurrentGroup = match.Group
				if matchedFullName {
					result.Status = uiucStatusChangedInferred
					result.Evidence = fmt.Sprintf(
						"current label %q matches source full name; source now selects short name %q",
						matchedName, result.Source.Name,
					)
				} else {
					result.Status = uiucStatusUnchanged
					result.Evidence = fmt.Sprintf(
						"exact department-name match under current college %q",
						match.Group,
					)
				}
				return result
			}
		}
	}

	if len(matches) == 1 {
		match := matches[0]
		result.CurrentLabel = match.Label
		result.CurrentGroup = match.Group
		if parent.CurrentValue == "" {
			result.Status = uiucStatusReview
			result.Evidence = fmt.Sprintf(
				"department name exists under current college %q, but source college %q has no reliable current mapping",
				match.Group, result.Source.ParentCollegeName,
			)
		} else {
			result.Status = uiucStatusChangedInferred
			result.Evidence = fmt.Sprintf(
				"same department name exists under current college %q; source places it under %q",
				match.Group, result.Source.ParentCollegeName,
			)
		}
		if matchedFullName {
			result.Evidence += fmt.Sprintf("; current label matches source full name, while output uses %q", result.Source.Name)
		}
		return result
	}

	result.Status = uiucStatusReview
	result.Evidence = fmt.Sprintf(
		"department name %q appears in %d current groups and the source parent cannot resolve it uniquely",
		matchedName, len(matches),
	)
	return result
}

func collegeSourceExplanation(source uiucDMICollege) string {
	if len(source.DirectSourceLines) > 0 {
		return fmt.Sprintf(
			"college/admin because Dean=T on source line(s) %s",
			formatSourceLines(source.DirectSourceLines),
		)
	}
	return fmt.Sprintf(
		"parent college required by Dean=F department line(s) %s; no direct Dean=T row",
		formatSourceLines(source.ParentSourceLines),
	)
}

func printUIUCAttributeComparison(
	output io.Writer,
	analysis *uiucDMIAnalysis,
	current *currentUIUCAttributes,
	comparison *uiucAttributeComparison,
) error {
	writer := bufio.NewWriter(output)

	fmt.Fprintln(writer, "================================================================================")
	fmt.Fprintln(writer, "UIUC ATTRIBUTES: LIVE SOURCE COMPARED WITH CURRENT CONTENT ITEM")
	fmt.Fprintln(writer, "================================================================================")
	fmt.Fprintf(writer, "Live source:      %s\n", analysis.SourceURL)
	fmt.Fprintf(writer, "Fetched at:       %s UTC\n", analysis.FetchedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(writer, "Current snapshot: %s\n", current.Path)
	fmt.Fprintf(writer, "\nLive result:      %d colleges, %d departments\n", len(analysis.Colleges), len(analysis.Departments))
	fmt.Fprintf(writer, "Current database:  %d colleges, %d departments\n", len(current.Colleges), len(current.Departments))
	printStatusSummary(writer, comparison)

	printCollegeComparisons(writer, comparison.Colleges)
	printDepartmentComparisons(writer, comparison.Departments)
	printCurrentOnlyItems(writer, comparison)

	fmt.Fprintln(writer, "\n================================================================================")
	fmt.Fprintln(writer, "END OF REPORT")
	fmt.Fprintln(writer, "================================================================================")
	return writer.Flush()
}

func printStatusSummary(writer *bufio.Writer, comparison *uiucAttributeComparison) {
	collegeCounts := countCollegeStatuses(comparison.Colleges)
	departmentCounts := countDepartmentStatuses(comparison.Departments)
	fmt.Fprintln(writer, "\nComparison summary:")
	for _, status := range []string{uiucStatusUnchanged, uiucStatusChangedInferred, uiucStatusNew, uiucStatusReview} {
		fmt.Fprintf(writer, "  %-20s colleges=%3d  departments=%3d\n", status, collegeCounts[status], departmentCounts[status])
	}
	fmt.Fprintf(writer, "  %-20s colleges=%3d  departments=%3d\n",
		"CURRENT ONLY", len(comparison.CurrentOnlyColleges), len(comparison.CurrentOnlyDepartments))
	fmt.Fprintln(writer, "\nLegend: the current JSON contains names but no unit_id or department_id.")
	fmt.Fprintln(writer, "        CHANGED (INFERRED) is supported by name/group overlap, but is not an ID match.")
	fmt.Fprintln(writer, "        NEW means no reliable current match was found; it can still be an unrecognized rename.")
	fmt.Fprintln(writer, "        CURRENT ONLY does not automatically mean removed.")
}

func printCollegeComparisons(writer *bufio.Writer, items []uiucCollegeComparison) {
	fmt.Fprintln(writer, "\n================================================================================")
	fmt.Fprintf(writer, "LIVE COLLEGES (%d)\n", len(items))
	fmt.Fprintln(writer, "================================================================================")
	for _, item := range items {
		fmt.Fprintf(writer, "\n%s\n", item.Source.Name)
		fmt.Fprintln(writer, "  COMPARISON")
		fmt.Fprintf(writer, "    status:   %s\n", item.Status)
		if item.CurrentDisplay != "" {
			fmt.Fprintf(writer, "    current:  %s", item.CurrentDisplay)
			if item.CurrentValue != "" && item.CurrentValue != item.CurrentDisplay {
				fmt.Fprintf(writer, " (stored value: %s)", item.CurrentValue)
			}
			fmt.Fprintln(writer)
		}
		fmt.Fprintf(writer, "    evidence: %s\n", item.Evidence)
		fmt.Fprintln(writer, "  SOURCE")
		fmt.Fprintf(writer, "    id:       %s\n", item.Source.ID)
		fmt.Fprintf(writer, "    rule:     %s\n", item.SourceExplanation)
	}
}

func printDepartmentComparisons(writer *bufio.Writer, items []uiucDepartmentComparison) {
	fmt.Fprintln(writer, "\n================================================================================")
	fmt.Fprintf(writer, "LIVE DEPARTMENTS (%d)\n", len(items))
	fmt.Fprintln(writer, "================================================================================")
	for _, item := range items {
		fmt.Fprintf(writer, "\n%s\n", item.Source.Name)
		fmt.Fprintln(writer, "  COMPARISON")
		fmt.Fprintf(writer, "    status:   %s\n", item.Status)
		if item.CurrentLabel != "" {
			fmt.Fprintf(writer, "    current:  %s  [college: %s]\n", item.CurrentLabel, item.CurrentGroup)
		}
		fmt.Fprintf(writer, "    evidence: %s\n", item.Evidence)
		fmt.Fprintln(writer, "  SOURCE")
		fmt.Fprintf(writer, "    college:  %s\n", item.Source.ParentCollegeName)
		fmt.Fprintf(writer, "    id:       %s\n", item.Source.ID)
		fmt.Fprintf(writer, "    rule:     %s\n", item.SourceExplanation)
	}
}

func printCurrentOnlyItems(writer *bufio.Writer, comparison *uiucAttributeComparison) {
	fmt.Fprintln(writer, "\n================================================================================")
	fmt.Fprintln(writer, "CURRENT DATABASE ITEMS WITH NO RELIABLE LIVE MATCH")
	fmt.Fprintln(writer, "================================================================================")
	fmt.Fprintln(writer, "These are not automatically declared removed. They can also be renamed or part of a split/merge.")

	fmt.Fprintf(writer, "\nCOLLEGES (%d)\n", len(comparison.CurrentOnlyColleges))
	for _, college := range comparison.CurrentOnlyColleges {
		fmt.Fprintf(writer, "\n%s\n", college.displayName())
		fmt.Fprintln(writer, "  COMPARISON")
		fmt.Fprintln(writer, "    status:   CURRENT ONLY / REVIEW")
		if college.Value != college.displayName() {
			fmt.Fprintf(writer, "    value:    %s\n", college.Value)
		}
		fmt.Fprintln(writer, "    evidence: no reliable live-source item was matched to this current value")
	}

	fmt.Fprintf(writer, "\nDEPARTMENTS (%d)\n", len(comparison.CurrentOnlyDepartments))
	for _, department := range comparison.CurrentOnlyDepartments {
		fmt.Fprintf(writer, "\n%s\n", department.Label)
		fmt.Fprintln(writer, "  COMPARISON")
		fmt.Fprintln(writer, "    status:   CURRENT ONLY / REVIEW")
		fmt.Fprintf(writer, "    college:  %s\n", department.Group)
		fmt.Fprintln(writer, "    evidence: no reliable live-source item was matched to this current value")
	}
}

func countCollegeStatuses(items []uiucCollegeComparison) map[string]int {
	result := make(map[string]int)
	for _, item := range items {
		result[item.Status]++
	}
	return result
}

func countDepartmentStatuses(items []uiucDepartmentComparison) map[string]int {
	result := make(map[string]int)
	for _, item := range items {
		result[item.Status]++
	}
	return result
}

func uniqueNonEmptyStrings(values ...string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]bool)
	for _, value := range values {
		if value != "" && !seen[value] {
			result = append(result, value)
			seen[value] = true
		}
	}
	return result
}
