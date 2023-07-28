package middlewares

import "net/http"

func CORSHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "http://127.0.0.1")
		w.Header().Add("Access-Control-Allow-Header", "Origin, X-Requested-With, Content-Type, Accept")
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
