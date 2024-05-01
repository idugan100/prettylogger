/*
prettylogger is a minimalistic server logging framework that prints out the request time, duration, http verb, url path, and color coded response code.
You start by passing all the places you want to log to that satify the io.Writer interface to the NewPrettyLogger function. Common places might be STDOUT or a log file.
This will return a PrettyLogger struct.
You can then use the PrettyLoggerMiddleWare function on the struct to wrap any http.Handlerfunc.
*/
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

/* This struct is the main object of the prettylogger package. You can create one of these yourself with you output locations as writers, or use the NewPrettyLogger function */
type PrettyLogger struct {
	Output io.Writer
}

/* Pass all the places you want to log to as writers and this will  return a new PrettyLogger struct for you to use */
func NewPrettyLogger(writers ...io.Writer) PrettyLogger {
	if len(writers) == 0 {
		return PrettyLogger{Output: os.Stdout}
	}
	return PrettyLogger{Output: io.MultiWriter(writers...)}
}

/* Once you have a PrettyLogger struct, use this function on the struct to wrap http.Handlefunc in order to log all the requests that go though it. */
func (p PrettyLogger) PrettyLoggerMiddleWare(f func(http.ResponseWriter, *http.Request)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		f(wrapped, r)

		fmt.Fprintf(p.Output, "%s | %s %s %s %s\n", start.Format("2006/01/02 15:04:05"), time.Since(start).String(), r.Method, r.URL.Path, wrapped.FormatCode())

	}
}
