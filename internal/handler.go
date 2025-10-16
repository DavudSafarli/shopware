package internal

import (
	"context"
	"log/slog"
	"net/http"
)

// RedirectStorage defines the interface for managing redirect rules and welcome page URL
type RedirectStorage interface {
	FindFullMatchRule(ctx context.Context, canonicalPathQuery string) (r FullMatchRule, ok bool, err error)
	GetWelcomePageURL(ctx context.Context) (string, error)
}

type Handler struct {
	storage RedirectStorage
}

func NewHandler(storage RedirectStorage) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	canonicalPathQuery := canonicalPathQuery(r.URL)

	rule, found, err := h.storage.FindFullMatchRule(r.Context(), canonicalPathQuery)
	if err != nil {
		slog.Error("failed to find full match rule", "canonical", canonicalPathQuery, "err", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if found {
		http.Redirect(w, r, rule.Target, http.StatusMovedPermanently)
		return
	}

	// No match found
	// todo: clarify: requirements says return "hello world" page. But, this might already be a valid page.
	welcomeURL, err := h.storage.GetWelcomePageURL(r.Context())
	if err != nil {
		slog.Error("failed to get welcome page url", "err", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if welcomeURL != "" {
		http.Redirect(w, r, welcomeURL, http.StatusMovedPermanently)
		return
	}

	// if it's hosted on the same domain, then maybe we can return to the root domain?
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
