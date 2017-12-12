package ksync

import (
	"regexp"
)

// LineStatus makes it easy to set status based off lines of input.
type LineStatus struct {
	statuses map[ServiceStatus]*regexp.Regexp
}

// NewLineStatus create a LineStatus that is ready to receive lines.
func NewLineStatus(statuses map[ServiceStatus]string) (*LineStatus, error) {
	lineHandler := &LineStatus{statuses: map[ServiceStatus]*regexp.Regexp{}}

	for status, line := range statuses {
		re, err := regexp.Compile(line)
		if err != nil {
			return nil, err
		}

		lineHandler.statuses[status] = re
	}

	return lineHandler, nil
}

// Next takes a line of text and returns the status corresponding to that line.
func (ls *LineStatus) Next(line string) ServiceStatus {
	for status, re := range ls.statuses {
		if re.MatchString(line) {
			return status
		}
	}

	return ""
}
