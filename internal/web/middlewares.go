package web

import (
	"log"
	"net/http"
)

// SessionMiddleware attaches session to each request
func SessionMiddleware(sm *InMemorySessionManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil || cookie.Value == "" {
			// no session → create one
			session, _ := sm.New()
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    session.ID,
				HttpOnly: true,
				Path:     "/",
			})
			log.Printf("New session created: %s", session.ID)
			// inject session into context if needed
			r = r.WithContext(SetSession(r.Context(), session))
		} else {
			// existing session
			session, err := sm.Get(cookie.Value)
			if err != nil {
				// expired / invalid → create new
				session, _ = sm.New()
				http.SetCookie(w, &http.Cookie{
					Name:     "session_id",
					Value:    session.ID,
					HttpOnly: true,
					Path:     "/",
				})
				log.Printf("Session refreshed: %s", session.ID)
			}
			r = r.WithContext(SetSession(r.Context(), session))
		}
		next.ServeHTTP(w, r)
	})
}
