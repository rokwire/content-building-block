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
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	uiucDMIDirectoryURL = "https://www.dmi.illinois.edu/ddd/mktextdirectory.asp"
	uiucCampusCode      = "1"
	uiucCampusID        = "illinois.edu"
)

// uiucDMIPrototypeOnly temporarily stops the normal service startup after the
// report is printed. Set it to false when this prototype is integrated into the
// normal Content BB lifecycle.
var uiucDMIPrototypeOnly = true

var (
	uiucUnitIDPattern       = regexp.MustCompile(`^[A-Z]{2}$`)
	uiucDepartmentIDPattern = regexp.MustCompile(`^\d{3}$`)
)

type uiucDMICollege struct {
	ID                         string
	CampusID                   string
	UnitID                     string
	Name                       string
	DirectCollegeOrAdminLabels []string
	DirectSourceLines          []int
	ParentSourceLines          []int
	SourceLines                []int
	Reasons                    []string
}

type uiucDMIDepartment struct {
	ID                string
	CampusID          string
	UnitID            string
	DepartmentID      string
	Name              string
	FullName          string
	ParentCollegeID   string
	ParentCollegeName string
	SourceLine        int
	Reason            string
}

type uiucDMIAnalysis struct {
	SourceURL            string
	FetchedAt            time.Time
	LastModified         string
	SourceRowCount       int
	IncludedSourceRows   int
	ExcludedSourceRows   int
	ExcludedReasonCounts map[string]int
	Colleges             []uiucDMICollege
	Departments          []uiucDMIDepartment
}

type uiucDMIRawRow struct {
	SourceLine   int
	LegacyCode   string
	DeptName     string
	DeptFullName string
	Dean         string
	Active       string
	Campus       string
	College      string
	CollegeName  string
	Dept         string
}

type uiucDMITSVSchema struct {
	leftIndexes  map[string]int
	rightOffsets map[string]int
}

func runUIUCDMIPrototype(ctx context.Context, output io.Writer) error {
	current, err := loadCurrentUIUCAttributes(currentUIUCAttributesFile())
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	analysis, err := fetchUIUCDMIAnalysis(ctx, client, uiucDMIDirectoryURL)
	if err != nil {
		return err
	}
	comparison := compareUIUCAttributes(analysis, current)
	return printUIUCAttributeComparison(output, analysis, current, comparison)
}

func fetchUIUCDMIAnalysis(ctx context.Context, client *http.Client, sourceURL string) (*uiucDMIAnalysis, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create UIUC DMI request: %w", err)
	}
	request.Header.Set("Accept", "text/plain, text/tab-separated-values")
	request.Header.Set("User-Agent", "content-building-block-uiuc-dmi-prototype/1.0")

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("download UIUC DMI directory: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 4<<10))
		return nil, fmt.Errorf("download UIUC DMI directory: unexpected HTTP status %s", response.Status)
	}

	analysis, err := analyzeUIUCDMIDirectory(response.Body)
	if err != nil {
		return nil, fmt.Errorf("analyze UIUC DMI directory: %w", err)
	}
	analysis.SourceURL = sourceURL
	analysis.FetchedAt = time.Now().UTC()
	analysis.LastModified = response.Header.Get("Last-Modified")
	return analysis, nil
}

