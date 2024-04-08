package uptimerobottest

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/clevyr/uptime-robot-operator/internal/uptimerobot/uptimerobottest/responses"
)

func NewServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/{endpoint}", func(w http.ResponseWriter, r *http.Request) {
		endpoint := r.PathValue("endpoint")
		f, err := responses.FS.Open(endpoint + ".json")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		_, _ = io.Copy(w, f)
		_ = f.Close()
	})

	return httptest.NewServer(mux)
}
