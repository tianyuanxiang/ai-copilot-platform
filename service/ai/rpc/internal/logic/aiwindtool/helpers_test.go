package aiwindtool

import (
	"testing"
	"time"
)

func TestDecodeArgsEmptyJSON(t *testing.T) {
	var args GetTurbineMetadataArgs
	if err := decodeArgs("", &args); err != nil {
		t.Fatalf("decodeArgs() error = %v", err)
	}
}

func TestDecodeArgsRejectsInvalidJSON(t *testing.T) {
	var args GetTurbineMetadataArgs
	if err := decodeArgs("{", &args); err == nil {
		t.Fatal("decodeArgs() expected error, got nil")
	}
}

func TestNormalizeTimeRangeDefaultsToLast24Hours(t *testing.T) {
	startRaw, endRaw, err := normalizeTimeRange("", "")
	if err != nil {
		t.Fatalf("normalizeTimeRange() error = %v", err)
	}
	start, err := time.Parse(toolTimeLayout, startRaw)
	if err != nil {
		t.Fatalf("parse start time: %v", err)
	}
	end, err := time.Parse(toolTimeLayout, endRaw)
	if err != nil {
		t.Fatalf("parse end time: %v", err)
	}
	if diff := end.Sub(start); diff != 24*time.Hour {
		t.Fatalf("default range = %v, want %v", diff, 24*time.Hour)
	}
}

func TestNormalizeTimeRangeRequiresBothEnds(t *testing.T) {
	if _, _, err := normalizeTimeRange("2026-05-01 00:00:00", ""); err == nil {
		t.Fatal("normalizeTimeRange() expected error for missing end time")
	}
	if _, _, err := normalizeTimeRange("", "2026-05-01 00:00:00"); err == nil {
		t.Fatal("normalizeTimeRange() expected error for missing start time")
	}
}

func TestNormalizeTimeRangeRejectsMoreThan31Days(t *testing.T) {
	if _, _, err := normalizeTimeRange("2026-01-01 00:00:00", "2026-02-02 00:00:00"); err == nil {
		t.Fatal("normalizeTimeRange() expected error for range over 31 days")
	}
}

func TestNormalizePageCapsPageSize(t *testing.T) {
	page, pageSize := normalizePage(0, 101)
	if page != 1 {
		t.Fatalf("page = %d, want 1", page)
	}
	if pageSize != maxPageSize {
		t.Fatalf("pageSize = %d, want %d", pageSize, maxPageSize)
	}
}

func TestValidateFieldsRejectsMoreThanEightFields(t *testing.T) {
	fields := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	if err := validateFields(fields); err == nil {
		t.Fatal("validateFields() expected error, got nil")
	}
}
