package formats

import (
	"testing"
)

func TestGetFormatFromExtension(t *testing.T) {
	tests := []struct {
		filename string
		expected string
		want     bool
	}{
		{"song.mp3", "mp3", true},
		{"song.MP3", "mp3", true},
		{"song.flac", "flac", true},
		{"song.wav", "wav", true},
		{"song.ogg", "ogg", true},
		{"song.opus", "opus", true},
		{"song.m4a", "alac", true},
		{"song.wma", "wma", true},
		{"song.unknown", "", false},
		{"song", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got, ok := GetFormatFromExtension(tt.filename)
			if got != tt.expected || ok != tt.want {
				t.Errorf("GetFormatFromExtension(%s) = (%s, %v), want (%s, %v)",
					tt.filename, got, ok, tt.expected, tt.want)
			}
		})
	}
}

func TestIsSupportedFormat(t *testing.T) {
	tests := []struct {
		format string
		want   bool
	}{
		{"mp3", true},
		{"flac", true},
		{"wav", true},
		{"ogg", true},
		{"opus", true},
		{"aac", true},
		{"wma", true},
		{"alac", true},
		{"unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			if got := IsSupportedFormat(tt.format); got != tt.want {
				t.Errorf("IsSupportedFormat(%s) = %v, want %v", tt.format, got, tt.want)
			}
		})
	}
}

func TestGetSupportedFormats(t *testing.T) {
	formats := GetSupportedFormats()

	// Check that we have some expected formats
	expected := []string{"mp3", "flac", "wav", "ogg", "opus", "aac", "wma", "alac"}
	expectedMap := make(map[string]bool)
	for _, f := range expected {
		expectedMap[f] = true
	}

	if len(formats) != len(expected) {
		t.Errorf("GetSupportedFormats() returned %d formats, expected %d", len(formats), len(expected))
	}

	for _, format := range formats {
		if !expectedMap[format] {
			t.Errorf("GetSupportedFormats() returned unexpected format: %s", format)
		}
	}
}
