package middlewares

import (
  "strings"
  "net/http"
)

// Heartbeat
// Router.Use(middlewares.Heartbeat("/healthz", "ok"))
func Heartbeat(endpoint string, text string) func(http.Handler) http.Handler {
  f := func(h http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
      if r.Method == "GET" && strings.EqualFold(r.URL.Path, endpoint) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(text))
        return
      }
      h.ServeHTTP(w, r)
    }
    return http.HandlerFunc(fn)
  }
  return f
}
