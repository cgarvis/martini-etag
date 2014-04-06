package etag

import "crypto/md5"
import "fmt"
import "github.com/go-martini/martini"
import "io"
import "net/http"

func ETag() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		resp := &etagResponseWriter{res, req, 200}
		c.MapTo(resp, (*http.ResponseWriter)(nil))
		c.Next()
	}
}

type etagResponseWriter struct {
	rw		http.ResponseWriter
	req		*http.Request
	status	int
}

func (e *etagResponseWriter) Header() http.Header {
	return e.rw.Header()
}

func (e *etagResponseWriter) Write(b []byte) (int, error) {
	if e.req.Method == "GET" && len(b) > 1024 {
		if e.status == 0 || e.status == http.StatusOK {
			hash := md5.New()
			for header := range e.rw.Header() {
				io.WriteString(hash, header)
			}
			io.WriteString(hash, string(b))
			etag := fmt.Sprintf("%x", hash.Sum(nil))

			if match, ok := e.req.Header["If-None-Match"]; ok {
				if match[0] == etag {
					e.WriteHeader(http.StatusNotModified)
					b = []byte{}
				}
			} else {
				e.Header().Set("ETag", etag)
			}
		}
	}

	size, err := e.rw.Write(b)
	return size, err
}

func (e *etagResponseWriter) WriteHeader(s int) {
	e.rw.WriteHeader(s)
	e.status = s
}

