package exportlistmapping

import (
	"testing"
)

func TestNewSetting(t *testing.T) {
	const jsonPath = "setting.json"

	s, err := NewSetting(jsonPath)
	if err != nil {
		t.Errorf("NewSetting: %v\n", err)
		return
	}

	if s.Sheet != "Sheet1" {
		t.Errorf("Sheet: %v\n", s.Sheet)
		return
	}

	if s.Kata != 3 {
		t.Errorf("Kata: %v\n", s.Kata)
		return
	}

	if s.Start != 8 {
		t.Errorf("Start: %v\n", s.Start)
		return
	}

	if s.Date.Row != 0 {
		t.Errorf("Date.Row: %v\n", s.Date.Row)
		return
	}
}
