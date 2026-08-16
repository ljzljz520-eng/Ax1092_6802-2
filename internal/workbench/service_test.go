package workbench

import "testing"

func TestTransitionAllowed(t *testing.T) {
	tests := []struct {
		name    string
		current Status
		target  Status
		want    bool
	}{
		{name: "draft to review", current: StatusDraft, target: StatusPendingReview, want: true},
		{name: "returned to review", current: StatusReturned, target: StatusPendingReview, want: true},
		{name: "review to published", current: StatusPendingReview, target: StatusPublished, want: true},
		{name: "review to returned", current: StatusPendingReview, target: StatusReturned, want: true},
		{name: "published to archived", current: StatusPublished, target: StatusArchived, want: true},
		{name: "archived to completed", current: StatusArchived, target: StatusCompleted, want: false},
		{name: "completed to archived", current: StatusCompleted, target: StatusArchived, want: false},
		{name: "same status", current: StatusDraft, target: StatusDraft, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := transitionAllowed(test.current, test.target); got != test.want {
				t.Fatalf("transitionAllowed(%q, %q) = %t, want %t", test.current, test.target, got, test.want)
			}
		})
	}
}
