package bitwarden

import (
	"context"
	"reflect"
	"testing"
	"time"
)

type decryptStructTestTarget struct {
	Nested *decryptStructTestNested
	When   *time.Time
}

type decryptStructTestNested struct {
	Value string
}

func TestDecryptStructSkipsExternalStructFields(t *testing.T) {
	now := time.Now()
	target := &decryptStructTestTarget{
		Nested: &decryptStructTestNested{},
		When:   &now,
	}

	if err := (&vault{}).decryptStruct(context.Background(), target, nil, nil); err != nil {
		t.Fatalf("decryptStruct() returned error: %v", err)
	}
}

func TestIsBasePackageType(t *testing.T) {
	if !isBasePackageType(reflect.TypeFor[decryptStructTestNested]()) {
		t.Fatal("local package type was not treated as a base package type")
	}
	if isBasePackageType(reflect.TypeFor[time.Time]()) {
		t.Fatal("external package type was treated as a base package type")
	}
}
