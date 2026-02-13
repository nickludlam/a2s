package a2s

import (
	"fmt"
	"strings"
)

// hexDump creates a hex+ASCII dump of data with highlighted bytes.
// highlightStart and highlightEnd define the range of bytes to highlight in bold.
// If highlightEnd is -1, highlights from highlightStart to the end.
func hexDump(data []byte, highlightStart, highlightEnd int) string {
	if len(data) == 0 {
		return ""
	}

	// If highlightEnd is -1, highlight to the end
	if highlightEnd == -1 {
		highlightEnd = len(data)
	}

	var sb strings.Builder
	sb.WriteString("\nRaw response hex dump:\n")
	sb.WriteString("Offset  00 01 02 03 04 05 06 07  08 09 0A 0B 0C 0D 0E 0F  ASCII\n")
	sb.WriteString("------  -----------------------------------------------  ----------------\n")

	for offset := 0; offset < len(data); offset += 16 {
		// Offset
		sb.WriteString(fmt.Sprintf("0x%04X  ", offset))

		// Hex bytes (first 8 bytes)
		for i := 0; i < 8; i++ {
			pos := offset + i
			if pos < len(data) {
				if pos >= highlightStart && pos < highlightEnd {
					sb.WriteString(fmt.Sprintf("\033[1m%02X\033[0m ", data[pos]))
				} else {
					sb.WriteString(fmt.Sprintf("%02X ", data[pos]))
				}
			} else {
				sb.WriteString("   ")
			}
		}

		sb.WriteString(" ")

		// Hex bytes (second 8 bytes)
		for i := 8; i < 16; i++ {
			pos := offset + i
			if pos < len(data) {
				if pos >= highlightStart && pos < highlightEnd {
					sb.WriteString(fmt.Sprintf("\033[1m%02X\033[0m ", data[pos]))
				} else {
					sb.WriteString(fmt.Sprintf("%02X ", data[pos]))
				}
			} else {
				sb.WriteString("   ")
			}
		}

		sb.WriteString(" ")

		// ASCII representation
		for i := 0; i < 16; i++ {
			pos := offset + i
			if pos < len(data) {
				b := data[pos]
				var char string
				if b >= 32 && b <= 126 {
					char = string(b)
				} else {
					char = "."
				}
				if pos >= highlightStart && pos < highlightEnd {
					sb.WriteString(fmt.Sprintf("\033[1m%s\033[0m", char))
				} else {
					sb.WriteString(char)
				}
			}
		}

		sb.WriteString("\n")
	}

	return sb.String()
}
