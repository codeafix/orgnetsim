package srvr

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/spaceweasel/mango"
)

// StaticRouter is a derived version of mango.Router that also serves an embedded
// copy of the orgnetsim UI
type StaticRouter struct {
	Router       *mango.Router
	embedHandler http.Handler
	webpath      string
}

// ListenAndServe launches the web server
// rootpath is the root directory where data served by this server is persisted
// webpath is the root directory of the static website served by this server
// webfs is the embedded copy of the static website that will be served if the
// webpath is an empty string
// port is the port to listen on
func ListenAndServe(rootpath string, webpath string, webfs fs.FS, port string) {
	fm := NewFileManager(rootpath)
	r := CreateRouter(fm)

	r.RequestLogger = func(l *mango.RequestLog) {
		fmt.Println(l.CombinedFormat())
	}
	r.ErrorLogger = func(err error) {
		fmt.Println(err.Error())
	}
	corsConfig := mango.CORSConfig{
		Origins: []string{"*"},
		Methods: []string{"GET", "POST", "PUT", "DELETE"},
	}
	r.SetGlobalCORS(corsConfig)

	sr := StaticRouter{
		Router:  r,
		webpath: webpath,
	}

	if len(webpath) > 0 {
		r.StaticDir(webpath)
	} else {
		var staticFS = http.FS(webfs)
		sr.embedHandler = http.FileServer(staticFS)
	}

	http.ListenAndServe(":"+port, &sr)
}

// CreateRouter registers the route handlers. This function allows the route handlers
// to be tested with the mango.Browser
func CreateRouter(fm FileManager) *mango.Router {
	r := mango.NewRouter()

	r.RegisterModules([]mango.Registerer{
		NewSimListHandler(fm),
		NewSimHandler(fm),
		NewStepHandler(fm),
	})

	return r
}

// ServeHTTP Decides based on the route path whether to use the API router or an embedded
// file server that serves the static files of the embedded website
func (sr *StaticRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := strings.Split(req.URL.Path, "/")
	//Ignore optional first /
	i := 0
	if len(path[i]) == 0 {
		i++
	}
	//If the request path does not begin with /api/, a webpath is not specified on the
	//command line, and there is an embed handler, then try and serve the embedded web files
	if path[i] != "api" && len(sr.webpath) == 0 && sr.embedHandler != nil {
		sr.embedHandler.ServeHTTP(w, req)
		return
	}
	sr.Router.ServeHTTP(w, req)
}
