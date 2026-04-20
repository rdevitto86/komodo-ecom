package models

import (
	"encoding/json"
	"testing"
)

// TestServiceUnmarshalJSON_defaultsServiceType verifies that existing S3 data
// lacking a service_type field is defaulted to "service" on unmarshal.
func TestServiceUnmarshalJSON_defaultsServiceType(t *testing.T) {
	t.Parallel()

	raw := `{"id":"svc-1","slug":"haircut","name":"Haircut","description":"A haircut","category":"grooming","status":"active","price":25.00,"locationTypes":["in-store"]}`
	var svc Service
	if err := json.Unmarshal([]byte(raw), &svc); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if svc.ServiceType != ServiceTypeService {
		t.Errorf("ServiceType: got %q, want %q", svc.ServiceType, ServiceTypeService)
	}
}

// TestServiceUnmarshalJSON_preservesExplicitServiceType verifies that an
// explicitly set service_type is not overwritten.
func TestServiceUnmarshalJSON_preservesExplicitServiceType(t *testing.T) {
	t.Parallel()

	raw := `{"id":"rep-1","slug":"screen-repair","name":"Screen Repair","description":"Fixes cracked screens","category":"repair","status":"active","price":99.00,"locationTypes":["in-store"],"service_type":"repair","accepted_device_types":["phone","tablet"],"estimated_turnaround_days":2,"warranty_on_repair":"90 days"}`
	var svc Service
	if err := json.Unmarshal([]byte(raw), &svc); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if svc.ServiceType != ServiceTypeRepair {
		t.Errorf("ServiceType: got %q, want %q", svc.ServiceType, ServiceTypeRepair)
	}
	if len(svc.AcceptedDeviceTypes) != 2 {
		t.Errorf("AcceptedDeviceTypes: got %v, want 2 entries", svc.AcceptedDeviceTypes)
	}
	if svc.EstimatedTurnaroundDays != 2 {
		t.Errorf("EstimatedTurnaroundDays: got %d, want 2", svc.EstimatedTurnaroundDays)
	}
	if svc.WarrantyOnRepair != "90 days" {
		t.Errorf("WarrantyOnRepair: got %q, want %q", svc.WarrantyOnRepair, "90 days")
	}
}

// TestServiceUnmarshalJSON_roundTrip verifies that marshalling a Service with
// repair fields and unmarshalling it produces identical output.
func TestServiceUnmarshalJSON_roundTrip(t *testing.T) {
	t.Parallel()

	original := Service{
		ID:                      "rep-2",
		Slug:                    "battery-replacement",
		Name:                    "Battery Replacement",
		Description:             "Replace phone battery",
		Category:                "repair",
		Status:                  "active",
		Price:                   49.99,
		LocationTypes:           []string{"in-store"},
		ServiceType:             ServiceTypeRepair,
		AcceptedDeviceTypes:     []string{"phone"},
		EstimatedTurnaroundDays: 1,
		WarrantyOnRepair:        "30 days",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var restored Service
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if restored.ServiceType != original.ServiceType {
		t.Errorf("ServiceType mismatch: got %q, want %q", restored.ServiceType, original.ServiceType)
	}
	if restored.EstimatedTurnaroundDays != original.EstimatedTurnaroundDays {
		t.Errorf("EstimatedTurnaroundDays mismatch: got %d, want %d", restored.EstimatedTurnaroundDays, original.EstimatedTurnaroundDays)
	}
}
