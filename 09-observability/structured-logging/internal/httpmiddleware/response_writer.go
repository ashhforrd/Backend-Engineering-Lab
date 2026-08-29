package httpmiddleware

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	statusCode  int
	bytes       int
	wroteHeader bool
}

func NewResponseWriter(
	writer http.ResponseWriter,
) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: writer,
		statusCode:     http.StatusOK,
	}
}

func (w *ResponseWriter) WriteHeader(
	statusCode int,
) {
	if w.wroteHeader {
		return
	}

	w.wroteHeader = true
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *ResponseWriter) Write(
	value []byte,
) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}

	written, err := w.ResponseWriter.Write(value)
	w.bytes += written

	return written, err
}

func (w *ResponseWriter) StatusCode() int {
	return w.statusCode
}

func (w *ResponseWriter) BytesWritten() int {
	return w.bytes
}

func (w *ResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
