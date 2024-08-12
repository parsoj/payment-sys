package id

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	generator := &SpecialIdGenerator{}

	id, err := generator.New()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(id) > 20 {
		t.Fatalf("Expected ID length <= 20, got %d", len(id))
	}

	if valid, err := generator.Validate(id); !valid || err != nil {
		fmt.Println(id)
		t.Fatalf("Expected valid ID, got invalid with error %v", err)
	}
}

func TestFromString(t *testing.T) {
	generator := &SpecialIdGenerator{}

	idStr := "abcd1234efgh5678ijkl"
	id, err := generator.FromString(idStr)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if string(id) != idStr {
		t.Fatalf("Expected ID %s, got %s", idStr, id)
	}

	if valid, err := generator.Validate(id); !valid || err != nil {
		t.Fatalf("Expected valid ID, got invalid with error %v", err)
	}
}

func TestFromBytes(t *testing.T) {
	generator := &SpecialIdGenerator{}

	bytes := []byte{0xde, 0xad, 0xbe, 0xef}
	id, err := generator.FromBytes(bytes)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(id) > 20 {
		t.Fatalf("Expected ID length <= 20, got %d", len(id))
	}

	if valid, err := generator.Validate(id); !valid || err != nil {
		t.Fatalf("Expected valid ID, got invalid with error %v", err)
	}
}

func TestValidate(t *testing.T) {
	generator := &SpecialIdGenerator{}

	validID := "abcd1234efgh5678ijkl"
	invalidID := "abcd1234efgh5678ijklmnopqrstuvwxyz" // 26 characters

	if valid, err := generator.Validate(SpecialId(validID)); !valid || err != nil {
		t.Fatalf("Expected valid ID, got invalid with error %v", err)
	}

	if valid, _ := generator.Validate(SpecialId(invalidID)); valid {
		t.Fatalf("Expected invalid ID due to length, but got valid")
	}
}