func analyzeUIUCDMIDirectory(input io.Reader) (*uiucDMIAnalysis, error) {
	reader := csv.NewReader(bufio.NewReader(input))
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read TSV header: %w", err)
	}
	if len(header) > 0 {
		header[0] = strings.TrimPrefix(header[0], "\ufeff")
	}

	schema, err := newUIUCDMITSVSchema(header)
	if err != nil {
		return nil, err
	}

	analysis := &uiucDMIAnalysis{ExcludedReasonCounts: make(map[string]int)}
	collegesByID := make(map[string]*uiucDMICollege)
	departmentsByID := make(map[string]uiucDMIDepartment)

	for sourceLine := 2; ; sourceLine++ {
		record, readErr := reader.Read()
		if errors.Is(readErr, io.EOF) {
			break
		}
		if readErr != nil {
			return nil, fmt.Errorf("read TSV source line %d: %w", sourceLine, readErr)
		}
		if isEmptyTSVRecord(record) {
			continue
		}

		raw, err := schema.parseRow(record, sourceLine)
		if err != nil {
			return nil, err
		}
		analysis.SourceRowCount++

		exclusionReasons := validateUIUCDMIRow(raw)
		if len(exclusionReasons) > 0 {
			analysis.ExcludedSourceRows++
			for _, reason := range exclusionReasons {
				analysis.ExcludedReasonCounts[reason]++
			}
			continue
		}
		analysis.IncludedSourceRows++

		collegeID := fmt.Sprintf("%s:%s", uiucCampusID, raw.College)
		college, exists := collegesByID[collegeID]
		if !exists {
			college = &uiucDMICollege{
				ID:       collegeID,
				CampusID: uiucCampusID,
				UnitID:   raw.College,
				Name:     raw.CollegeName,
			}
			collegesByID[collegeID] = college
		} else if college.Name != raw.CollegeName {
			return nil, fmt.Errorf(
				"source line %d: unit %s has conflicting names %q and %q",
				raw.SourceLine, raw.College, college.Name, raw.CollegeName,
			)
		}

		college.SourceLines = append(college.SourceLines, raw.SourceLine)
		if raw.Dean == "T" {
			college.DirectCollegeOrAdminLabels = appendUniqueString(college.DirectCollegeOrAdminLabels, raw.DeptName)
			college.DirectSourceLines = append(college.DirectSourceLines, raw.SourceLine)
			college.Reasons = append(college.Reasons, fmt.Sprintf(
				"source line %d: direct college/admin evidence because Campus=1, Active=Y and Dean=T; the stored identifier uses campus_id + unit_id and omits department_id; direct row label=%q",
				raw.SourceLine, raw.DeptName,
			))
			continue
		}

		college.ParentSourceLines = append(college.ParentSourceLines, raw.SourceLine)
		college.Reasons = append(college.Reasons, fmt.Sprintf(
			"source line %d: parent college is ensured from College=%s and CollegeName=%q because this included Dean=F department row requires that parent; the same row also creates department %s/%q",
			raw.SourceLine, raw.College, raw.CollegeName, raw.Dept, raw.DeptName,
		))

		departmentID := fmt.Sprintf("%s:%s", collegeID, raw.Dept)
		if existing, duplicate := departmentsByID[departmentID]; duplicate {
			return nil, fmt.Errorf(
				"source line %d: duplicate department identifier %s (already defined by source line %d)",
				raw.SourceLine, departmentID, existing.SourceLine,
			)
		}
		departmentsByID[departmentID] = uiucDMIDepartment{
			ID:                departmentID,
			CampusID:          uiucCampusID,
			UnitID:            raw.College,
			DepartmentID:      raw.Dept,
			Name:              raw.DeptName,
			FullName:          raw.DeptFullName,
			ParentCollegeID:   collegeID,
			ParentCollegeName: raw.CollegeName,
			SourceLine:        raw.SourceLine,
			Reason: fmt.Sprintf(
				"source line %d: included as department because Campus=1, Active=Y, Dean=F and Dept=%s is a three-digit department_id; parent college=%s/%q",
				raw.SourceLine, raw.Dept, raw.College, raw.CollegeName,
			),
		}
	}

	analysis.Colleges = make([]uiucDMICollege, 0, len(collegesByID))
	for _, college := range collegesByID {
		analysis.Colleges = append(analysis.Colleges, *college)
	}
	sort.Slice(analysis.Colleges, func(i, j int) bool {
		if analysis.Colleges[i].Name == analysis.Colleges[j].Name {
			return analysis.Colleges[i].UnitID < analysis.Colleges[j].UnitID
		}
		return analysis.Colleges[i].Name < analysis.Colleges[j].Name
	})

	analysis.Departments = make([]uiucDMIDepartment, 0, len(departmentsByID))
	for _, department := range departmentsByID {
		analysis.Departments = append(analysis.Departments, department)
	}
	sort.Slice(analysis.Departments, func(i, j int) bool {
		left := analysis.Departments[i]
		right := analysis.Departments[j]
		if left.ParentCollegeName != right.ParentCollegeName {
			return left.ParentCollegeName < right.ParentCollegeName
		}
		if left.Name != right.Name {
			return left.Name < right.Name
		}
		return left.ID < right.ID
	})

	return analysis, nil
}

