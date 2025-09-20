package ginadapter

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoggerMiddleware_WritesLog(t *testing.T) {
	var buf bytes.Buffer

	prev := log.Writer()
	log.SetOutput(&buf)
	t.Cleanup(func() {
		log.SetOutput(prev)
	})

	r := gin.New()
	r.Use(LoggerMiddleware())
	r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })

	req := httptest.NewRequest("GET", "/ping", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
	if got := buf.String(); got == "" {
		t.Fatalf("expected log output, got empty")
	}
}
