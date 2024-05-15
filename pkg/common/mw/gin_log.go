

package mw

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/tools/log"
)

type responseWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			c.Abort()
			return
		}
		start := time.Now()
		log.ZDebug(c, "gin request", "method", c.Request.Method, "uri", c.Request.RequestURI, "req", string(req))
		c.Request.Body = io.NopCloser(bytes.NewReader(req))
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			buf:            bytes.NewBuffer(nil),
		}
		c.Writer = writer
		c.Next()
		resp := writer.buf.Bytes()
		log.ZDebug(c, "gin response", "time", time.Since(start), "status", c.Writer.Status(), "resp", string(resp))
	}
}