func printUIUCDMIAnalysis(output io.Writer, analysis *uiucDMIAnalysis) error {
	writer := bufio.NewWriter(output)

	fmt.Fprintln(writer, "================================================================================")
	fmt.Fprintln(writer, "UIUC DMI COLLEGE / DEPARTMENT ANALYSIS")
	fmt.Fprintln(writer, "================================================================================")
	fmt.Fprintf(writer, "Source URL:           %s\n", analysis.SourceURL)
	fmt.Fprintf(writer, "Fetched at (UTC):     %s\n", analysis.FetchedAt.Format(time.RFC3339))
	if analysis.LastModified != "" {
		fmt.Fprintf(writer, "HTTP Last-Modified:   %s\n", analysis.LastModified)
	}
	fmt.Fprintf(writer, "Source rows:          %d\n", analysis.SourceRowCount)
	fmt.Fprintf(writer, "Included source rows: %d\n", analysis.IncludedSourceRows)
	fmt.Fprintf(writer, "Excluded source rows: %d\n", analysis.ExcludedSourceRows)
	fmt.Fprintf(writer, "Unique colleges:      %d\n", len(analysis.Colleges))
	fmt.Fprintf(writer, "Unique departments:   %d\n", len(analysis.Departments))

	if len(analysis.ExcludedReasonCounts) > 0 {
		fmt.Fprintln(writer, "\nExcluded source-row reasons:")
		reasons := make([]string, 0, len(analysis.ExcludedReasonCounts))
		for reason := range analysis.ExcludedReasonCounts {
			reasons = append(reasons, reason)
		}
		sort.Strings(reasons)
		for _, reason := range reasons {
			fmt.Fprintf(writer, "  - %s: %d row(s)\n", reason, analysis.ExcludedReasonCounts[reason])
		}
	}

	fmt.Fprintln(writer, "\n================================================================================")
	fmt.Fprintf(writer, "COLLEGES (%d)\n", len(analysis.Colleges))
	fmt.Fprintln(writer, "================================================================================")
	for index, college := range analysis.Colleges {
		fmt.Fprintf(writer, "\n[%03d] %s | %s\n", index+1, college.ID, college.Name)
		fmt.Fprintf(writer, "      campus_id=%s | unit_id=%s | all source lines=%s\n",
			college.CampusID, college.UnitID, formatSourceLines(college.SourceLines))
		if len(college.DirectCollegeOrAdminLabels) > 0 {
			fmt.Fprintf(writer, "      why: direct college/admin evidence from Dean=T source line(s) %s, label(s)=%q",
				formatSourceLines(college.DirectSourceLines), strings.Join(college.DirectCollegeOrAdminLabels, " | "))
			if len(college.ParentSourceLines) > 0 {
				fmt.Fprintf(writer, "; also required as parent by Dean=F department source line(s) %s", formatSourceLines(college.ParentSourceLines))
			}
			fmt.Fprintln(writer)
		} else {
			fmt.Fprintf(writer, "      why: no Dean=T direct row; included because it is the required parent college for Dean=F department source line(s) %s\n",
				formatSourceLines(college.ParentSourceLines))
		}
	}

	fmt.Fprintln(writer, "\n================================================================================")
	fmt.Fprintf(writer, "DEPARTMENTS (%d)\n", len(analysis.Departments))
	fmt.Fprintln(writer, "================================================================================")
	for index, department := range analysis.Departments {
		fmt.Fprintf(writer, "\n[%03d] %s | %s\n", index+1, department.ID, department.Name)
		fmt.Fprintf(writer, "      campus_id=%s | unit_id=%s | department_id=%s | parent=%s / %s | source line=%d\n",
			department.CampusID, department.UnitID, department.DepartmentID,
			department.ParentCollegeID, department.ParentCollegeName, department.SourceLine)
		fmt.Fprintf(writer, "      why: %s\n", department.Reason)
	}

	fmt.Fprintln(writer, "\n================================================================================")
	fmt.Fprintln(writer, "END OF UIUC DMI ANALYSIS")
	fmt.Fprintln(writer, "================================================================================")
	return writer.Flush()
}

