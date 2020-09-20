package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type InstrumentedResponseWriter struct {
	w http.ResponseWriter
	Status int
	Bytes int
}

func (w *InstrumentedResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *InstrumentedResponseWriter) WriteHeader(status int) {
	w.Status = status
	w.w.WriteHeader(status)
}

func (w *InstrumentedResponseWriter) Write(data []byte) (int, error) {
	n, err := w.w.Write(data)
	w.Bytes += n
	return n, err
}

type AccessLogger struct {
	FileName string
	f io.WriteCloser
	lock *sync.Mutex
}

func NewAccessLogger(filename string) (*AccessLogger, error) {
	dn := filepath.Dir(filename)
	st, err := os.Stat(dn)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dn, 0775)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else if !st.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dn)
	}
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &AccessLogger{
		FileName: filename,
		f: f,
		lock: &sync.Mutex{},
	}, nil
}

func (l *AccessLogger) Close() error {
	l.lock.Lock()
	defer l.lock.Unlock()
	f := l.f
	l.f = nil
	if f == nil {
		return errors.New("logger already closed")
	}
	return f.Close()
}

func (l *AccessLogger) Log(iw *InstrumentedResponseWriter, r *http.Request, startTime time.Time) {
	endTime := time.Now()
	dur := endTime.Sub(startTime)
	line := fmt.Sprintf(
		"%s [%s] \"%s %s\" %d %d %.3f \"%s\" \"%s\"\n",
		r.RemoteAddr,
		startTime.Format("2006-01-02 15:04:05 -0700"),
		r.Method,
		r.RequestURI,
		iw.Status,
		iw.Bytes,
		float64(dur.Microseconds()) / 1000,
		r.Header.Get("Referer"),
		r.Header.Get("User-Agent"),
	)
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.f != nil {
		l.f.Write([]byte(line))
	}
}

