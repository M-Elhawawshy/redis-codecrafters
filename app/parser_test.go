package main

import "testing"

func TestParseRESP(t *testing.T) {
	input := "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"
	expected := []string{"ECHO", "hey"}

	result, err := parseRESP(input)
	if err != nil {
		t.Fatalf("parsing failed: %v", err)
	}
	if len(result) != len(expected) {
		t.Fatalf("expected %d elements, got %d", len(expected), len(result))
	}

	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("element %d mismatch: expected %q, got %q", i, expected[i], result[i])
		}
	}
}
