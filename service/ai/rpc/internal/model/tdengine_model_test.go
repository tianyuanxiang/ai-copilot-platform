package model

import (
	"context"
	"strings"
	"testing"
)

func TestSafeIdentifier(t *testing.T) {
	valid, err := SafeIdentifier("inclinometer")
	if err != nil {
		t.Fatalf("expected valid identifier: %v", err)
	}
	if valid != "inclinometer" {
		t.Fatalf("unexpected identifier: %s", valid)
	}

	if _, err := SafeIdentifier("alarm;drop table"); err == nil {
		t.Fatal("expected unsafe identifier to be rejected")
	}
}

func TestTDengineModelWithoutConnectionReturnsEmpty(t *testing.T) {
	model := NewTdengineModel(nil)
	rows, fields, err := model.QueryRows(context.Background(), "fuyu", "inclinometer", []string{"x"}, "1=1", 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rows != nil {
		t.Fatalf("expected nil rows when TDengine is not configured")
	}
	if len(fields) != 1 || fields[0] != "x" {
		t.Fatalf("unexpected fields: %#v", fields)
	}
}

func TestTDengineWhereHelpers(t *testing.T) {
	where := JoinWhere(append(TimeWhere("2026-05-25 00:00:00", "2026-05-26 00:00:00"), DeviceWhere("03", "7")...))
	for _, part := range []string{"ts >=", "ts <", "tower_id=3", "device_channel=7"} {
		if !strings.Contains(where, part) {
			t.Fatalf("expected where to contain %q, got %s", part, where)
		}
	}
}

func TestWindMetadataFieldsUseWhitelist(t *testing.T) {
	metadata := NewWindDeviceTypeModel(nil, nil)
	fields := metadata.FieldsForDeviceType(context.Background(), "INSX", "x")
	if len(fields) != 1 || fields[0] != "x" {
		t.Fatalf("expected whitelisted field x, got %#v", fields)
	}

	fields = metadata.FieldsForDeviceType(context.Background(), "INSX", "not_allowed")
	if len(fields) != 0 {
		t.Fatalf("expected non-whitelisted field to be rejected, got %#v", fields)
	}
}
