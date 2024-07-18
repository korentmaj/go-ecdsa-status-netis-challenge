package status

import (
	"testing"
)

func TestStatusList(t *testing.T) {
	sl := NewStatusList()

	// Add statuses
	index1 := sl.AddStatus(true)
	index2 := sl.AddStatus(false)

	if index1 != 0 {
		t.Fatalf("Expected index1 to be 0, got %d", index1)
	}

	if index2 != 8 {
		t.Fatalf("Expected index2 to be 8, got %d", index2)
	}

	// Set statuses
	if err := sl.SetStatus(index2, true); err != nil {
		t.Fatalf("Error setting status: %v", err)
	}

	// Encode
	encoded, err := sl.Encode()
	if err != nil {
		t.Fatalf("Error encoding status list: %v", err)
	}

	if encoded == "" {
		t.Fatalf("Expected encoded string to be non-empty")
	}
}
