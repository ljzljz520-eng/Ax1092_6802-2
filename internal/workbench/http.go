package workbench

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

type transitionRequest struct {
	Status Status `json:"status"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes(static http.Handler) http.Handler {
	router := chi.NewRouter()
	router.Route("/api", func(api chi.Router) {
		api.Get("/articles", h.listArticles)
		api.Get("/articles/{articleID}", h.getArticle)
		api.Patch("/articles/{articleID}", h.updateArticle)
		api.Post("/articles/{articleID}/transitions", h.transitionArticle)
		api.Get("/queues/completed", h.completedQueue)
	})
	if static != nil {
		router.Handle("/*", static)
	}
	return router
}

func (h *Handler) listArticles(w http.ResponseWriter, r *http.Request) {
	section := Section(r.URL.Query().Get("section"))
	status := Status(r.URL.Query().Get("status"))
	writeJSON(w, http.StatusOK, map[string]any{"articles": h.service.List(section, status)})
}

func (h *Handler) getArticle(w http.ResponseWriter, r *http.Request) {
	article, err := h.service.Get(chi.URLParam(r, "articleID"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, article)
}

func (h *Handler) updateArticle(w http.ResponseWriter, r *http.Request) {
	var update ContentUpdate
	if err := decodeJSON(r, &update); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	article, err := h.service.UpdateContent(chi.URLParam(r, "articleID"), update)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, article)
}

func (h *Handler) transitionArticle(w http.ResponseWriter, r *http.Request) {
	var request transitionRequest
	if err := decodeJSON(r, &request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	article, err := h.service.Transition(chi.URLParam(r, "articleID"), request.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, article)
}

func (h *Handler) completedQueue(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"articles": h.service.Completed()})
}

func decodeJSON(r *http.Request, destination any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(destination)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrInvalidTransition):
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrInvalidContent), errors.Is(err, ErrInvalidStatus):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
