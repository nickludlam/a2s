package a2s

import (
	"bytes"
	"os"
	"testing"
)

// TestHexDump tests the hexDump function with various scenarios
func TestHexDump(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		highlightStart int
		highlightEnd   int
		wantContains   []string
	}{
		{
			name:           "empty data",
			data:           []byte{},
			highlightStart: 0,
			highlightEnd:   0,
			wantContains:   []string{},
		},
		{
			name:           "small data with highlight",
			data:           []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x00},
			highlightStart: 0,
			highlightEnd:   4,
			wantContains:   []string{"Raw response hex dump", "0x0000", "FF"},
		},
		{
			name:           "player response with wrong header",
			data:           []byte{0xAB, 0xCD, 0xEF, 0x12, 0x44}, // Wrong header, correct player response byte
			highlightStart: 0,
			highlightEnd:   4,
			wantContains:   []string{"Raw response hex dump", "0x0000", "AB", "CD", "EF"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hexDump(tt.data, tt.highlightStart, tt.highlightEnd)

			if len(tt.data) == 0 {
				if got != "" {
					t.Errorf("hexDump() with empty data should return empty string, got %q", got)
				}
				return
			}

			for _, want := range tt.wantContains {
				if !bytes.Contains([]byte(got), []byte(want)) {
					t.Errorf("hexDump() output missing expected string %q\nGot:\n%s", want, got)
				}
			}
		})
	}
}

// TestValidatorWithHexDump tests that validators output hex dumps when ShowRawResponses is true
func TestValidatorWithHexDump(t *testing.T) {
	// Test isMultiPacket with wrong header
	wrongHeaderData := []byte{0xAB, 0xCD, 0xEF, 0x12, 0x44}

	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Call with ShowRawResponses = true
	_, err := isMultiPacket(wrongHeaderData, true)

	// Restore stderr
	w.Close()
	os.Stderr = oldStderr

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if err == nil {
		t.Error("isMultiPacket() should return error for wrong header")
	}

	if !bytes.Contains([]byte(output), []byte("Raw response hex dump")) {
		t.Errorf("Expected hex dump in output when ShowRawResponses=true, got:\n%s", output)
	}

	// Check that the hex bytes are present (don't check exact format due to ANSI codes)
	if !bytes.Contains([]byte(output), []byte("AB")) || !bytes.Contains([]byte(output), []byte("CD")) {
		t.Errorf("Expected header bytes in hex dump, got:\n%s", output)
	}
}

// TestValidatorWithoutHexDump tests that validators don't output hex dumps when ShowRawResponses is false
func TestValidatorWithoutHexDump(t *testing.T) {
	wrongHeaderData := []byte{0xAB, 0xCD, 0xEF, 0x12, 0x44}

	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Call with ShowRawResponses = false
	_, err := isMultiPacket(wrongHeaderData, false)

	// Restore stderr
	w.Close()
	os.Stderr = oldStderr

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if err == nil {
		t.Error("isMultiPacket() should return error for wrong header")
	}

	if bytes.Contains([]byte(output), []byte("Raw response hex dump")) {
		t.Errorf("Should not output hex dump when ShowRawResponses=false, got:\n%s", output)
	}
}

// TestValidateResponseTypeWithHexDump tests validateResponseType with hex dump output
func TestValidateResponseTypeWithHexDump(t *testing.T) {
	// Create a response with wrong player response type
	wrongResponseData := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x99, 0x00} // Byte 4 is wrong

	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Call with ShowRawResponses = true
	err := validateResponseType(PlayerRequest, Flag(0x99), wrongResponseData, true)

	// Restore stderr
	w.Close()
	os.Stderr = oldStderr

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if err == nil {
		t.Error("validateResponseType() should return error for wrong response type")
	}

	if !bytes.Contains([]byte(output), []byte("Raw response hex dump")) {
		t.Errorf("Expected hex dump in output when ShowRawResponses=true, got:\n%s", output)
	}
}
