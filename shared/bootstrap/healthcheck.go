package bootstrap

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

func HealthCheck(port string, logger *zap.Logger) {
	handler := http.NewServeMux()

	handler.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = fmt.Fprintf(w, "Alive")
	})

	go func() {
		logger.Info("starting healthcheck", zap.String("addr", ":"+port))

		_ = http.ListenAndServe(":"+port, handler)
	}()
}
