package ticket

import "testing"

func TestFromBranch(t *testing.T) {
	tests := []struct {
		name   string
		branch string
		want   string
		ok     bool
	}{
		{name: "simple", branch: "ABC-123", want: "ABC-123", ok: true},
		{name: "with prefix", branch: "feature/abc-123-add-thing", want: "ABC-123", ok: true},
		{name: "letters suffix", branch: "feature/GITCOMMIT-XXXX-something", want: "GITCOMMIT-XXXX", ok: true},
		{name: "no match", branch: "main", want: "", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := FromBranch(tt.branch)
			if ok != tt.ok {
				t.Fatalf("ok=%v, want %v (got=%q)", ok, tt.ok, got)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
