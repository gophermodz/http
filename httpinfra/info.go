package httpinfra

import (
	"net/http"
	"os"
	"runtime"
	"strconv"
)

var hostname, _ = os.Hostname()

func (s *Server) infoHandler(w http.ResponseWriter, r *http.Request) {
	_, span := s.tracer.Start(r.Context(), "infoHandler")
	defer span.End()

	data := RuntimeResponse{
		Hostname:     hostname,
		Version:      s.config.Version,
		Revision:     s.config.Revision,
		GOOS:         runtime.GOOS,
		GOARCH:       runtime.GOARCH,
		Runtime:      runtime.Version(),
		NumGoroutine: strconv.FormatInt(int64(runtime.NumGoroutine()), 10),
		NumCPU:       strconv.FormatInt(int64(runtime.NumCPU()), 10),
	}

	s.JSONResponse(w, r, data)
}

type RuntimeResponse struct {
	Hostname     string `json:"hostname"`
	Version      string `json:"version"`
	Revision     string `json:"revision"`
	GOOS         string `json:"goos"`
	GOARCH       string `json:"goarch"`
	Runtime      string `json:"runtime"`
	NumGoroutine string `json:"num_goroutine"`
	NumCPU       string `json:"num_cpu"`
}
