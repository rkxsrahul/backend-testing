package method

import (
	"bytes"
	"fmt"
	"strconv"
)

// percentEncode percent encodes a string according to RFC 3986 2.1.
func percentEncode(input string) string {
	var buf bytes.Buffer
	for _, b := range []byte(input) {
		// if in unreserved set
		ok, _ := strconv.ParseBool(shouldEscape(b))
		if ok {
			buf.Write([]byte(fmt.Sprintf("%%%02X", b)))
		} else {
			// do not escape, write byte as-is
			buf.WriteByte(b)
		}
	}
	return buf.String()
}

// shouldEscape returns false if the byte is an unreserved character that
// should not be escaped and true otherwise, according to RFC 3986 2.1.

func shouldEscape(c byte) string {
	// RFC3986 2.3 unreserved characters
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return "false"
	}
	switch c {
	case '-', '.', '_', '~':
		return "false"
	}
	// all other bytes must be escaped
	return "true"
}
