package workbench

import "strings"

type Section string

const (
	SectionNews      Section = "news"
	SectionInterview Section = "interview"
	SectionEditorial Section = "editorial"
	SectionEvent     Section = "event"
)

type Status string

const (
	StatusDraft         Status = "draft"
	StatusPendingReview Status = "pending_review"
	StatusPublished     Status = "published"
	StatusReturned      Status = "returned"
	StatusArchived      Status = "archived"
	StatusCompleted     Status = "completed"
)

type Article struct {
	ID           string  `json:"id"`
	Section      Section `json:"section"`
	Title        string  `json:"title"`
	Summary      string  `json:"summary"`
	Body         string  `json:"body"`
	Author       string  `json:"author"`
	Status       Status  `json:"status"`
	Edition      string  `json:"edition"`
	UpdatedLabel string  `json:"updatedLabel"`
}

type ContentUpdate struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Body    string `json:"body"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusDraft, StatusPendingReview, StatusPublished, StatusReturned, StatusArchived, StatusCompleted:
		return true
	default:
		return false
	}
}

func (u ContentUpdate) Valid() bool {
	return strings.TrimSpace(u.Title) != "" && strings.TrimSpace(u.Body) != ""
}
