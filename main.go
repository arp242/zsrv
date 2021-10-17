package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"zgo.at/zhttp"
	"zgo.at/zstd/znet"
)

var version = ""

//go:embed www/*
var webroot embed.FS

func main() {
	zhttp.Static404 = func(w http.ResponseWriter, r *http.Request) {
		r.RequestURI = strings.ToLower(r.RequestURI)

		// Check redirect.
		if redirs := redirects[znet.RemovePort(strings.ToLower(r.Host))]; redirs != nil {
			// Exact match.
			if r, ok := redirs[r.RequestURI]; ok {
				w.Header().Add("Location", r)
				w.WriteHeader(301)
				return
			}

			// TODO: can optimize by doing this on startup.
			for k, v := range redirs {
				k = strings.ToLower(k)
				if strings.HasSuffix(k, "*") && strings.HasPrefix(r.RequestURI, k[:len(k)-1]) {
					w.Header().Add("Location", strings.Replace(v, "{}", r.RequestURI[len(k):], -1))
					w.WriteHeader(301)
					return
				}
			}
		}

		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(404)
		fmt.Fprintf(w, `Page not found. <a href="/">Back home</a>`)
	}

	var (
		domains []string
		routers = make(map[string]http.Handler)
	)
	dirs, err := fs.ReadDir(webroot, "www")
	if err != nil {
		fmt.Fprintf(os.Stderr, "zsrv: reading www: %s", err)
		os.Exit(1)
	}
	for _, d := range dirs {
		domains = append(domains, d.Name())
		routers[strings.ReplaceAll(d.Name(), "STAR", "*")] = static(d.Name())
	}

	for k, v := range domainRedirects {
		routers[k] = zhttp.RedirectHost("//" + v)
	}

	fmt.Printf("zsrv %s listening on %s\n", version, listen)
	fmt.Printf("    serving domains: %s\n", domains)
	w, err := zhttp.Serve(0, nil, &http.Server{
		Addr:    listen,
		Handler: zhttp.HostRoute(routers),
	})
	if err != nil {
		panic(err)
	}
	<-w
	<-w
}

// TODO: be smarter about cache (per-file and per-filetype)
func static(dir string) http.Handler {
	fsys, err := fs.Sub(webroot, "www/"+dir)
	if err != nil {
		panic(err)
	}

	stat := zhttp.NewStatic("", fsys, map[string]int{"": 86400})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" || r.URL.Path == "/" {
			r.URL.Path = "/index.html"
		}
		stat.ServeHTTP(w, r)
	})
}
