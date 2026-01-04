package ticket

import (
	"regexp"
	"strings"
)

var ticketRe = regexp.MustCompile(`(?i)\b([a-z][a-z0-9]+-[a-z0-9]+)\b`)

func FromBranch(branch string) (string, bool) {
	m := ticketRe.FindStringSubmatch(branch)
	if len(m) < 2 {
		return "", false
	}
	return strings.ToUpper(m[1]), true
}
