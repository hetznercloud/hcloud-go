package hcloud

import (
	"encoding/json"
	"testing"
)

func TestServerTypeUnmarshalJSON(t *testing.T) {
	data := []byte(`{
		"id": 1,
		"name": "cx10",
		"description": "description",
		"cores": 4,
		"memory": 1.0,
		"disk": 20,
		"storage_type": "local"
	}`)

	var serverType ServerType
	if err := json.Unmarshal(data, &serverType); err != nil {
		t.Fatal(err)
	}

	if serverType.ID != 1 {
		t.Errorf("unexpected ID: %v", serverType.ID)
	}
	if serverType.Name != "cx10" {
		t.Errorf("unexpected name: %q", serverType.Name)
	}
	if serverType.Description != "description" {
		t.Errorf("unexpected description: %q", serverType.Description)
	}
	if serverType.Cores != 4 {
		t.Errorf("unexpected cores: %v", serverType.Cores)
	}
	if serverType.Memory != 1.0 {
		t.Errorf("unexpected memory: %v", serverType.Memory)
	}
	if serverType.Disk != 20 {
		t.Errorf("unexpected disk: %v", serverType.Disk)
	}
	if serverType.StorageType != StorageTypeLocal {
		t.Errorf("unexpected storage type: %q", serverType.StorageType)
	}
}
