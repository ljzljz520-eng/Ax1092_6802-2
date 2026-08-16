package workbench

import "errors"

var (
	ErrInvalidContent    = errors.New("title and body are required")
	ErrInvalidStatus     = errors.New("invalid target status")
	ErrInvalidTransition = errors.New("status transition is not allowed")
)

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) List(section Section, status Status) []Article {
	return s.store.List(section, status)
}

func (s *Service) Get(id string) (Article, error) {
	return s.store.Get(id)
}

func (s *Service) UpdateContent(id string, update ContentUpdate) (Article, error) {
	if !update.Valid() {
		return Article{}, ErrInvalidContent
	}
	return s.store.UpdateContent(id, update)
}

func (s *Service) Transition(id string, target Status) (Article, error) {
	if !target.Valid() {
		return Article{}, ErrInvalidStatus
	}
	article, err := s.store.Get(id)
	if err != nil {
		return Article{}, err
	}
	if !transitionAllowed(article.Status, target) {
		return Article{}, ErrInvalidTransition
	}
	return s.store.SetStatus(id, target)
}

func transitionAllowed(current, target Status) bool {
	switch current {
	case StatusDraft, StatusReturned:
		return target == StatusPendingReview
	case StatusPendingReview:
		return target == StatusPublished || target == StatusReturned
	case StatusPublished:
		return target == StatusArchived
	default:
		return false
	}
}

func (s *Service) Completed() []Article {
	return s.store.List("", StatusCompleted)
}
