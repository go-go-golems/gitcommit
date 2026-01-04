package commitmsg

import "strings"

func EnsureTicketPrefix(ticketID string, message string) string {
	message = strings.TrimSpace(message)
	if message == "" {
		return message
	}

	ticketID = strings.TrimSpace(ticketID)
	if ticketID == "" {
		return message
	}

	candidates := []string{
		ticketID + ":",
		ticketID + " ",
		"[" + ticketID + "]",
		"(" + ticketID + ")",
	}
	for _, c := range candidates {
		if strings.HasPrefix(message, c) {
			return message
		}
	}

	return ticketID + ": " + message
}
