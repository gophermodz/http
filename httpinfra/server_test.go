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
	server := New(context.Background(), l)
	t.Run("healthz", func(t *testing.T) {
		scenarios := []httptest.APIScenario{
			{
				Name:            "success",
				Method:          http.MethodGet,
				URL:             "/healthz",
				ExpectedStatus:  http.StatusOK,
				ExpectedContent: []string{`{"status":"ok"}`},
				Handler:         server,
			},
			{
				Name:           "wrong method",
				Method:         http.MethodPost,
				URL:            "/healthz",
				ExpectedStatus: http.StatusMethodNotAllowed,
				Handler:        server,
			},
		}
		for _, scenario := range scenarios {
			scenario.Test(t)
		}
	})
}
