package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type (
	WindAlarmAnalysis struct {
		Id           int64
		TraceId      string
		FarmCode     string
		TowerCode    string
		AlarmCode    string
		Title        string
		Content      string
		EvidenceJson string
		Status       string
	}

	WindHealthReport struct {
		Id           int64
		TraceId      string
		ReportType   string
		FarmCode     string
		TowerCode    string
		Title        string
		Content      string
		EvidenceJson string
		Status       string
	}

	WindMaintenanceTicketDraft struct {
		Id           int64
		TraceId      string
		FarmCode     string
		TowerCode    string
		AlarmCode    string
		Title        string
		Content      string
		EvidenceJson string
		Status       string
	}
)

func WindEvidenceJSON(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func WindTraceID() string {
	return fmt.Sprintf("wind-%d", time.Now().UnixNano())
}
