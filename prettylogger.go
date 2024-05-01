package prettylogger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var blue = "\033[34m"
var purple = "\033[35m"
var cyan = "\033[36m"
var gray = "\033[37m"
var white = "\033[97m"

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}
func (w *wrappedWriter) FormatCode() string {
	statusString := strconv.Itoa(w.statusCode)
	if w.statusCode < 300 {
		return green + statusString + reset
	} else if w.statusCode < 500 {
		return yellow + statusString + reset
	}
	return red + statusString + reset
}

type PrettyLogger struct {
	Output io.Writer
}

func NewPrettyLogger(writers ...io.Writer) PrettyLogger {
	if len(writers) == 0 {
		return PrettyLogger{Output: os.Stdout}
	}
	return PrettyLogger{Output: io.MultiWriter(writers...)}
}

func (p PrettyLogger) PrettyLoggerMiddleWare(f func(http.ResponseWriter, *http.Request)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		f(wrapped, r)

		fmt.Fprintf(p.Output, "%s %s %s %s %s", start.Format("2006/01/02 15:04:05"), time.Since(start).String(), r.Method, r.URL.Path, wrapped.FormatCode())

	}
}
