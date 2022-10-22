package httpinfra

import "net/http"

func (s *Server) healthzHandler(writer http.ResponseWriter, r *http.Request) {
	_, span := s.tracer.Start(r.Context(), "healthzHandler")
	defer span.End()
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte(`{"status":"ok"}`))
}
