package hcloud

import (
	"encoding/json"
	"testing"
)

func TestISOUnmarshalJSON(t *testing.T) {
	data := []byte(`{
		"id": 4711,
		"name": "FreeBSD-11.0-RELEASE-amd64-dvd1",
		"description": "FreeBSD 11.0 x64",
		"type": "public"
	}`)

	var v ISO
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}

	if v.ID != 4711 {
		t.Errorf("unexpected ID: %v", v.ID)
	}
	if v.Name != "FreeBSD-11.0-RELEASE-amd64-dvd1" {
		t.Errorf("unexpected name: %v", v.Name)
	}
	if v.Description != "FreeBSD 11.0 x64" {
		t.Errorf("unexpected description: %v", v.Description)
	}
	if v.Type != ISOTypePublic {
		t.Errorf("unexpected type: %v", v.Type)
	}
}
