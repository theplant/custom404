package custom404

import (
	"context"
	"net/http"
)

type notfoundHandler struct {
	mux       http.Handler
	custom404 http.Handler
}

type notFoundWriter struct {
	http.ResponseWriter
	notfound bool
}

func (nfw *notFoundWriter) WriteHeader(status int) {
	if status == 404 {
		nfw.notfound = true
		return
	}
	nfw.ResponseWriter.WriteHeader(status)
}

func (nfw *notFoundWriter) Write(b []byte) (int, error) {
	if nfw.notfound {
		return 0, nil
	}
	return nfw.ResponseWriter.Write(b)
}

func (nf *notfoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nfw := &notFoundWriter{ResponseWriter: w}
	nf.mux.ServeHTTP(nfw, r)
	if nfw.notfound {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		newr := r.WithContext(context.WithValue(r.Context(), "statuscode", http.StatusNotFound))
		nf.custom404.ServeHTTP(w, newr)
	}
}

func WithCustom404(mux http.Handler, custom404 http.Handler) (newMux http.Handler) {
	return &notfoundHandler{mux: mux, custom404: custom404}
}
