package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Yamaguchi-Katsuya/ca-tech-dojo-go/model"
)

func ActionLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.UserAgent()

		al := &model.ActionLog{
			Timestamp: time.Now(),
			Path:      r.URL.Path,
			OS:        ua,
		}

		next.ServeHTTP(w, r)
		al.ElapsedTime = time.Since(al.Timestamp)

		fmt.Printf(
			"[%s] Path: %s | OS: %s | ElapsedTime: %s \n",
			al.Timestamp.Format(time.DateTime),
			al.Path,
			al.OS,
			al.ElapsedTime,
		)
	})
}
