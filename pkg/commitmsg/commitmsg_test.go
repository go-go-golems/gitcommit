package commitmsg

import "testing"

func TestEnsureTicketPrefix(t *testing.T) {
	tests := []struct {
		name     string
		ticketID string
		message  string
		want     string
	}{
		{
			name:     "empty message",
			ticketID: "ABC-123",
			message:  "   ",
			want:     "",
		},
		{
			name:     "empty ticket leaves message",
			ticketID: "",
			message:  "Fix thing",
			want:     "Fix thing",
		},
		{
			name:     "prefix when missing",
			ticketID: "ABC-123",
			message:  "Fix thing",
			want:     "ABC-123: Fix thing",
		},
		{
			name:     "already prefixed with colon",
			ticketID: "ABC-123",
			message:  "ABC-123: Fix thing",
			want:     "ABC-123: Fix thing",
		},
		{
			name:     "already prefixed with bracket",
			ticketID: "ABC-123",
			message:  "[ABC-123] Fix thing",
			want:     "[ABC-123] Fix thing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EnsureTicketPrefix(tt.ticketID, tt.message)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
