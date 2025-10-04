package middlewares

import (
	"log/slog"
	"net/http"
)

const USER = "admin"
const PASSWORD = "admin"

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pwd, ok := r.BasicAuth()
		if !ok || USER != user || PASSWORD != pwd {
			http.Error(w, "401 Unauthorized", 401)
			slog.Info("Basic Auth failed: ", "ok", ok, "user", user, "pwd", pwd)
			return
		}
		slog.Info("User authenticated", "user", user)
		next.ServeHTTP(w, r)
	})
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request Received.. ")
		var (
			contentType = r.Header.Get("Content-Type")
			ip          = r.RemoteAddr
			method      = r.Method
			proto       = r.Proto
			headers     = r.Header
		)
		userAttrs := slog.Group("user", "ip", ip)
		requestAttrs := slog.Group("request", "method", method, "content-type", contentType, "protocol", proto, "headers", headers)
		slog.Info("Request Data: ", "userAttrs", userAttrs, "requestAttrs", requestAttrs)
		next.ServeHTTP(w, r)
	})
}
