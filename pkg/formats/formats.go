package formats

import "strings"

// AudioFormat represents an audio format with its properties
type AudioFormat struct {
	Name        string
	Extensions  []string
	Codec       string
	Description string
	TypicalUse  string
	IsLossy     bool
}

// SupportedFormats contains all supported audio formats
var SupportedFormats = map[string]AudioFormat{
	"mp3": {
		Name:        "MP3",
		Extensions:  []string{".mp3"},
		Codec:       "libmp3lame",
		Description: "MPEG Audio Layer 3, lossy",
		TypicalUse:  "música, podcasts",
		IsLossy:     true,
	},
	"aac": {
		Name:        "AAC",
		Extensions:  []string{".aac", ".m4a"},
		Codec:       "aac",
		Description: "Advanced Audio Coding, lossy",
		TypicalUse:  "streaming, apple devices",
		IsLossy:     true,
	},
	"flac": {
		Name:        "FLAC",
		Extensions:  []string{".flac"},
		Codec:       "flac",
		Description: "Free Lossless Audio Codec",
		TypicalUse:  "archival, audiófilos",
		IsLossy:     false,
	},
	"wav": {
		Name:        "WAV",
		Extensions:  []string{".wav"},
		Codec:       "pcm_s16le",
		Description: "Waveform Audio, sin compresión",
		TypicalUse:  "producción, edición",
		IsLossy:     false,
	},
	"ogg": {
		Name:        "OGG Vorbis",
		Extensions:  []string{".ogg"},
		Codec:       "libvorbis",
		Description: "Ogg container con Vorbis codec, lossy",
		TypicalUse:  "open source, gaming",
		IsLossy:     true,
	},
	"opus": {
		Name:        "Opus",
		Extensions:  []string{".opus"},
		Codec:       "libopus",
		Description: "Codec moderno de alta eficiencia, lossy",
		TypicalUse:  "VoIP, streaming, mejor calidad/tamaño",
		IsLossy:     true,
	},
	"wma": {
		Name:        "WMA",
		Extensions:  []string{".wma"},
		Codec:       "wmav2",
		Description: "Windows Media Audio, lossy",
		TypicalUse:  "legacy windows",
		IsLossy:     true,
	},
	"alac": {
		Name:        "ALAC",
		Extensions:  []string{".m4a"},
		Codec:       "alac",
		Description: "Apple Lossless Audio Codec",
		TypicalUse:  "Apple ecosystem, lossless",
		IsLossy:     false,
	},
}

// GetFormatFromExtension returns the format name from file extension
func GetFormatFromExtension(filename string) (string, bool) {
	if !strings.Contains(filename, ".") {
		return "", false
	}
	ext := strings.ToLower(strings.TrimPrefix(strings.ToLower(filename[strings.LastIndex(filename, "."):]), "."))
	for format, info := range SupportedFormats {
		for _, e := range info.Extensions {
			if strings.TrimPrefix(e, ".") == ext {
				return format, true
			}
		}
	}
	return "", false
}

// IsSupportedFormat checks if a format is supported
func IsSupportedFormat(format string) bool {
	_, exists := SupportedFormats[strings.ToLower(format)]
	return exists
}

// GetSupportedFormats returns a list of all supported format names
func GetSupportedFormats() []string {
	var formats []string
	for format := range SupportedFormats {
		formats = append(formats, format)
	}
	return formats
}
