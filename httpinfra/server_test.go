package httpinfra

import (
	"context"
	"net/http"
	"testing"

	"go.uber.org/zap"

	"github.com/gophermodz/http/httptest"
)

func TestServer(t *testing.T) {
	l := zap.NewNop()
	server, err := New(context.Background(), l)
	if err != nil {
		return
	}
	t.Run("healthz", func(t *testing.T) {
		scenarios := []httptest.APIScenario{
			{
				Name:            "success",
				Method:          http.MethodGet,
				URL:             "/api/healthz",
				ExpectedStatus:  http.StatusOK,
				ExpectedContent: []string{`{"status":"ok"}`},
				Handler:         server,
			},
			{
				Name:           "wrong method",
				Method:         http.MethodPost,
				URL:            "/api/healthz",
				ExpectedStatus: http.StatusMethodNotAllowed,
				Handler:        server,
			},
		}
		for _, scenario := range scenarios {
			scenario.Test(t)
		}
	})
}
