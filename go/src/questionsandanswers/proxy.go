package main

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"
)

func parseAddr(addr string) string {
	parts := strings.Split(addr, ":")
	parts = parts[:len(parts) - 1]
	return strings.Join(parts, ":")
}

var fwdRe = regexp.MustCompile(`([\s,;=]+)=([^\s,;=]+|"[^"]*")([;,$])`)

func parseForwarded(h string) [][]string {
	if h == "" {
		return [][]string{}
	}
	fwd := [][]string{[]string{}}
	idx := 0
	ms := fwdRe.FindAllStringSubmatch(strings.TrimSpace(h), -1)
	for _, m := range ms {
		fwd[idx] = append(fwd[idx], fmt.Sprintf("%s=%s", m[0], m[1]))
		if m[2] == "," {
			fwd = append(fwd, []string{})
			idx += 1
		}
	}
	if len(fwd[idx]) == 0 {
		fwd = fwd[:idx - 1]
	}
	return fwd
}

func formatForwarded(fwd [][]string) string {
	parts := make([]string, len(fwd))
	for i, part := range fwd {
		parts[i] = strings.Join(part, ";")
	}
	return strings.Join(parts, ", ")
}

func (router *Router) Proxy(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{ Timeout: 30 * time.Second }
	defer client.CloseIdleConnections()
	proxyUrl := *router.DefaultProxy
	proxyUrl.Path = path.Join(proxyUrl.Path, strings.TrimLeft(r.URL.Path, "/"))
	req, err := http.NewRequest(r.Method, proxyUrl.String(), r.Body)
	if err != nil {
		SendError(w, BadRequest(err))
		return
	}
	fwd := parseForwarded(req.Header.Get("Forwarded"))
	ip := parseAddr(r.RemoteAddr)
	req.Header.Set("X-Forwarded-Host", r.Host)
	req.Header.Set("X-Forwarded-Proto", r.URL.Scheme)
	req.Header.Set("X-Real-IP", ip)
	for k, vs := range r.Header {
		switch k {
		case "Host", "Forwarded":
		default:
			req.Header[k] = vs
		}
	}
	xff := req.Header.Get("X-Forwarded-For")
	if xff == "" {
		req.Header.Set("X-Forwarded-For", ip)
	} else {
		req.Header.Set("X-Forwarded-For", fmt.Sprintf("%s, %s", xff, ip))
	}
	req.Header.Set("Forwarded", formatForwarded(fwd))
	req.Header.Set("X-Forwarded-For", ip)
	res, err := client.Do(req)
	if err != nil {
		SendError(w, BadGateway(err))
		return
	}
	wh := w.Header()
	for k, vs := range res.Header {
		wh[k] = vs
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
	res.Body.Close()
}
