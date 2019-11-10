//go:generate go run gen.go

package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"zgo.at/zhttp"
	"zgo.at/zlog"
)

var version = ""

func main() {
	zhttp.ErrPage = func(w http.ResponseWriter, r *http.Request, code int, reported error) {
		w.WriteHeader(code)
		if code >= 500 {
			zlog.FieldsRequest(r).Error(reported)
			reported = errors.New("internal server error")
		} else {
			zlog.Field("code", code).Print(reported.Error())
		}

		fmt.Fprintf(w, "%d: %s", code, reported)
	}

	zhttp.Static404 = func(w http.ResponseWriter, r *http.Request) {
		r.RequestURI = strings.ToLower(r.RequestURI)

		// Check redirect.
		if redirs := redirects[zhttp.RemovePort(strings.ToLower(r.Host))]; redirs != nil {
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

	var domains []string
	routers := make(map[string]chi.Router)
	for d := range packmap {
		domains = append(domains, d)
		routers[d] = static(d)
	}

	for k, v := range domainRedirects {
		routers[k] = zhttp.RedirectHost("//" + v)
	}

	fmt.Printf("zsrv %s listening on %s\n", version, listen)
	fmt.Printf("    serving domains: %s\n", domains)
	zhttp.Serve(&http.Server{Addr: listen, Handler: zhttp.HostRoute(routers)}, nil)
}

// TODO: be smarter about cache (per-file and per-filetype)
func static(dir string) chi.Router {
	r := chi.NewRouter()
	p := packmap[dir]
	if p == nil {
		panic(fmt.Sprintf("packmap[%q] is nil", dir))
	}
	r.Get("/*", zhttp.NewStatic(dir, dir, 86400, p).ServeHTTP)
	zhttp.MountACME(r, certdir)
	return r
}
