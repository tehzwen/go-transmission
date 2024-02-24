package go_transmission

import (
	"strings"
)

type Torrent struct {
	Id       int     `json:"id,omitempty"`
	Name     string  `json:"name,omitempty"`
	Progress float32 `json:"percentDone,omitempty"`
	Status   int     `json:"status,omitempty"`
}

func WithCleanedMagnet() func(*TransmissionRequest) {
	return func(tr *TransmissionRequest) {
		if tr.Arguments.Filename != "" {
			tr.Arguments.Filename = cleanMagnet(tr.Arguments.Filename)
		}
	}
}

// cleanMagnet takes a full torrent magnet link and returns the
// link without any tracker urls
func cleanMagnet(original string) string {
	parts := strings.Split(original, "&tr=")
	return parts[0]
}
