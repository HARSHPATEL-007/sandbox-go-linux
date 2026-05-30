package api

import (
	"encoding/json"
	"net/http"
	
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	
	"goboxd/internal/config"
	"goboxd/internal/engine"
)

func NewRouter(registry map[string]config.Language, eng *engine.Engine) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Post("/run", func(w http.ResponseWriter, r *http.Request) {
		// Security Hole #4: Request size limits
		// Capped at 256KiB at the HTTP layer
		r.Body = http.MaxBytesReader(w, r.Body, 256*1024)

		var req engine.RunRequest // Map to spec
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":{"code":"bad_json"}}`, http.StatusBadRequest)
			return
		}

		// Validation & Queue logic
		// (Assume req.Language is parsed)
		lang, ok := registry["cpp"] // Hardcoded for snippet brevity
		if !ok {
			http.Error(w, `{"error":{"code":"unknown_language"}}`, http.StatusBadRequest)
			return
		}

		resChan := make(chan engine.RunResponse)
		eng.JobQueue <- engine.Job{Req: req, Lang: lang, Result: resChan}
		
		res := <-resChan
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	})

	return r
}