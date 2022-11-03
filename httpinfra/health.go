package httpinfra

import "net/http"

func (s *Server) healthzHandler(writer http.ResponseWriter, r *http.Request) {
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte(`{"status":"ok"}`))
}
