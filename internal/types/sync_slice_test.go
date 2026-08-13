package types

import (
	"errors"
	"testing"
)

func TestSyncSlice(t *testing.T) {
	t.Parallel()

	values := NewSyncSlice[int](4, 2)
	if values.Len() != 2 {
		t.Fatalf("Len() = %d, want 2", values.Len())
	}
	if err := values.Insert(0, 10); err != nil {
		t.Fatalf("Insert(0) returned error: %v", err)
	}
	values.Append(30)

	if got, ok := values.At(0); !ok || got != 10 {
		t.Fatalf("At(0) = %d, %v, want 10, true", got, ok)
	}
	if got, ok := values.At(1); !ok || got != 0 {
		t.Fatalf("At(1) = %d, %v, want zero value, true", got, ok)
	}
	if got, ok := values.At(2); !ok || got != 30 {
		t.Fatalf("At(2) = %d, %v, want 30, true", got, ok)
	}
	if got, ok := values.At(3); ok || got != 0 {
		t.Fatalf("At(3) = %d, %v, want zero value, false", got, ok)
	}

	snapshot := values.ToSlice()
	if len(snapshot) != 3 || snapshot[0] != 10 || snapshot[1] != 0 || snapshot[2] != 30 {
		t.Fatalf("ToSlice() = %#v, want [10 0 30]", snapshot)
	}
}

func TestSyncSliceInsertTooLarge(t *testing.T) {
	t.Parallel()

	values := NewSyncSlice[int](1, 1)
	err := values.Insert(1, 10)
	if !errors.Is(err, ErrSliceTooSmall) {
		t.Fatalf("Insert() error = %v, want ErrSliceTooSmall", err)
	}
}
