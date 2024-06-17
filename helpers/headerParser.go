package helpers

import "strings"

func ParseHeaders(headerStr string) map[string]string {
	headers := make(map[string]string)
	for _, h := range strings.Split(headerStr, "\n") {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return headers
}