func newUIUCDMITSVSchema(header []string) (*uiucDMITSVSchema, error) {
	indexes := make(map[string]int, len(header))
	for index, name := range header {
		indexes[strings.TrimSpace(name)] = index
	}

	leftFields := []string{"Legacy_Code", "Deptname", "deptFullName"}
	rightFields := []string{"Dean", "Active", "Campus", "College", "CollegeName", "Dept"}
	schema := &uiucDMITSVSchema{
		leftIndexes:  make(map[string]int, len(leftFields)),
		rightOffsets: make(map[string]int, len(rightFields)),
	}
	for _, field := range leftFields {
		index, ok := indexes[field]
		if !ok {
			return nil, fmt.Errorf("TSV header is missing required field %q", field)
		}
		schema.leftIndexes[field] = index
	}
	for _, field := range rightFields {
		index, ok := indexes[field]
		if !ok {
			return nil, fmt.Errorf("TSV header is missing required field %q", field)
		}
		schema.rightOffsets[field] = len(header) - 1 - index
	}
	return schema, nil
}

func (schema *uiucDMITSVSchema) parseRow(record []string, sourceLine int) (uiucDMIRawRow, error) {
	left := func(field string) (string, error) {
		index := schema.leftIndexes[field]
		if index < 0 || index >= len(record) {
			return "", fmt.Errorf("source line %d: missing left-anchored field %s", sourceLine, field)
		}
		return strings.TrimSpace(record[index]), nil
	}
	right := func(field string) (string, error) {
		index := len(record) - 1 - schema.rightOffsets[field]
		if index < 0 || index >= len(record) {
			return "", fmt.Errorf("source line %d: missing right-anchored field %s", sourceLine, field)
		}
		return strings.TrimSpace(record[index]), nil
	}

	legacyCode, err := left("Legacy_Code")
	if err != nil {
		return uiucDMIRawRow{}, err
	}
	deptName, err := left("Deptname")
	if err != nil {
		return uiucDMIRawRow{}, err
	}
	deptFullName, err := left("deptFullName")
	if err != nil {
		return uiucDMIRawRow{}, err
	}
	dean, err := right("Dean")
	if err != nil {
		return uiucDMIRawRow{}, err
	}
	active, err := right("Active")
	if err != nil {
		return uiucDMIRawRow{}, err
	}
	campus, err := right("Campus")
	if err != nil {
		return uiucDMIRawRow{}, err
	}
	college, err := right("College")
	if err != nil {
		return uiucDMIRawRow{}, err
	}
	collegeName, err := right("CollegeName")
	if err != nil {
		return uiucDMIRawRow{}, err
	}
	dept, err := right("Dept")
	if err != nil {
		return uiucDMIRawRow{}, err
	}

	return uiucDMIRawRow{
		SourceLine:   sourceLine,
		LegacyCode:   legacyCode,
		DeptName:     deptName,
		DeptFullName: deptFullName,
		Dean:         dean,
		Active:       active,
		Campus:       campus,
		College:      college,
		CollegeName:  collegeName,
		Dept:         dept,
	}, nil
}

func validateUIUCDMIRow(row uiucDMIRawRow) []string {
	reasons := make([]string, 0, 6)
	if row.Campus != uiucCampusCode {
		reasons = append(reasons, fmt.Sprintf("Campus=%s (expected 1 for UIUC)", row.Campus))
	}
	if row.Active != "Y" {
		reasons = append(reasons, fmt.Sprintf("Active=%s (expected Y)", row.Active))
	}
	if !uiucUnitIDPattern.MatchString(row.College) {
		reasons = append(reasons, fmt.Sprintf("College=%s (expected two uppercase characters)", row.College))
	}
	if row.CollegeName == "" {
		reasons = append(reasons, "CollegeName is empty")
	}
	if row.DeptName == "" {
		reasons = append(reasons, "Deptname is empty")
	}
	if row.Dean != "T" && row.Dean != "F" {
		reasons = append(reasons, fmt.Sprintf("Dean=%s (expected T or F)", row.Dean))
	}
	if row.Dean == "F" && !uiucDepartmentIDPattern.MatchString(row.Dept) {
		reasons = append(reasons, fmt.Sprintf("Dept=%s (expected a three-digit department_id)", row.Dept))
	}
	return reasons
}

func isEmptyTSVRecord(record []string) bool {
	for _, value := range record {
		if strings.TrimSpace(value) != "" {
			return false
		}
	}
	return true
}

func appendUniqueString(values []string, candidate string) []string {
	for _, value := range values {
		if value == candidate {
			return values
		}
	}
	return append(values, candidate)
}

func formatSourceLines(lines []int) string {
	values := make([]string, len(lines))
	for index, line := range lines {
		values[index] = fmt.Sprintf("%d", line)
	}
	return strings.Join(values, ", ")
}
